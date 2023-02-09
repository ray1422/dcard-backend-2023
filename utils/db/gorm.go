package db

import (
	"log"
	"os"
	"time"

	"github.com/ray1422/dcard-backend-2023/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db  *gorm.DB
	dsn = "host=" + utils.Getenv("DB_HOST", "localhost") + " user=" + utils.Getenv("DB_USER", "postgres") + " password=" + utils.Getenv("DB_PASSWORD", "no_password") + " dbname=" + utils.Getenv("DB_NAME", "dcard") + " port=" + utils.Getenv("DB_PORT", "5432") + " sslmode=" + utils.Getenv("DB_SSL_MODE", "disable") + " TimeZone=" + utils.Getenv("DB_TIMEZONE", "Asia/Taipei") + ""
)

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
			LogLevel:                  logger.Info,            // Log level
			IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,                  // Disable color
		},
	)

	if db == nil {
		db2, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			// Logger: logger.Default.LogMode(logger.Silent),
			Logger: newLogger,
		})
		if err != nil {
			panic(err)
		}
		db = db2
	}

}

// GormDB return gorm.DB instance.
func GormDB() *gorm.DB {
	return db
}
