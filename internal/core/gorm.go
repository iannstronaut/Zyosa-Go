package core

import (
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDB(viper *viper.Viper) (*gorm.DB, error) {
	username := viper.GetString("database.user")
	password := viper.GetString("database.password")
	host := viper.GetString("database.host")
	port := viper.GetInt("database.port")
	dbname := viper.GetString("database.name")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, dbname)

	dbLogger := logger.New(log.New(os.Stdout, "\r\n", 0), logger.Config{
		SlowThreshold:             0,
		LogLevel:                  logger.Silent,
		IgnoreRecordNotFoundError: true,
		ParameterizedQueries:      true,
		Colorful:                  false,
	})

	db, err := gorm.Open(
		mysql.New(mysql.Config{DSN: dsn}),
		&gorm.Config{Logger: dbLogger},
	)

	if err != nil {
		logrus.Fatal(err)
	}

	return db, nil
}
