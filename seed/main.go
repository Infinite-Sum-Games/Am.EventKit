package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"os"
)

func initDB() (*pgx.Conn, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("error loading .env file: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), cmd.Env.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	return conn, nil
}

func main() {
	seedFlag := flag.Bool("s", false, "Run seeding process")
	clearFlag := flag.Bool("c", false, "Run truncation/clear process")
	flag.Parse()

	if *clearFlag && *seedFlag {
		fmt.Fprintf(os.Stderr, "Error: Cannot run both seeding (-s) and clearing (-c) together\n")
		os.Exit(1)
	}

	if *clearFlag {
		truncate()
	} else if *seedFlag {
		seed()
	} else {
		fmt.Fprintf(os.Stderr, "Error: Please specify either -s (seed) or -c (clear)\n")
		flag.Usage()
		os.Exit(1)
	}
}
