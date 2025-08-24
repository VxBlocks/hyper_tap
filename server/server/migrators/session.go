package migrators

import (
	"hyperliquid-server/models"
	"timescale"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(SessionMigrator{})
}

type SessionMigrator struct {
}

// Migrations implements timescale.Table.
func (s SessionMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_sessions_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&models.SessionORMV1{})
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE OR REPLACE FUNCTION is_authenticated(address text, session text) 
					RETURNS boolean AS $$
					BEGIN
						RETURN EXISTS (
							SELECT 1 FROM sessions 
							WHERE sessions.address = is_authenticated.address 
							AND sessions.session = is_authenticated.session
						);
					END;
					$$ LANGUAGE plpgsql;`).Error
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE OR REPLACE FUNCTION user_id_from_session(session text)
						RETURNS text AS $$
					DECLARE
						result text;
					BEGIN
						SELECT address INTO result FROM sessions 
						WHERE sessions.session = user_id_from_session.session;
						
						IF result IS NULL THEN
							RAISE EXCEPTION 'Session not found: %', user_id_from_session.session;
						END IF;
						
						RETURN result;
					END;
					$$ LANGUAGE plpgsql;`).Error
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				err := tx.Exec(`DROP FUNCTION IF EXISTS user_id_from_session(text)`).Error
				if err != nil {
					return err
				}
				return tx.Migrator().DropTable(&models.SessionORMV1{})
			},
		},
	}
}

// TableName implements timescale.Table.
func (s SessionMigrator) TableName() string {
	return "sessions"
}

var _ timescale.Migrator = SessionMigrator{}
