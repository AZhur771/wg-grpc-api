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
				device_id,
				private_key,
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
				:device_id,
				:private_key,
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

	type allowedIP struct {
		PeerID   string `db:"peer_id" sql:",type:uuid"`
		DeviceID string `db:"device_id" sql:",type:uuid"`
		Addr     string `sql:",type:inet"`
	}

	queryAddr := `
		INSERT INTO peer_address (peer_id, device_id, address)
		VALUES (:peer_id, :device_id, :addr);
	`

	allowedIPs := make([]allowedIP, 0, len(peer.AllowedIPs))

	for _, addr := range peer.AllowedIPs {
		allowedIPs = append(allowedIPs, allowedIP{
			PeerID:   model.ID.String(),
			DeviceID: model.DeviceID.String(),
			Addr:     addr,
		})
	}

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
		UPDATE peer
		SET preshared_key = :preshared_key,
			"name" = :name,
			description = :description,
			email = :email,
			dns = :dns,
			mtu = :mtu,
			persistent_keep_alive = :persistent_keep_alive,
			is_enabled = :is_enabled
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
		_, err = p.db.ExecContext(ctx, query, id.String())
	} else {
		_, err = tx.Exec(query, id.String())
	}

	if err != nil {
		return fmt.Errorf("peer repo: %w", err)
	}

	return nil
}

func (p *PeerRepo) Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.Peer, error) {
	model := NewModel()

	rows, err := p.db.QueryxContext(ctx, `
			SELECT p.id as id,
					p.device_id as device_id,
					p.private_key as private_key,
					p.preshared_key as preshared_key,
					p."name" as "name",
					p.description as description,
					p.email as email,
					p.persistent_keep_alive as persistent_keep_alive,
					p.is_enabled as is_enabled,
					p.mtu as mtu,
					p.dns as dns,
					pa.address as address
			FROM peer p
				JOIN peer_address pa
				ON   p.id = pa.peer_id
			WHERE p.id = $1;
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
	mapper := make(map[uuid.UUID]*PeerModel)

	query := `
		SELECT t.id as id,
				t.device_id as device_id,
				t.private_key as private_key,
				t.preshared_key as preshared_key,
				t."name" as "name",
				t.description as description,
				t.email as email,
				t.persistent_keep_alive as persistent_keep_alive,
				t.is_enabled as is_enabled,
				t.mtu as mtu,
				t.dns as dns,
				pa.address as address
		FROM (
				SELECT * FROM peer
				WHERE true
				###device_id###
				###search###
				OFFSET :skip LIMIT :limit
			) t
		JOIN peer_address pa
		ON 	 t.id = pa.peer_id;
	`

	if deviceID != uuid.Nil {
		query = strings.Replace(query, "###device_id###", "AND device_id = :device_id", 1)
	} else {
		query = strings.Replace(query, "###device_id###", "", 1)
	}

	if search != "" {
		query = strings.Replace(query,
			"###search###",
			"AND \"name\" ILIKE '%' || :search || '%' OR description ILIKE '%' || :search || '%'", 1)
	} else {
		query = strings.Replace(query, "###search###", "", 1)
	}

	var rows *sqlx.Rows
	var err error

	args := struct {
		Search   string
		Skip     int
		Limit    int
		DeviceID string `db:"device_id" sql:",type:uuid"`
	}{
		Search:   search,
		Skip:     skip,
		Limit:    limit,
		DeviceID: deviceID.String(),
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
		model := NewModel()

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

		p, ok := mapper[model.ID]
		if !ok {
			model.AllowedIPs = append(model.AllowedIPs, allowedIP)
			mapper[model.ID] = model
		} else {
			p.AllowedIPs = append(p.AllowedIPs, allowedIP)
		}
	}

	peers := make([]*entity.Peer, 0, len(mapper))

	for _, model := range mapper {
		peer, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("device repo: %w", err)
		}

		peers = append(peers, peer)
	}

	return peers, nil
}

func (p *PeerRepo) Count(ctx context.Context, tx *sqlx.Tx, deviceID uuid.UUID) (int, error) {
	query := "SELECT count(1) FROM peer ###device_id###;"

	if deviceID != uuid.Nil {
		query = strings.Replace(query, "###device_id###", "WHERE device_id = :device_id", 1)
	} else {
		query = strings.Replace(query, "###device_id###", "", 1)
	}

	var count int
	var rows *sqlx.Rows
	var err error

	args := struct {
		DeviceID string `db:"device_id" sql:",type:uuid"`
	}{
		DeviceID: deviceID.String(),
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
