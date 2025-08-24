package migrators

import (
	userwatchv1 "hyperliquid-server/gen/userwatch/v1"
	"hyperliquid-server/monitor"
	"timescale"

	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

func init() {
	register(UserfillProgressMigrator{})
	register(UserfillMsgMigrator{})
}

type UserfillProgressMigrator struct {
}

// Migrations implements timescale.Table.
func (s UserfillProgressMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_userfill_progress_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&monitor.AddressFillProgress{})
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&monitor.AddressFillProgress{})
			},
		},
	}
}

// TableName implements timescale.Table.
func (s UserfillProgressMigrator) TableName() string {
	return monitor.AddressFillProgress{}.TableName()
}

var _ timescale.Migrator = UserfillProgressMigrator{}

type UserfillMsgMigrator struct {
}

// Migrations implements timescale.Table.
func (s UserfillMsgMigrator) Migrations() []*gormigrate.Migration {
	return []*gormigrate.Migration{
		{
			ID: "create_userfill_msg_table",
			Migrate: func(tx *gorm.DB) error {

				err := tx.AutoMigrate(&userwatchv1.UserWatchMsgORM{})
				if err != nil {
					return err
				}

				return nil
			},
			Rollback: func(tx *gorm.DB) error {
				return tx.Migrator().DropTable(&userwatchv1.UserWatchMsgORM{})
			},
		},
	}
}

// TableName implements timescale.Table.
func (s UserfillMsgMigrator) TableName() string {
	return userwatchv1.UserWatchMsgORM{}.TableName()
}

var _ timescale.Migrator = UserfillProgressMigrator{}
