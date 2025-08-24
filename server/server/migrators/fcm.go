package migrators

import (
	"hyperliquid-server/models"
	"timescale"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(FcmTokenMigrator{})
}

type FcmTokenMigrator struct {
}

// Migrations implements timescale.Table.
func (s FcmTokenMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_fcm_token_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&models.FcmToken{})
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {

				return tx.Migrator().DropTable(&models.FcmToken{})
			},
		},
	}
}

// TableName implements timescale.Table.
func (s FcmTokenMigrator) TableName() string {
	return models.FcmToken{}.TableName()
}

var _ timescale.Migrator = FcmTokenMigrator{}
