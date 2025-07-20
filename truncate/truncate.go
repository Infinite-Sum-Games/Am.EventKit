package main

import (
	"context"
	"fmt"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
)

var q *db.Queries

func initDB() (error, *pgx.Conn) {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %v", err), nil
	}

	dbURL := os.Getenv("database_url")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL is not set in .env file"), nil
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err), nil
	}
	q = db.New() // Initialize with the connection
	return nil, conn
}

func main() {
	err, conn := initDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close(context.Background())

	if err := q.TruncateAllTables(context.Background(), conn); err != nil {
		fmt.Printf("Error truncating tables: %v\n", err)
		return
	}

	fmt.Println("All tables truncated successfully.")
}
