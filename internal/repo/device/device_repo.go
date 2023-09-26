package devicerepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AZhur771/wg-grpc-api/internal/entity"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type DeviceRepo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *DeviceRepo {
	return &DeviceRepo{
		db: db,
	}
}

func (d *DeviceRepo) Add(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (*entity.Device, error) {
	model := d.toModel(dev)

	query := `
		INSERT INTO device (
				private_key,
				description,
				endpoint,
				fw_mark,
				address,
				mtu,
				dns,
				tble,
				pre_up,
				post_up,
				pre_down,
				post_down
			)
		VALUES (
				:private_key,
				:description,
				:endpoint,
				:fw_mark,
				:address,
				:mtu,
				:dns,
				:tble,
				:pre_up,
				:post_up,
				:pre_down,
				:post_down
			)
		RETURNING *;
	`

	var rows *sqlx.Rows
	var err error

	if tx == nil {
		rows, err = d.db.NamedQueryContext(ctx, query, model)
	} else {
		rows, err = tx.NamedQuery(query, model)
	}

	if err != nil {
		return nil, fmt.Errorf("device repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(model); err != nil {
			return nil, fmt.Errorf("device repo: %w", err)
		}
	}

	return model.ToEntity()
}

func (d *DeviceRepo) Update(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (*entity.Device, error) {
	model := d.toModel(dev)

	query := `
		UPDATE device
		SET description = :description,
			endpoint = :endpoint,
			fw_mark = :fw_mark,
			address = :address,
			mtu = :mtu,
			dns = :dns,
			tble = :tble,
			pre_up = :pre_up,
			post_up = :post_up,
			pre_down = :pre_down,
			post_down = :post_down;
	`

	var rows *sqlx.Rows
	var err error

	if tx == nil {
		rows, err = d.db.NamedQueryContext(ctx, query, model)
	} else {
		rows, err = tx.NamedQuery(query, model)
	}

	if err != nil {
		return nil, fmt.Errorf("device repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.StructScan(model); err != nil {
			return nil, fmt.Errorf("device repo: %w", err)
		}
	}

	return model.ToEntity()
}

func (d *DeviceRepo) Remove(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) error {
	query := "DELETE FROM device WHERE id = $1;"

	var err error

	if tx == nil {
		_, err = d.db.ExecContext(ctx, query, id)
	} else {
		_, err = tx.Exec(query, id)
	}

	if err != nil {
		return fmt.Errorf("device repo: %w", err)
	}

	return nil
}

func (d *DeviceRepo) Get(ctx context.Context, tx *sqlx.Tx, id uuid.UUID) (*entity.Device, error) {
	model := NewModel()

	query := `
		SELECT id,
			private_key,
			"name",
			description,
			endpoint,
			fw_mark,
			address,
			mtu,
			dns,
			tble,
			pre_up,
			post_up,
			pre_down,
			post_down
		FROM device
		WHERE id = $1;
	`

	var row *sqlx.Row

	if tx == nil {
		row = d.db.QueryRowxContext(ctx, query, id)
	} else {
		row = tx.QueryRowx(query, id)
	}

	if err := row.StructScan(model); err != nil {
		return nil, fmt.Errorf("device repo: %w", err)
	}

	return model.ToEntity()
}

func (d *DeviceRepo) GetAll(ctx context.Context, tx *sqlx.Tx, skip, limit int, search string) ([]*entity.Device, error) {
	devs := make([]*entity.Device, 0)

	query := `
		SELECT id,
			private_key,
			"name",
			description,
			endpoint,
			fw_mark,
			address,
			mtu,
			dns,
			tble,
			pre_up,
			post_up,
			pre_down,
			post_down
		FROM device
	`

	if search != "" {
		query += `
			WHERE \"name\" ILIKE '%' || :search || '%'
				OR description ILIKE '%' || :search || '%'
		`
	}

	query += "OFFSET :skip LIMIT :limit;"

	var rows *sqlx.Rows
	var err error

	args := struct {
		Search string
		Skip   int
		Limit  int
	}{
		Search: search,
		Skip:   skip,
		Limit:  limit,
	}

	if tx == nil {
		rows, err = d.db.NamedQueryContext(ctx, query, args)
	} else {
		rows, err = tx.NamedQuery(query, args)
	}

	if err != nil {
		return nil, fmt.Errorf("device repo: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		model := NewModel()

		if err := rows.StructScan(model); err != nil {
			return nil, fmt.Errorf("device repo: %w", err)
		}

		dev, err := model.ToEntity()
		if err != nil {
			return nil, fmt.Errorf("device repo: %w", err)
		}

		devs = append(devs, dev)
	}

	return devs, nil
}

func (d *DeviceRepo) Count(ctx context.Context, tx *sqlx.Tx) (int, error) {
	var count int
	var err error

	query := "SELECT count(1) FROM device;"

	if tx == nil {
		err = d.db.GetContext(ctx, &count, query)
	} else {
		err = tx.Get(&count, query)
	}

	if err != nil {
		return 0, fmt.Errorf("device repo: %w", err)
	}

	return count, nil
}

func (d *DeviceRepo) GenerateAddress(ctx context.Context, tx *sqlx.Tx, dev *entity.Device) (string, error) {
	query := `
		SELECT sub.ip
		FROM (
				SELECT Set_masklen(
						(
							(
								Generate_series(
									1,
									(2 ^ (32 - Masklen($1::cidr)))::INTEGER - 2
								) + $1::cidr
							)::inet
						),
						32
					) AS ip
			) AS sub
		WHERE sub.ip NOT IN (
				SELECT address
				FROM peer_address
				WHERE device_id = $2
			)
			AND sub.ip > Set_masklen($1, 32);
	`

	var addr string
	var err error

	if tx == nil {
		err = d.db.GetContext(ctx, &addr, query, dev.Address, dev.ID)
	} else {
		err = tx.Get(&addr, query, dev.Address, dev.ID)
	}

	if errors.Is(err, sql.ErrNoRows) {
		return addr, ErrRunOutOfAddresses
	} else if err != nil {
		return addr, fmt.Errorf("device repo: %w", err)
	}

	return addr, nil
}

func (d *DeviceRepo) BeginTxx(ctx context.Context, options *sql.TxOptions) (*sqlx.Tx, error) {
	return d.db.BeginTxx(ctx, options)
}

func (d *DeviceRepo) toModel(dev *entity.Device) *DeviceModel {
	return NewModel().FromEntity(dev)
}
