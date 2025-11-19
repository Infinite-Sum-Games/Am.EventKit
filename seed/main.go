package main

import (
	"context"
	"flag"
	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
)

func initDB() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), cmd.Env.DatabaseURL)
	if err != nil {
		if pkg.Log != nil {
			pkg.Log.Info("Failed to connect to the database")
		} else {
			pkg.Log.Error("[CRASH] Failed to connect to database: %v\n", err)
		}
		return nil, err
	}
	return conn, nil
}

func main() {
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Printf("[CRASH] Failed to load environment variables: %v", err)
		return
	}
	cmd.Env = config
	log.Println("[OK]: Environment variables loaded successfully.")

	pkg.Log, err = pkg.InitLogger(cmd.Env.Environment)
	if err != nil {
		log.Printf("[CRASH]: Logger initialization failed: %v", err)
		return
	}
	pkg.Log.Info("[OK]: Logger initiation successful")

	seedFlag := flag.Bool("s", false, "Run seeding process")
	clearFlag := flag.Bool("c", false, "Run truncation/clear process")
	flag.Parse()

	if *clearFlag && *seedFlag {
		pkg.Log.Info("Error: Cannot run both seeding (-s) and clearing (-c) together\n,")
		os.Exit(1)
	}

	if *clearFlag {
		truncate()
	} else if *seedFlag {
		seed()
	} else {
		pkg.Log.Info("Error: Please specify either -s (seed) or -c (clear)\n")
		flag.Usage()
		os.Exit(1)
	}
}
