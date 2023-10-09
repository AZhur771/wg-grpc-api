package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upInit, downInit)
}

func upInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("CREATE SEQUENCE IF NOT EXISTS device_num;")
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS device
			(
				id          		  UUID DEFAULT Gen_random_uuid() PRIMARY KEY,
				name        		  TEXT DEFAULT ('wg' || Nextval('device_num') - 1),
				private_key 		  TEXT NOT NULL,
				description 		  TEXT,
				endpoint    		  TEXT NOT NULL,
				fw_mark     		  INT,
				address     		  INET NOT NULL,
				dns         		  TEXT,
				mtu         		  INT,
				persistent_keep_alive INT,
				tble                  TEXT,
				pre_up                TEXT,
				post_up               TEXT,
				pre_down              TEXT,
				post_down   		  TEXT,
				UNIQUE(name),
				UNIQUE(endpoint)
			); 
		`,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS peer
			(
				id                    UUID DEFAULT Gen_random_uuid() PRIMARY KEY,
				device_id             UUID NOT NULL,
				name                  TEXT NOT NULL,
				private_key           TEXT NOT NULL,
				preshared_key		  TEXT,
				description           TEXT,
				email                 TEXT,
				dns                   TEXT,
				mtu                   INT,
				persistent_keep_alive INT,
				is_enabled            BOOLEAN,
				FOREIGN KEY (device_id) REFERENCES device (id) ON DELETE CASCADE
			); 
		`,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
			CREATE TABLE IF NOT EXISTS peer_address
			(
				id        UUID DEFAULT Gen_random_uuid() PRIMARY KEY,
				peer_id   UUID NOT NULL,
				device_id UUID NOT NULL,
				address   INET NOT NULL,
				FOREIGN KEY (peer_id) REFERENCES peer (id) ON DELETE CASCADE,
				FOREIGN KEY (device_id) REFERENCES device (id) ON DELETE CASCADE,
				UNIQUE(address, device_id)
			);
		`,
	)
	if err != nil {
		return err
	}

	return nil
}

func downInit(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Exec("DROP TABLE device;")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DROP SEQUENCE device_num;")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DROP TABLE peer;")
	if err != nil {
		return err
	}

	_, err = tx.Exec("DROP TABLE peer_address;")
	if err != nil {
		return err
	}

	return nil
}
