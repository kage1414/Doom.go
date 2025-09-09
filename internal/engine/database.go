package engine

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"doomlike/ent"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	client *ent.Client
}

// NewDatabase creates a new database connection and initializes the schema
func NewDatabase() (*Database, error) {
	// Create data directory if it doesn't exist
	dataDir := "data"
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open SQLite database
	dbPath := filepath.Join(dataDir, "doomlike.db")
	db, err := sql.Open("sqlite3", dbPath+"?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create ent client
	drv := entsql.OpenDB(dialect.SQLite, db)
	client := ent.NewClient(ent.Driver(drv))

	// Run migrations
	ctx := context.Background()
	if err := client.Schema.Create(ctx); err != nil {
		return nil, fmt.Errorf("failed to create schema: %w", err)
	}

	return &Database{client: client}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.client.Close()
}

// LoadSettings loads game settings from the database
func (db *Database) LoadSettings() (*gameSettings, error) {
	ctx := context.Background()

	// Try to get existing settings
	settings, err := db.client.GameSettings.Get(ctx, "default")
	if err != nil {
		if ent.IsNotFound(err) {
			// Create default settings if none exist
			return db.createDefaultSettings(ctx)
		}
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	return &gameSettings{
		fireRate:    settings.FireRate,
		bulletSpeed: settings.BulletSpeed,
		levelCount:  settings.LevelCount,
	}, nil
}

// SaveSettings saves game settings to the database
func (db *Database) SaveSettings(settings *gameSettings) error {
	ctx := context.Background()

	// Check if settings exist
	_, err := db.client.GameSettings.Get(ctx, "default")
	if err != nil {
		if ent.IsNotFound(err) {
			// Create new settings
			_, err = db.client.GameSettings.Create().
				SetID("default").
				SetFireRate(settings.fireRate).
				SetBulletSpeed(settings.bulletSpeed).
				SetLevelCount(settings.levelCount).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()).
				Save(ctx)
			if err != nil {
				return fmt.Errorf("failed to create settings: %w", err)
			}
		} else {
			return fmt.Errorf("failed to check existing settings: %w", err)
		}
	} else {
		// Update existing settings
		_, err = db.client.GameSettings.UpdateOneID("default").
			SetFireRate(settings.fireRate).
			SetBulletSpeed(settings.bulletSpeed).
			SetLevelCount(settings.levelCount).
			SetUpdatedAt(time.Now()).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update settings: %w", err)
		}
	}

	return nil
}

// createDefaultSettings creates default settings in the database
func (db *Database) createDefaultSettings(ctx context.Context) (*gameSettings, error) {
	settings := &gameSettings{
		fireRate:    defaultFireRate,
		bulletSpeed: defaultBulletSpeed,
		levelCount:  defaultLevelCount,
	}

	_, err := db.client.GameSettings.Create().
		SetID("default").
		SetFireRate(settings.fireRate).
		SetBulletSpeed(settings.bulletSpeed).
		SetLevelCount(settings.levelCount).
		SetCreatedAt(time.Now()).
		SetUpdatedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create default settings: %w", err)
	}

	log.Println("Created default game settings in database")
	return settings, nil
}
