package main

import (
	"context"
	"fmt"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
)

func initDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	dbURL := os.Getenv("database_url")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set in .env file")
	}

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

func main() {
	conn, err := initDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			fmt.Printf("Error closing connection: %v\n", err)
		}
	}()

	q := db.New()

	if err := q.TruncateAllTables(context.Background(), conn); err != nil {
		fmt.Printf("Error truncating tables: %v\n", err)
		return
	}

	fmt.Println("All tables truncated successfully.")
}
