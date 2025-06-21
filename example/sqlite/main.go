package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"template-golang/database"
)

func main() {
	// Example 1: Basic SQLite database usage
	fmt.Println("=== SQLite Database Migration Example ===")

	// Create a temporary directory for this example
	tempDir, err := os.MkdirTemp("", "sqlite_migration_example")
	if err != nil {
		log.Fatal("Failed to create temp directory:", err)
	}
	defer os.RemoveAll(tempDir)

	// Create database file path
	dbPath := filepath.Join(tempDir, "example.db")
	fmt.Printf("Using database file: %s\n", dbPath)

	// Example 2: Create SQLite database with default migration path
	fmt.Println("\n--- Creating SQLite database with default settings ---")
	
	config := &database.Config{
		DSN:           dbPath,
		LogMode:       true,
		MigrationPath: "", // Will use default "file://db/migrations"
	}

	db, err := database.NewSqliteDatabase(config)
	if err != nil {
		log.Fatal("Failed to create SQLite database:", err)
	}
	defer db.Close()

	fmt.Printf("✓ SQLite database created successfully\n")

	// Example 3: Check migration version (will show error if no migrations table exists yet)
	fmt.Println("\n--- Checking migration version ---")
	version, dirty, err := db.GetVersion()
	if err != nil {
		fmt.Printf("Migration version check failed (expected for fresh DB): %v\n", err)
	} else {
		fmt.Printf("Current migration version: %d, dirty: %t\n", version, dirty)
	}

	// Example 4: Using provider functions
	fmt.Println("\n--- Using provider functions ---")
	
	db2, err := database.ProvideSqliteDatabase(":memory:", false)
	if err != nil {
		log.Fatal("Failed to create SQLite database via provider:", err)
	}
	defer db2.Close()

	fmt.Printf("✓ SQLite in-memory database created via provider\n")

	// Example 5: Using provider with custom migration path
	fmt.Println("\n--- Using provider with custom migration path ---")
	
	customMigrationPath := "file://" + filepath.Join(tempDir, "custom_migrations")
	db3, err := database.ProvideSqliteDatabaseWithMigrationPath(
		":memory:", 
		true, 
		customMigrationPath,
	)
	if err != nil {
		log.Fatal("Failed to create SQLite database with custom migration path:", err)
	}
	defer db3.Close()

	fmt.Printf("✓ SQLite database created with custom migration path: %s\n", customMigrationPath)

	// Example 6: Demonstrate GORM operations
	fmt.Println("\n--- Testing GORM operations ---")
	
	gormDB := db2.GetDb()
	
	// Create a simple table using GORM AutoMigrate (not related to golang-migrate)
	type User struct {
		ID   uint   `gorm:"primaryKey"`
		Name string
	}

	if err := gormDB.AutoMigrate(&User{}); err != nil {
		log.Fatal("Failed to auto-migrate:", err)
	}

	// Insert a test record
	user := User{Name: "Test User"}
	if err := gormDB.Create(&user).Error; err != nil {
		log.Fatal("Failed to create user:", err)
	}

	// Query the record
	var count int64
	gormDB.Model(&User{}).Count(&count)
	fmt.Printf("✓ Created and queried User table, record count: %d\n", count)

	fmt.Println("\n=== Example completed successfully! ===")
	fmt.Println("\nNote: To use migrations properly, you would need to:")
	fmt.Println("1. Create migration files in db/migrations/ directory")
	fmt.Println("2. Name them like: 000001_create_users.up.sql and 000001_create_users.down.sql")
	fmt.Println("3. Run db.MigrateUp() to apply migrations")
	fmt.Println("4. Run db.MigrateDown(n) to rollback n migrations")
}
