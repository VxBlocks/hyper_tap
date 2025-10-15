package migrators

import (
	"hyperliquid-server/handler"
	"hyperliquid-server/monitor"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(PriceMigrator{})
}

type PriceMigrator struct {
}

// Migrations implements timescale.Table.
func (p PriceMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_price_alerts_tables",
			Migrate: func(tx *gorm.DB) error {
				err := tx.AutoMigrate(&monitor.PriceAlert{}, &monitor.PriceAlertState{}, &handler.PriceAlertReadOrm{})
				if err != nil {
					return err
				}

				err = tx.Exec(`CREATE OR REPLACE FUNCTION price_alert_is_read(address text, alert_id text) 
					RETURNS boolean AS $$
					BEGIN
						RETURN EXISTS (
							SELECT 1 FROM news_read 
							WHERE price_alert_read.alert_id = price_alert_is_read.alert_id 
							AND price_alert_read.user_id = price_alert_is_read.address
						);
					END;
					$$ LANGUAGE plpgsql;`).Error
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				err := tx.Exec(`DROP FUNCTION IF EXISTS price_alert_is_read;`).Error
				if err != nil {
					return err
				}

				err = tx.Migrator().DropTable(&monitor.PriceAlert{}, &monitor.PriceAlertState{}, &handler.PriceAlertReadOrm{})
				if err != nil {
					return err
				}

				return nil
			},
		},
	}
}

// TableName implements timescale.Table.
func (p PriceMigrator) TableName() string {
	return "price_alerts"
}
