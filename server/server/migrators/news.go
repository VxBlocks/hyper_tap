package migrators

import (
	"hyperliquid-server/handler"
	"hyperliquid-server/monitor"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(NewsMigrator{})
}

type NewsMigrator struct {
}

func (s NewsMigrator) TableName() string {
	return monitor.NewsOrm{}.TableName()
}

func (s NewsMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_news_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&monitor.NewsOrm{}, &handler.NewsReadOrm{})
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE OR REPLACE FUNCTION news_is_read(address text, uuid text) 
					RETURNS boolean AS $$
					BEGIN
						RETURN EXISTS (
							SELECT 1 FROM news_read 
							WHERE news_read.news_id = news_is_read.uuid 
							AND news_read.user_id = news_is_read.address
						);
					END;
					$$ LANGUAGE plpgsql;`).Error
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&monitor.NewsOrm{}, &handler.NewsReadOrm{})
			},
		},
	}
}
