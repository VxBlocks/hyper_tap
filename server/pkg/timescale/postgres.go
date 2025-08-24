package timescale

import (
	"context"
	"fmt"
	"sync"
	"time"

	"config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var _DB *gorm.DB

var _DB_once sync.Once

func GetPostgresGormDB(ctx ...context.Context) *gorm.DB {
	_DB_once.Do(func() {
		// https://github.com/go-gorm/postgres
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN: fmt.Sprintf("user=%s password=%s dbname=%s port=%s  host=%s sslmode=disable TimeZone=Asia/Shanghai",
				config.C.Postgres.User,
				config.C.Postgres.Password,
				config.C.Postgres.DBName,
				config.C.Postgres.Port,
				config.C.Postgres.Url),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{
			TranslateError:           false,
			SkipDefaultTransaction:   true,
			DisableNestedTransaction: true,
			Logger: &GormLogger{
				IgnoreRecordNotFoundError: true,
				LogLevel:                  gormlogger.Warn,
				SlowThreshold:             200 * time.Millisecond,
			},
		})
		if err != nil {
			panic(fmt.Errorf("error openning database: %w", err))
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Errorf("error read sqlDB: %w", err))
		}
		sqlDB.SetMaxOpenConns(40)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(time.Hour)

		_DB = db

	})

	if len(ctx) > 0 {
		return _DB.WithContext(ctx[0])
	}

	return _DB
}
