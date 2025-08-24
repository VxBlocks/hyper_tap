package migrators

import "timescale"

var migrators = []timescale.Migrator{}

func register(migrator timescale.Migrator) {
	migrators = append(migrators, migrator)
}

func Migrate() error {
	for _, migrator := range migrators {
		err := timescale.MigrateTable(migrator)
		if err != nil {
			return err
		}
	}
	return nil
}
