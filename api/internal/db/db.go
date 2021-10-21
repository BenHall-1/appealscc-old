package db

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/benhall-1/appealscc/api/internal/models/model"
	"github.com/getsentry/sentry-go"
)

var DB *gorm.DB

func Open() error {
	var err error

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)

	DB, err = gorm.Open(mysql.Open(os.Getenv("DB_CONNECTION")), &gorm.Config{
		// Uncomment the following line when running the first migration for the database
		// DisableForeignKeyConstraintWhenMigrating: true,
		Logger: newLogger,
	})
	if err != nil {
		sentry.CaptureException(err)
		return err
	}

	return nil
}

func Migrate() {
	DB.AutoMigrate(model.User{}, model.Organisation{}, model.Appeal{}, model.AppealResponse{}, model.AppealTemplate{}, model.AppealTemplateField{})
}
