package main

import (
	"fmt"
	"os"

	"github.com/ThuraMinThein/seats/config"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	config.LoadConfig()

	config := config.Config
	sslMode := "require"
	if config.Environment == "development" {
		sslMode = "disable"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.DBHost,
		config.DBUser,
		config.DBPassword,
		config.DBName,
		config.DBPort,
		sslMode,
	)

	var err error

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("Database connection failed in migration", err)
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		logrus.Fatal("Failed to get sql.DB from GORM", err)
	}

	driver, err := migrate_postgres.WithInstance(sqlDB, &migrate_postgres.Config{})
	if err != nil {
		logrus.Fatal("Failed to create migration driver", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgresql",
		driver,
	)

	if err != nil {
		logrus.Fatal("Migration initialization failed", err)
	}

	cmd := os.Args[(len(os.Args) - 1)]

	switch cmd {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			logrus.Fatal("Migration failed: ", err)
		}
		logrus.Info("Migrations ran successfully")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			logrus.Fatal("Migration failed: ", err)
		}
		logrus.Info("Migrations rolled back successfully")
	default:
		logrus.Fatal("Invalid command. Use 'up' or 'down'.")
	}

	logrus.Info("Migrations ran successfully")
}
