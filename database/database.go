package database

import (
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// PUBLIC TYPES
// ========================================================================

/*
Migration callback function
*/
type MigrationFn func(*gorm.DB)

/*
Configuration for database connection
*/
type DatabaseConnectionConfig struct {
	DCS             string      // Connection string
	MigrateCallback MigrationFn // Callback function for migrating database content
}

/*
Database connection struct that provides automatic health checks
*/
type DatabaseConnection struct {
	mx       sync.RWMutex
	DB       *gorm.DB
	cfg      DatabaseConnectionConfig
	migrated bool
}

// PRIVATE TYPES
// ========================================================================

// PUBLIC FUNCTIONS
// ========================================================================

/*
Initialize a database connection with a given database connection string, returns
error on error, nil otherwise
*/
func (dc *DatabaseConnection) Initialize(cfg DatabaseConnectionConfig) error {
	dc.DB = nil
	dc.cfg = cfg
	dc.migrated = false

	return nil
}

/*
Middleware for gin that requires a database connection
*/
func RequireDatabaseConnectionMiddleware(dc *DatabaseConnection) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := RequireDatabaseConnection(dc)
		if err != nil {
			log.Println(err)
			ctx.Abort()
			ctx.Status(http.StatusServiceUnavailable)
			return
		}

		ctx.Next()
	}
}

/*
Require the database connection to be valid

Returns nil if valid, an error otherwise
*/
func RequireDatabaseConnection(dc *DatabaseConnection) error {
	// Avoid processing requests since service unavailable
	if !dc.connect() {
		return errors.New("can't connect to database")
	}

	if !dc.migrated && dc.cfg.MigrateCallback != nil {
		// Migrate
		log.Println("Migrating database")
		dc.mx.Lock()
		defer dc.mx.Unlock()

		dc.cfg.MigrateCallback(dc.DB)
		dc.migrated = true
	}

	return nil
}

// PRIVATE FUNCTIONS
// ========================================================================

/*
Connect to database
*/
func (dc *DatabaseConnection) connect() bool {
	var err error

	if dc.DB == nil {
		// Open
		dc.DB, err = gorm.Open(mysql.Open(dc.cfg.DCS), &gorm.Config{})
		if err != nil {
			// Failed to open connection
			log.Println(err)
			dc.DB = nil
			return false
		}

		return true
	} else {
		// Ping
		db, err := dc.DB.DB()

		if err != nil {
			// We failed to get a DB instance?
			log.Println(err)
			return false
		}

		err = db.Ping()
		if err != nil {
			dc.DB = nil
			return false
		}

		return true
	}
}
