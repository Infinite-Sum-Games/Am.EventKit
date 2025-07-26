package main

import (
	"context"
	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	pkg "github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"math/big"
	"os"
	"strconv"
	"time"
)

func SeedOrganizers(conn *pgx.Conn) error {
	q := db.New()

	organizers, _ := q.ListOrganizersQuery(context.Background(), conn)
	if len(organizers) > 0 {
		pkg.Log.Info("Organizers already seeded, skipping...")
		return nil
	}

	// Manually seed a few organizers for reference
	manualOrganizers := []db.InsertOrganizerQueryParams{
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
		err := q.InsertOrganizerQuery(context.Background(), conn, organizer)
		if err != nil {
			pkg.Log.Error("Error inserting manual organizer: %v\n", err)
			return err
		}
	}

	// Seed additional organizers with random data
	for i := 2; i < 5; i++ {
		organizer := db.InsertOrganizerQueryParams{
			Name:        gofakeit.Company() + " " + strconv.Itoa(i+1),
			Abbr:        gofakeit.LetterN(3) + strconv.Itoa(i+1),
			OrgType:     db.OrganizerTypeEnum(gofakeit.RandomString([]string{"DEPARTMENT", "CLUB"})),
			StudentHead: gofakeit.Name(),
			FacultyHead: gofakeit.Name(),
		}
		err := q.InsertOrganizerQuery(context.Background(), conn, organizer)
		if err != nil {
			pkg.Log.Error("Error inserting organizer: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 5 organizers.")
	return nil
}

func SeedPeople(conn *pgx.Conn) error {
	q := db.New()

	people, _ := q.ListPeopleQuery(context.Background(), conn)
	if len(people) > 0 {
		pkg.Log.Info("People already seeded, skipping...")
		return nil
	}

	// Manually seed a few people for reference
	manualPeople := []db.InsertPeopleQueryParams{
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
		err := q.InsertPeopleQuery(context.Background(), conn, person)
		if err != nil {
			pkg.Log.Error("Error inserting manual person: %v\n", err)
			return err
		}
	}

	// Seed additional people with random data
	for i := 2; i < 20; i++ {
		person := db.InsertPeopleQueryParams{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Phone(),
		}
		err := q.InsertPeopleQuery(context.Background(), conn, person)
		if err != nil {
			pkg.Log.Error("Error inserting person: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 20 people.")
	return nil
}

func SeedTags(conn *pgx.Conn) error {
	q := db.New()

	tags, _ := q.ListTagsQuery(context.Background(), conn)
	if len(tags) > 0 {
		pkg.Log.Info("Tags already seeded, skipping...")
		return nil
	}

	// Manually seed a few tags for reference
	manualTags := []db.InsertTagsQueryParams{
		{
			Name:         "Tech",
			Abbreviation: "TEC",
		},
		{
			Name:         "Art",
			Abbreviation: "ART",
		},
	}
	for _, tag := range manualTags {
		err := q.InsertTagsQuery(context.Background(), conn, tag)
		if err != nil {
			pkg.Log.Error("Error inserting manual tag: %v\n", err)
			return err
		}
	}

	// Seed additional tags with random data
	for i := 2; i < 10; i++ {
		tag := db.InsertTagsQueryParams{
			Name:         gofakeit.Word() + "Tag" + strconv.Itoa(i+1),
			Abbreviation: gofakeit.LetterN(3) + strconv.Itoa(i+1),
		}
		err := q.InsertTagsQuery(context.Background(), conn, tag)
		if err != nil {
			pkg.Log.Error("Error inserting tag: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 10 tags.")
	return nil
}

func SeedEvents(conn *pgx.Conn) error {
	q := db.New()

	events, _ := q.ListEventsQuery(context.Background(), conn)
	if len(events) > 0 {
		pkg.Log.Info("Events already seeded, skipping...")
		return nil
	}

	// Manually seed a few events for reference
	manualEvents := []db.InsertEventQueryParams{
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
		err := q.InsertEventQuery(context.Background(), conn, event)
		if err != nil {
			pkg.Log.Error("Error inserting manual event: %v\n", err)
			return err
		}
	}

	// Seed additional events with random data
	for i := 2; i < 10; i++ {
		event := db.InsertEventQueryParams{
			Name:           gofakeit.BeerName() + " " + strconv.Itoa(i+1),
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
		err := q.InsertEventQuery(context.Background(), conn, event)
		if err != nil {
			pkg.Log.Error("Error inserting event: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 10 events.")
	return nil
}

func SeedEventToOrganizerMapping(conn *pgx.Conn) error {
	q := db.New()

	records, err := q.ListEventToOrganizerMappingQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event to Organizer mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event to organizer mappings: %v", err)
	}

	events, err := q.ListEventsQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	organizers, err := q.ListOrganizersQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing organizers: %v\n", err)
		return err
	}

	// Manually map the first two events to organizers
	manualMappings := []db.InsertEventToOrganizerMappingQueryParams{
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
		err := q.InsertEventToOrganizerMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual event to organizer mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		organizer := organizers[i%len(organizers)]
		mapping := db.InsertEventToOrganizerMappingQueryParams{
			EventID:     event.ID,
			OrganizerID: organizer.ID,
		}
		err := q.InsertEventToOrganizerMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting event to organizer mapping: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded event to organizer mappings.")
	return nil
}

func SeedEventSchedule(conn *pgx.Conn) error {
	q := db.New()

	records, err := q.ListEventScheduleQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event schedules already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event schedules: %v", err)
	}

	events, err := q.ListEventsQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	// Manually schedule the first two events
	manualSchedules := []db.InsertEventScheduleQueryParams{
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
		err := q.InsertEventScheduleQuery(context.Background(), conn, schedule)
		if err != nil {
			pkg.Log.Error("Error inserting manual event schedule: %v\n", err)
			return err
		}
	}

	// Seed additional schedules with random data
	for i := 2; i < len(events); i++ {
		schedule := db.InsertEventScheduleQueryParams{
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
		err := q.InsertEventScheduleQuery(context.Background(), conn, schedule)
		if err != nil {
			pkg.Log.Error("Error inserting event schedule: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded event schedules.")
	return nil
}

func SeedPeopleToEventMapping(conn *pgx.Conn) error {
	q := db.New()

	records, err := q.ListPeopleToEventMappingQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("People to Event mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list people to event mappings: %v", err)
	}

	events, err := q.ListEventsQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	people, err := q.ListPeopleQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing people: %v\n", err)
		return err
	}

	// Manually map the first two events to people
	manualMappings := []db.InsertPeopleToEventMappingQueryParams{
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
		err := q.InsertPeopleToEventMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual people to event mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		person := people[i%len(people)]
		mapping := db.InsertPeopleToEventMappingQueryParams{
			EventID:  event.ID,
			PersonID: person.ID,
		}
		err := q.InsertPeopleToEventMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting people to event mapping: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded people to event mappings.")
	return nil
}

func SeedEventTagMapping(conn *pgx.Conn) error {
	q := db.New()

	records, err := q.ListEventTagMappingQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event to Tag mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event to tag mappings: %v", err)
	}

	events, err := q.ListEventsQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	tags, err := q.ListTagsQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing tags: %v\n", err)
		return err
	}

	// Manually map the first two events to tags
	manualMappings := []db.InsertEventTagMappingQueryParams{
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
		err := q.InsertEventTagMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual event tag mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		tag := tags[i%len(tags)]
		mapping := db.InsertEventTagMappingQueryParams{
			TagID:   tag.ID,
			EventID: event.ID,
		}
		err := q.InsertEventTagMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting event tag mapping: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded event tag mappings.")
	return nil
}

func seed() {
	conn, err := initDB()
	if err != nil {
		pkg.Log.Error("Database initialization failed: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			pkg.Log.Error("Error closing database connection: %v\n", err)
		}
	}()

	if err := SeedOrganizers(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedPeople(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	if err := SeedTags(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	if err := SeedEvents(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventToOrganizerMapping(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventSchedule(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedPeopleToEventMapping(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}
	if err := SeedEventTagMapping(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	pkg.Log.Info("Database seeding completed successfully.")
}
