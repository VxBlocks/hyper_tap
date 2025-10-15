package timescale

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	slog "github.com/charmbracelet/log"
	"github.com/go-gormigrate/gormigrate/v2"

	"gorm.io/gorm"
)

var onces sync.Map

func GetPostgresGormTypedDB(ctx context.Context, model Migrator) *gorm.DB {
	if reflect.ValueOf(model).Kind() != reflect.Pointer {
		panic("model must be a pointer")
	}

	table := model

	once, _ := onces.LoadOrStore(table.TableName(), new(sync.Once))
	once.(*sync.Once).Do(func() {
		err := MigrateTable(model)
		if err != nil {
			panic(err)
		}
		slog.Info("created table", "name", table.TableName())
	})

	if ctx == nil {
		return GetPostgresGormDB().Model(model)
	}

	return GetPostgresGormDB().Model(model).WithContext(ctx)
}

func MigrateTable(model Migrator) error {

	db := GetPostgresGormDB()
	m := gormigrate.New(db, gormigrate.DefaultOptions, model.Migrations())
	err := m.Migrate()
	slog.Info("migrate table", "name", model.TableName())
	if err != nil {
		m.RollbackLast()
		return fmt.Errorf("error migrate table :%w", err)
	}

	return nil
}

type Table interface {
	TableName() string
}

type Migrator interface {
	TableName() string
	Migrations() []*gormigrate.Migration
}
