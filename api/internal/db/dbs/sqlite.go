package dbs

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/andebugan/sheetstruct/internal/app/dtos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Initializes and returns a sqlite database connection
func InitSqliteDatabase() (*gorm.DB, error) {
	// Ensure database directory exists
	if err := os.MkdirAll("database", 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Open SQLite database
	db, err := gorm.Open(sqlite.Open("database/sqlite.db"), &gorm.Config{
		Logger:                                   newLogger,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// AutoMigrate - Creates tables automatically
	err = db.AutoMigrate(
		&dtos.User{},
		&dtos.Sheet{},
		&dtos.Component{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	// Get generic database object for connection pool configuration
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully")
	return db, nil
}

// Closes the database connection
func CloseSqliteDatabase(db *gorm.DB) {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting database instance: %v", err)
			return
		}
		sqlDB.Close()
		log.Println("Database connection closed")
	}
}
