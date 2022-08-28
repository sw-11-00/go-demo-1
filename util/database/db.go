package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Open(user, pass, host string, port int64, dbName string) *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		logrus.Fatal(err)
	}
	return db
}

func TransactionTemplate(db *gorm.DB, f func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if err := f(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func TransactionTemplateWithResult(db *gorm.DB, f func(tx *gorm.DB) (interface{}, error)) (interface{}, error) {
	tx := db.Begin()
	result, err := f(tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return result, tx.Commit().Error
}
