package migrators

import (
	userwatchv1 "hyperliquid-server/gen/userwatch/v1"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(UserwatchMigrator{})
}

type UserwatchMigrator struct {
}

func (s UserwatchMigrator) TableName() string {
	return userwatchv1.UserWatchORM{}.TableName()
}

func (s UserwatchMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_userwatch_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&userwatchv1.UserWatchORM{})
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&userwatchv1.UserWatchORM{})
			},
		},
	}
}
