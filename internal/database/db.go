package database

import (
	"database/sql"
	"fmt"
	"golang-grcp-user-services/internal/logger"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(dbFile string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := seedDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to seed database: %v", err)
	}

	return db, nil
}

func seedDatabase(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		fname TEXT NOT NULL,
		city TEXT NOT NULL,
		phone INTEGER NOT NULL,
		height REAL NOT NULL,
		married BOOLEAN NOT NULL
	);`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}

	users := []struct {
		Fname   string
		City    string
		Phone   int64
		Height  float32
		Married bool
	}{
		{"Rahul", "Mumbai", 1234567890, 5.9, true},
		{"Priya", "Delhi", 9876543210, 5.5, true},
		{"Suresh", "Chennai", 5551234567, 6.0, false},
		{"Deepika", "Bangalore", 9998887776, 5.7, false},
		{"Amit", "Kolkata", 1112223334, 5.10, true},
		{"Pooja", "Hyderabad", 4445556667, 5.6, true},
		{"Ravi", "Pune", 7778889990, 5.8, false},
		{"Anjali", "Jaipur", 3334445556, 5.4, false},
		{"Vikram", "Ahmedabad", 6667778889, 6.1, true},
		{"Neha", "Surat", 2223334445, 5.7, true},
		{"Sanjay", "Lucknow", 8889990001, 5.9, true},
		{"Aarti", "Chandigarh", 5556667778, 5.9, false},
		{"Rahul", "Indore", 3332221114, 5.10, true},
		{"Vinay", "Bhopal", 9998887776, 5.8, true},
		{"Anita", "Nagpur", 4445556667, 5.6, false},
		{"Sandeep", "Patna", 7778889990, 5.7, false},
		{"Geeta", "Vadodara", 3334445556, 5.5, false},
		{"Sumit", "Ghaziabad", 6667778889, 6.1, true},
		{"Sangeeta", "Ludhiana", 2223334445, 5.8, true},
		{"Vivek", "Agra", 8889990001, 5.9, true},
	}

	insertUserSQL := `INSERT INTO users (fname, city, phone, height, married) VALUES (?, ?, ?, ?, ?)`
	for _, user := range users {
		_, err := db.Exec(insertUserSQL, user.Fname, user.City, user.Phone, user.Height, user.Married)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to insert seed data: %v", err)
			return fmt.Errorf("failed to insert seed data: %v", err)
		}
	}

	logger.InfoLogger.Println("Database seeded successfully")
	return nil
}
