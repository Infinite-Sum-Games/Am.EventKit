package main

import (
	"context"
	"fmt"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/joho/godotenv"
	"math/big"
	"os"
	"time"
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
	q = db.New()
	return nil, conn
}

func SeedOrganizers(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	organizers, _ := q.ListOrganizers(context.Background(), conn)
	if len(organizers) > 0 {
		fmt.Println("Organizers already seeded, skipping...")
		return nil
	}

	// Manually seed a few organizers for reference
	manualOrganizers := []db.InsertOrganizerParams{
		{
			Name:        "Engineering Department",
			Abbr:        "ENG",
			OrgType:     db.OrganizerTypeEnum("DEPARTMENT"),
			StudentHead: "Alice Smith",
			FacultyHead: "Dr. Emily Brown",
		},
		{
			Name:        "Music Club",
			Abbr:        "MUS",
			OrgType:     db.OrganizerTypeEnum("CLUB"),
			StudentHead: "Charlie Davis",
			FacultyHead: "Prof. John Green",
		},
	}
	for _, organizer := range manualOrganizers {
		err := q.InsertOrganizer(context.Background(), conn, organizer)
		if err != nil {
			fmt.Printf("Error inserting manual organizer: %v\n", err)
			return err
		}
	}

	// Seed additional organizers with random data
	for i := 2; i < 5; i++ {
		organizer := db.InsertOrganizerParams{
			Name:        gofakeit.Company() + " " + string(i+1),
			Abbr:        gofakeit.LetterN(3) + string(i+1),
			OrgType:     db.OrganizerTypeEnum(gofakeit.RandomString([]string{"DEPARTMENT", "CLUB"})),
			StudentHead: gofakeit.Name(),
			FacultyHead: gofakeit.Name(),
		}
		err := q.InsertOrganizer(context.Background(), conn, organizer)
		if err != nil {
			fmt.Printf("Error inserting organizer: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded 5 organizers.")
	return nil
}

func SeedPeople(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	people, _ := q.ListPeople(context.Background(), conn)
	if len(people) > 0 {
		fmt.Println("People already seeded, skipping...")
		return nil
	}

	// Manually seed a few people for reference
	manualPeople := []db.InsertPeopleParams{
		{
			Name:        "Eve Wilson",
			PhoneNumber: "9876543210",
		},
		{
			Name:        "Frank Miller",
			PhoneNumber: "8765432109",
		},
	}
	for _, person := range manualPeople {
		err := q.InsertPeople(context.Background(), conn, person)
		if err != nil {
			fmt.Printf("Error inserting manual person: %v\n", err)
			return err
		}
	}

	// Seed additional people with random data
	for i := 2; i < 20; i++ {
		person := db.InsertPeopleParams{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Phone(),
		}
		err := q.InsertPeople(context.Background(), conn, person)
		if err != nil {
			fmt.Printf("Error inserting person: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded 20 people.")
	return nil
}

func SeedTags(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	tags, _ := q.ListTags(context.Background(), conn)
	if len(tags) > 0 {
		fmt.Println("Tags already seeded, skipping...")
		return nil
	}

	// Manually seed a few tags for reference
	manualTags := []db.InsertTagsParams{
		{
			Name:        "Tech",
			Abbrevation: "TEC",
		},
		{
			Name:        "Art",
			Abbrevation: "ART",
		},
	}
	for _, tag := range manualTags {
		err := q.InsertTags(context.Background(), conn, tag)
		if err != nil {
			fmt.Printf("Error inserting manual tag: %v\n", err)
			return err
		}
	}

	// Seed additional tags with random data
	for i := 2; i < 10; i++ {
		tag := db.InsertTagsParams{
			Name:        gofakeit.Word() + "Tag" + string(i+1),
			Abbrevation: gofakeit.LetterN(3) + string(i+1),
		}
		err := q.InsertTags(context.Background(), conn, tag)
		if err != nil {
			fmt.Printf("Error inserting tag: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded 10 tags.")
	return nil
}

func SeedEvents(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	events, _ := q.ListEvents(context.Background(), conn)
	if len(events) > 0 {
		fmt.Println("Events already seeded, skipping...")
		return nil
	}

	// Manually seed a few events for reference
	manualEvents := []db.InsertEventParams{
		{
			Name:           "Tech Workshop",
			Blurb:          "A workshop on technology trends.",
			Description:    "Learn about the latest tech innovations.",
			Price:          pgtype.Numeric{Int: big.NewInt(100), Exp: 0, Valid: true},
			IsPerHead:      true,
			Rules:          "No late entries.",
			EventType:      db.EventTypeEnum("WORKSHOP"),
			IsGroup:        false,
			TotalSeats:     50,
			SeatsFilled:    25,
			EventStatus:    db.EventStatusEnum("ACTIVE"),
			EventMode:      db.EventModeEnum("ONLINE"),
			AttendanceMode: db.AttendanceModeEnum("SOLO"),
		},
		{
			Name:           "Art Exhibition",
			Blurb:          "An exhibition of student artwork.",
			Description:    "Showcase your creativity.",
			Price:          pgtype.Numeric{Int: big.NewInt(50), Exp: 0, Valid: true},
			IsPerHead:      false,
			Rules:          "Bring your own materials.",
			EventType:      db.EventTypeEnum("EVENT"),
			IsGroup:        true,
			TotalSeats:     30,
			SeatsFilled:    15,
			EventStatus:    db.EventStatusEnum("ACTIVE"),
			EventMode:      db.EventModeEnum("OFFLINE"),
			AttendanceMode: db.AttendanceModeEnum("DUO"),
		},
	}
	for _, event := range manualEvents {
		err := q.InsertEvent(context.Background(), conn, event)
		if err != nil {
			fmt.Printf("Error inserting manual event: %v\n", err)
			return err
		}
	}

	// Seed additional events with random data
	for i := 2; i < 10; i++ {
		event := db.InsertEventParams{
			Name:           gofakeit.BeerName() + " " + string(i+1),
			Blurb:          gofakeit.Sentence(10),
			Description:    gofakeit.Paragraph(3, 5, 10, " "),
			Price:          pgtype.Numeric{Int: big.NewInt(int64(gofakeit.Int32())), Exp: 0, Valid: true},
			IsPerHead:      gofakeit.Bool(),
			Rules:          gofakeit.Sentence(5),
			EventType:      db.EventTypeEnum(gofakeit.RandomString([]string{"EVENT", "WORKSHOP"})),
			IsGroup:        gofakeit.Bool(),
			TotalSeats:     gofakeit.Int32(),
			SeatsFilled:    gofakeit.Int32(),
			EventStatus:    db.EventStatusEnum(gofakeit.RandomString([]string{"CLOSED", "ACTIVE", "COMPLETED"})),
			EventMode:      db.EventModeEnum(gofakeit.RandomString([]string{"ONLINE", "OFFLINE"})),
			AttendanceMode: db.AttendanceModeEnum(gofakeit.RandomString([]string{"SOLO", "DUO"})),
		}
		err := q.InsertEvent(context.Background(), conn, event)
		if err != nil {
			fmt.Printf("Error inserting event: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded 10 events.")
	return nil
}

func SeedEventToOrganizerMapping(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	records, err := q.ListEventToOrganizerMapping(context.Background(), conn)
	if len(records) > 0 {
		fmt.Println("Event to Organizer mappings already seeded, skipping...")
		return nil
	}

	events, err := q.ListEvents(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing events: %v\n", err)
		return err
	}

	organizers, err := q.ListOrganizers(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing organizers: %v\n", err)
		return err
	}

	// Manually map the first two events to organizers
	manualMappings := []db.InsertEventToOrganizerMappingParams{
		{
			EventID:     events[0].ID,
			OrganizerID: organizers[0].ID,
		},
		{
			EventID:     events[1].ID,
			OrganizerID: organizers[1].ID,
		},
	}
	for _, mapping := range manualMappings {
		err := q.InsertEventToOrganizerMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting manual event to organizer mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		organizer := organizers[i%len(organizers)]
		mapping := db.InsertEventToOrganizerMappingParams{
			EventID:     event.ID,
			OrganizerID: organizer.ID,
		}
		err := q.InsertEventToOrganizerMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting event to organizer mapping: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded event to organizer mappings.")
	return nil
}

func SeedEventSchedule(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	records, err := q.ListEventSchedule(context.Background(), conn)
	if len(records) > 0 {
		fmt.Println("Event schedules already seeded, skipping...")
		return nil
	}

	events, err := q.ListEvents(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing events: %v\n", err)
		return err
	}

	// Manually schedule the first two events
	manualSchedules := []db.InsertEventScheduleParams{
		{
			EventID: events[0].ID,
			EventDate: pgtype.Date{
				Time:  time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour),
				Valid: true,
			},
			StartTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 10, 0, 0, 0, time.UTC),
				Valid: true,
			},
			EndTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Venue: "Online Zoom Room",
		},
		{
			EventID: events[1].ID,
			EventDate: pgtype.Date{
				Time:  time.Now().AddDate(0, 0, 2).Truncate(24 * time.Hour),
				Valid: true,
			},
			StartTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 10, 0, 0, 0, time.UTC),
				Valid: true,
			},
			EndTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Venue: "Campus Hall",
		},
	}
	for _, schedule := range manualSchedules {
		err := q.InsertEventSchedule(context.Background(), conn, schedule)
		if err != nil {
			fmt.Printf("Error inserting manual event schedule: %v\n", err)
			return err
		}
	}

	// Seed additional schedules with random data
	for i := 2; i < len(events); i++ {
		schedule := db.InsertEventScheduleParams{
			EventID: events[i].ID,
			EventDate: pgtype.Date{
				Time:  time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour),
				Valid: true,
			},
			StartTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 10, 0, 0, 0, time.UTC),
				Valid: true,
			},
			EndTime: pgtype.Timestamp{
				Time:  time.Date(2025, 7, 21, 12, 0, 0, 0, time.UTC),
				Valid: true,
			},
			Venue: gofakeit.City(),
		}
		err := q.InsertEventSchedule(context.Background(), conn, schedule)
		if err != nil {
			fmt.Printf("Error inserting event schedule: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded event schedules.")
	return nil
}

func SeedPeopleToEventMapping(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	records, err := q.ListPeopleToEventMapping(context.Background(), conn)
	if len(records) > 0 {
		fmt.Println("People to Event mappings already seeded, skipping...")
		return nil
	}

	events, err := q.ListEvents(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing events: %v\n", err)
		return err
	}

	people, err := q.ListPeople(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing people: %v\n", err)
		return err
	}

	// Manually map the first two events to people
	manualMappings := []db.InsertPeopleToEventMappingParams{
		{
			EventID:  events[0].ID,
			PersonID: people[0].ID,
		},
		{
			EventID:  events[1].ID,
			PersonID: people[1].ID,
		},
	}
	for _, mapping := range manualMappings {
		err := q.InsertPeopleToEventMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting manual people to event mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		person := people[i%len(people)]
		mapping := db.InsertPeopleToEventMappingParams{
			EventID:  event.ID,
			PersonID: person.ID,
		}
		err := q.InsertPeopleToEventMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting people to event mapping: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded people to event mappings.")
	return nil
}

func SeedEventTagMapping(conn *pgx.Conn) error {
	var q *db.Queries
	q = db.New()

	records, err := q.ListEventTagMapping(context.Background(), conn)
	if len(records) > 0 {
		fmt.Println("Event to Tag mappings already seeded, skipping...")
		return nil
	}

	events, err := q.ListEvents(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing events: %v\n", err)
		return err
	}

	tags, err := q.ListTags(context.Background(), conn)
	if err != nil {
		fmt.Printf("Error listing tags: %v\n", err)
		return err
	}

	// Manually map the first two events to tags
	manualMappings := []db.InsertEventTagMappingParams{
		{
			EventID: events[0].ID,
			TagID:   tags[0].ID,
		},
		{
			EventID: events[1].ID,
			TagID:   tags[1].ID,
		},
	}
	for _, mapping := range manualMappings {
		err := q.InsertEventTagMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting manual event tag mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		tag := tags[i%len(tags)]
		mapping := db.InsertEventTagMappingParams{
			TagID:   tag.ID,
			EventID: event.ID,
		}
		err := q.InsertEventTagMapping(context.Background(), conn, mapping)
		if err != nil {
			fmt.Printf("Error inserting event tag mapping: %v\n", err)
			return err
		}
	}

	fmt.Println("Successfully seeded event tag mappings.")
	return nil
}

func main() {
	// Initialize database connection
	err, conn := initDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Database initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

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
