package main

import (
	"context"
	"fmt"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
)

func truncate() {
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
