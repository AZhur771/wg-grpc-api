package peerrepo

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PeerRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *PeerRepo {
	return &PeerRepo{
		db: db,
	}
}

func (p *PeerRepo) Add(ctx context.Context, tx *sqlx.Tx, peer *entity.Peer) (*entity.Peer, error) {
	model := p.toModel(peer)

	queryPeer := `
		INSERT INTO peer (
				device_id private_key,
				preshared_key,
				"name",
				description,
				email,
				persistent_keep_alive,
				dns,
				mtu,
				is_enabled
			)
		VALUES (
				:device_id :private_key,
				:preshared_key,
				:name,
				:description,
				:email,
				:persistent_keep_alive,
				:dns,
				:mtu,
				:is_enabled
			)
		RETURNING *;
	`

	var rows *sqlx.Rows
	var err error

	if tx == nil {
		rows, err = p.db.NamedQueryContext(ctx, queryPeer, model)
	} else {
		rows, err = tx.NamedQuery(queryPeer, model)
	}

	if err != nil {
		return nil, fmt.Errorf("peer repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(model); err != nil {
			return nil, fmt.Errorf("peer repo: %w", err)
		}
	}

	allowedIPs := make([]struct {
		peerID   string
		deviceID string
		addr     string
	}, len(peer.AllowedIPs))

	for i, allowedIP := range peer.AllowedIPs {
		allowedIPs[i] = struct {
			peerID   string
			deviceID string
			addr     string
		}{
			peerID:   model.ID.String(),
			deviceID: model.DeviceID.String(),
			addr:     allowedIP,
		}
	}

	queryAddr := `
		INSERT INTO peer (peer_id, device_id, address)
		VALUES (:peerID, :deviceID, :addr);
	`

	if tx == nil {
		_, err = tx.NamedExecContext(ctx, queryAddr, allowedIPs)
	} else {
		_, err = tx.NamedExec(queryAddr, allowedIPs)
	}

	if err != nil {
		return nil, fmt.Errorf("peer repo: %w", err)
	}

	return model.ToEntity()
}

func (p *PeerRepo) Update(ctx context.Context, tx *sqlx.Tx, peer *entity.Peer) (*entity.Peer, error) {
	model := p.toModel(peer)

	query := `
		UPDATE device
		SET preshared_key = :preshared_key,
			"name" = :name,
			description = :description,
			email = :email,
			dns = :dns,
			mtu = :mtu,
			persistent_keep_alive = :persistent_keep_alive is_enabled = :is_enabled
		RETURNING *;
	`

	var rows *sqlx.Rows
	var err error

	if tx == nil {
		rows, err = p.db.NamedQueryContext(ctx, query, model)
	} else {
		rows, err = tx.NamedQuery(query, model)
	}

	if err != nil {
		return nil, fmt.Errorf("peer repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(model); err != nil {
			return nil, fmt.Errorf("peer repo: %w", err)
		}
	}

	return model.ToEntity()
}

func (p *PeerRepo) Remove(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	query := "DELETE FROM peer WHERE id = $1;"

	var err error

	if tx == nil {
		_, err = p.db.ExecContext(ctx, query, id)
	} else {
		_, err = tx.Exec(query, id)
	}

	if err != nil {
		return fmt.Errorf("peer repo: %w", err)
	}

	return nil
}

func (p *PeerRepo) Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.Peer, error) {
	model := NewModel()

	rows, err := p.db.QueryxContext(ctx, `
			SELECT id,
					device_id,
					private_key,
					preshared_key,
					"name",
					description,
					email,
					persistent_keep_alive,
					is_enabled,
					mtu,
					dns,
					pa.address
			FROM   peer p
				JOIN peer_address pa
				ON   p.id = pa.peer_id
			WHERE  id = $1;
		`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var allowedIP string

		if err := rows.Scan(
			&model.ID,
			&model.DeviceID,
			&model.PrivateKey,
			&model.PresharedKey,
			&model.Name,
			&model.Description,
			&model.Email,
			&model.PersistentKeepAlive,
			&model.IsEnabled,
			&model.Mtu,
			&model.DNS,
			&allowedIP,
		); err != nil {
			return nil, fmt.Errorf("storage: %w", err)
		}

		model.AllowedIPs = append(model.AllowedIPs, allowedIP)
	}

	return model.ToEntity()
}

