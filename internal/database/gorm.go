package database

import (
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func NewGorm(masterDSN string, replicaDSN ...string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(masterDSN), &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		logrus.WithError(err).WithField("dsn", masterDSN).Error("couldn't connect to the database")
		return nil, err
	}
	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Replicas: lo.Map(append(replicaDSN, masterDSN), func(item string, _ int) gorm.Dialector {
			return mysql.Open(item)
		}),
	})); err != nil {
		logrus.WithError(err).Error("couldn't setup replica databases")
		return nil, err
	}

	sqlDB,_ := db.DB()
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
