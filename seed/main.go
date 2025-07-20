package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	// Initialize database connection
	err, conn := initDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background()) // Close connection when done

	// Run the seeding function
	// if err := SeedEvent(conn); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
	// 	os.Exit(1)
	// }
	if err := SeedOrganizers(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedPeople(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}

	if err := SeedTags(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}

	if err := SeedEvents(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventToOrganizerMapping(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventSchedule(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedPeopleToEventMapping(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventTagMapping(conn); err != nil {
		fmt.Fprintf(os.Stderr, "Seeding failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database seeding completed successfully.")
}