func (p *PeerRepo) GetAll(ctx context.Context, tx *sqlx.Tx, skip, limit int, search string, deviceID uuid.UUID) ([]*entity.Peer, error) {
	mapper := make(map[uuid.UUID]*entity.Peer)

	query := `
		SELECT *
		FROM (
				SELECT id,
					device_id,
					private_key,
					preshared_key,
					"name",
					description,
					email,
					persistent_keep_alive,
					is_enabled,
					mtu,
					dns,
					pa.address
				FROM peer p
					JOIN peer_address pa ON p.id = pa.peer_id
				WHERE true
				###device_id###
				###search###
			) t
	`

	if deviceID != uuid.Nil {
		query = strings.Replace(query, "###device_id###", "AND device_id := device_id", 1)
	}

	if search != "" {
		query = strings.Replace(query,
			"###search###",
			"AND \"name\" ILIKE '%' || :search || '%' OR description ILIKE '%' || :search || '%'", 1)
	}

	query += "OFFSET :skip LIMIT :limit;"

	var rows *sqlx.Rows
	var err error

	args := struct {
		Search   string
		Skip     int
		Limit    int
		DeviceID uuid.UUID `db:"device_id" sql:",type:uuid"`
	}{
		Search:   search,
		Skip:     skip,
		Limit:    limit,
		DeviceID: deviceID,
	}

	if tx == nil {
		rows, err = p.db.NamedQueryContext(ctx, query, args)
	} else {
		rows, err = tx.NamedQuery(query, args)
	}

	if err != nil {
		return nil, fmt.Errorf("storage: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var peer *entity.Peer
		var allowedIP string

		err := rows.Scan(&peer.ID, &peer.DeviceID, &peer.PrivateKey, &peer.PresharedKey, &peer.Name, &peer.Description,
			&peer.Email, &peer.PersistentKeepaliveInterval, &peer.IsEnabled, &allowedIP)

		p, ok := mapper[peer.ID]
		if !ok {
			peer.AllowedIPs = append(peer.AllowedIPs, allowedIP)
			mapper[peer.ID] = peer
		} else {
			p.AllowedIPs = append(p.AllowedIPs, allowedIP)
		}

		if err != nil {
			return nil, fmt.Errorf("storage: %w", err)
		}

		peer.AllowedIPs = append(peer.AllowedIPs, allowedIP)
	}

	peers := make([]*entity.Peer, 0, len(mapper))

	for _, p := range mapper {
		peers = append(peers, p)
	}

	return peers, nil
}

func (p *PeerRepo) Count(ctx context.Context, tx *sqlx.Tx, deviceID uuid.UUID) (int, error) {
	query := "SELECT count(1) FROM peer ###device_id###;"

	if deviceID != uuid.Nil {
		query = strings.Replace(query, "###device_id###", "WHERE device_id := device_id", 1)
	}

	var count int
	var rows *sqlx.Rows
	var err error

	args := struct {
		deviceID uuid.UUID `db:"device_id" sql:",type:uuid"`
	}{
		deviceID: deviceID,
	}

	if tx == nil {
		rows, err = p.db.NamedQueryContext(ctx, query, args)
	} else {
		rows, err = tx.NamedQuery(query, args)
	}

	if err != nil {
		return count, fmt.Errorf("peer repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return count, fmt.Errorf("peer repo: %w", err)
		}
	}

	return count, nil
}

func (p *PeerRepo) BeginTxx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return p.db.BeginTxx(ctx, options)
}

func (p *PeerRepo) toModel(peer *entity.Peer) *PeerModel {
	return NewModel().FromEntity(peer)
}
