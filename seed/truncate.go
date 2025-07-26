package main

import (
	"context"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
)

func truncate() {
	conn, err := initDB()
	if err != nil {
		pkg.Log.Error("Failed to connect to the database", err)
		return
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			pkg.Log.Error("Error closing connection: %v\n", err)
		}
	}()

	q := db.New()

	if err := q.TruncateAllTablesQuery(context.Background(), conn); err != nil {
		pkg.Log.Error("Error truncating tables: %v\n", err)
		return
	}

	pkg.Log.Info("All tables truncated successfully.")
}
