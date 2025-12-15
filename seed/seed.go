package main

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	db "github.com/Thanus-Kumaar/anokha-2025-backend/db/gen"
	pkg "github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func SeedAdmin(conn *pgx.Conn) error {
	q := db.New()

	admins, _ := q.ViewAdminSeedQuery(context.Background(), conn)
	if len(admins) > 0 {
		pkg.Log.Info("Admins already seeded, skipping...")
		return nil
	}

	hashedPassword, err := pkg.Hash("ff7bd97b1a7789ddd2775122fd6817f3173672da9f802ceec57f284325bf589f")
	if err != nil {
		pkg.Log.Error("Error inserting admin: %v\n", err)
		return err
	}

	manualAdmin := []db.SeedAdminQueryParams{

		{
			Name:     pkg.ToPgText("Naganathan M"),
			Email:    "naganathan@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Revanth Singothu"),
			Email:    "revanth@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Thanus Kumaar A"),
			Email:    "thanus@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Tharun Kumarr A"),
			Email:    "tharun@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Adithya Menon R"),
			Email:    "adukottan@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Nandgopal R Nair"),
			Email:    "nandu@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Vijay SB"),
			Email:    "vijay@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Akshay KS"),
			Email:    "akshay@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Saran Hiruthik"),
			Email:    "saran@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Keerthivasan Venkitajalam"),
			Email:    "keerthivasan@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Amrith B"),
			Email:    "amrith@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Kiran Rajeev KV"),
			Email:    "kiran@amrita.edu",
			Password: hashedPassword,
		},
		{
			Name:     pkg.ToPgText("Shivanesh"),
			Email:    "shivanesh@amrita.edu",
			Password: hashedPassword,
		},
	}

	for _, admin := range manualAdmin {
		err := q.SeedAdminQuery(context.Background(), conn, admin)
		if err != nil {
			pkg.Log.Error("Error inserting admin: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded admin.")
	return nil
}

func SeedStudents(conn *pgx.Conn) error {
	q := db.New()

	students, _ := q.ViewStudentSeedQuery(context.Background(), conn)
	if len(students) > 0 {
		pkg.Log.Info("Students already seeded, skipping...")
		return nil
	}

	hashedPassword, err := pkg.Hash("ff7bd97b1a7789ddd2775122fd6817f3173672da9f802ceec57f284325bf589f")
	if err != nil {
		pkg.Log.Error("Error inserting manual students: %v\n", err)
		return err
	}

	manualStudents := []db.SeedAmritaStudentQueryParams{
		{
			Name:        "Naganathan M",
			Email:       "naganathan@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.EN.U4CSE22240",
				Valid:  true,
			},
		},
		{
			Name:        "Revanth Singothu",
			Email:       "revanth@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.EN.U4CSE22149",
				Valid:  true,
			},
		},
		{
			Name:        "Thanus Kumaar A",
			Email:       "thanus@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.EN.U4CSE22447",
				Valid:  true,
			},
		},
		{
			Name:        "Tharun Kumarr A",
			Email:       "tharun@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.EN.U4CSE22253",
				Valid:  true,
			},
		},
		{
			Name:        "Adithya Menon R",
			Email:       "adukottan@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23506",
				Valid:  true,
			},
		},
		{
			Name:        "Nandgopal R Nair",
			Email:       "nandu@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23634",
				Valid:  true,
			},
		},
		{
			Name:        "Vijay SB",
			Email:       "vijay@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23249",
				Valid:  true,
			},
		},
		{
			Name:        "Akshay KS",
			Email:       "akshay@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23104",
				Valid:  true,
			},
		},
		{
			Name:        "Saran Hiruthik",
			Email:       "saran@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23745",
				Valid:  true,
			},
		},
		{
			Name:        "Keerthivasan Venkitajalam",
			Email:       "keerthivasan@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23037",
				Valid:  true,
			},
		},
		{
			Name:        "Amrith B",
			Email:       "amrith@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23105",
				Valid:  true,
			},
		},
		{
			Name:        "Kiran Rajeev KV",
			Email:       "kiran@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23624",
				Valid:  true,
			},
		},
		{
			Name:        "Shivanesh",
			Email:       "shivanesh@amrita.edu",
			Password:    hashedPassword,
			PhoneNumber: "9999911111",
			IsAmritaStudent: pgtype.Bool{
				Bool:  true,
				Valid: true,
			},
			AmritaRollNumber: pgtype.Text{
				String: "CB.SC.U4CSE23743",
				Valid:  true,
			},
		},
	}
	for _, student := range manualStudents {
		err := q.SeedAmritaStudentQuery(context.Background(), conn, student)
		if err != nil {
			pkg.Log.Error("Error inserting manual students: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 5 students.")
	return nil
}

func SeedOrganizers(conn *pgx.Conn) error {
	q := db.New()

	organizers, _ := q.ViewOrganizerSeedQuery(context.Background(), conn)
	if len(organizers) > 0 {
		pkg.Log.Info("Organizers already seeded, skipping...")
		return nil
	}

	hashedPassword, err := pkg.Hash("ff7bd97b1a7789ddd2775122fd6817f3173672da9f802ceec57f284325bf589f")
	if err != nil {
		pkg.Log.Error("Error inserting organizers: %v\n", err)
		return err
	}

	// Manually seed a few organizers for reference
	manualOrganizers := []db.SeedOrganizerQueryParams{
		{
			Name:        "Computer Science and Engineering",
			Email:       "cse@cb.amrita.edu",
			Password:    hashedPassword,
			OrgType:     db.OrganizerTypeEnum("DEPARTMENT"),
			StudentHead: "Tharun Kumarr A",
			FacultyHead: "Dr. Ritwik M",
		},
		{
			Name:        "Amrita Centre for Entrepreneurship",
			Email:       "ace@cb.amrita.edu",
			Password:    hashedPassword,
			OrgType:     db.OrganizerTypeEnum("CLUB"),
			StudentHead: "Thanus Kumaar A",
			FacultyHead: "Dr. Dhanya MD",
		},
	}
	for _, organizer := range manualOrganizers {
		err := q.SeedOrganizerQuery(context.Background(), conn, organizer)
		if err != nil {
			pkg.Log.Error("Error inserting manual organizer: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 5 organizers.")
	return nil
}

func SeedPeople(conn *pgx.Conn) error {
	q := db.New()

	people, _ := q.ViewPeopleSeedQuery(context.Background(), conn)
	if len(people) > 0 {
		pkg.Log.Info("People already seeded, skipping...")
		return nil
	}

	// Manually seed a few people for reference
	manualPeople := []db.SeedPeopleQueryParams{
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
		err := q.SeedPeopleQuery(context.Background(), conn, person)
		if err != nil {
			pkg.Log.Error("Error inserting manual person: %v\n", err)
			return err
		}
	}

	// Seed additional people with random data
	for i := 2; i < 20; i++ {
		person := db.SeedPeopleQueryParams{
			Name:        gofakeit.Name(),
			PhoneNumber: gofakeit.Phone(),
		}
		err := q.SeedPeopleQuery(context.Background(), conn, person)
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

	tags, _ := q.ViewTagSeedQuery(context.Background(), conn)
	if len(tags) > 0 {
		pkg.Log.Info("Tags already seeded, skipping...")
		return nil
	}

	// Manually seed a few tags for reference
	manualTags := []db.SeedTagsQueryParams{
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
		err := q.SeedTagsQuery(context.Background(), conn, tag)
		if err != nil {
			pkg.Log.Error("Error inserting manual tag: %v\n", err)
			return err
		}
	}

	// Seed additional tags with random data
	for i := 2; i < 10; i++ {
		tag := db.SeedTagsQueryParams{
			Name:         gofakeit.Word() + "Tag" + strconv.Itoa(i+1),
			Abbreviation: gofakeit.LetterN(3) + strconv.Itoa(i+1),
		}
		err := q.SeedTagsQuery(context.Background(), conn, tag)
		if err != nil {
			pkg.Log.Error("Error inserting tag: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded 10 tags.")
	return nil
}

func SeedEvents(conn *pgx.Conn) error {

	images := []string{
		"https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1451187580459-43490279c0fa?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1518770660439-4636190af475?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1555066931-4365d14bab8c?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1461749280684-dccba630e2f6?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1551288049-bebda4e38f71?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1515630278258-407f66498911?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1550745165-9bc0b252726f?q=80&w=1170&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1550751827-4bd374c3f58b?q=80&w=2070&auto=format&fit=crop",
		"https://images.unsplash.com/photo-1639153696879-e86e3a292685?q=80&w=1170&auto=format&fit=crop",
	}

	q := db.New()

	events, _ := q.ViewEventSeedQuery(context.Background(), conn)
	if len(events) > 0 {
		pkg.Log.Info("Events already seeded, skipping...")
		return nil
	}

	// Manually seed a few events for reference
	manualEvents := []db.SeedEventQueryParams{
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
			SeatsFilled:    0,
			EventStatus:    db.EventStatusEnum("ACTIVE"),
			EventMode:      db.EventModeEnum("ONLINE"),
			AttendanceMode: db.AttendanceModeEnum("SOLO"),
			CoverImageUrl:  pgtype.Text{String: images[0], Valid: true},
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
			SeatsFilled:    0,
			EventStatus:    db.EventStatusEnum("ACTIVE"),
			EventMode:      db.EventModeEnum("OFFLINE"),
			AttendanceMode: db.AttendanceModeEnum("DUO"),
			CoverImageUrl:  pgtype.Text{String: images[1], Valid: true},
		},
	}
	for _, event := range manualEvents {
		err := q.SeedEventQuery(context.Background(), conn, event)
		if err != nil {
			pkg.Log.Error("Error inserting manual event: %v\n", err)
			return err
		}
	}

	// Seed additional events with random data
	for i := 2; i < 10; i++ {
		event := db.SeedEventQueryParams{
			Name:           gofakeit.BeerName() + " " + strconv.Itoa(i+1),
			Blurb:          gofakeit.Sentence(10),
			Description:    gofakeit.Paragraph(3, 5, 10, " "),
			Price:          pgtype.Numeric{Int: big.NewInt(int64(gofakeit.Int32())), Exp: 0, Valid: true},
			IsPerHead:      gofakeit.Bool(),
			Rules:          gofakeit.Sentence(5),
			EventType:      db.EventTypeEnum(gofakeit.RandomString([]string{"EVENT", "WORKSHOP"})),
			IsGroup:        gofakeit.Bool(),
			TotalSeats:     100,
			SeatsFilled:    0,
			EventStatus:    db.EventStatusEnum(gofakeit.RandomString([]string{"CLOSED", "ACTIVE", "COMPLETED"})),
			EventMode:      db.EventModeEnum(gofakeit.RandomString([]string{"ONLINE", "OFFLINE"})),
			AttendanceMode: db.AttendanceModeEnum(gofakeit.RandomString([]string{"SOLO", "DUO"})),
			CoverImageUrl:  pgtype.Text{String: images[i], Valid: true},
		}
		err := q.SeedEventQuery(context.Background(), conn, event)
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

	records, err := q.ViewEventToOrganizerMappingSeedQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event to Organizer mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event to organizer mappings: %v", err)
	}

	events, err := q.ViewEventSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	organizers, err := q.ViewOrganizerSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing organizers: %v\n", err)
		return err
	}

	// Manually map the first two events to organizers
	manualMappings := []db.SeedEventToOrganizerMappingQueryParams{
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
		err := q.SeedEventToOrganizerMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual event to organizer mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		organizer := organizers[i%len(organizers)]
		mapping := db.SeedEventToOrganizerMappingQueryParams{
			EventID:     event.ID,
			OrganizerID: organizer.ID,
		}
		err := q.SeedEventToOrganizerMappingQuery(context.Background(), conn, mapping)
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

	records, err := q.ViewEventScheduleSeedQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event schedules already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event schedules: %v", err)
	}

	events, err := q.ViewEventSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	// Manually schedule the first two events
	manualSchedules := []db.SeedEventScheduleQueryParams{
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
		err := q.SeedEventScheduleQuery(context.Background(), conn, schedule)
		if err != nil {
			pkg.Log.Error("Error inserting manual event schedule: %v\n", err)
			return err
		}
	}

	// Seed additional schedules with random data
	for i := 2; i < len(events); i++ {
		schedule := db.SeedEventScheduleQueryParams{
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
		err := q.SeedEventScheduleQuery(context.Background(), conn, schedule)
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

	records, err := q.ViewPeopleToEventMappingSeedQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("People to Event mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list people to event mappings: %v", err)
	}

	events, err := q.ViewEventSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	people, err := q.ViewPeopleSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing people: %v\n", err)
		return err
	}

	// Manually map the first two events to people
	manualMappings := []db.SeedPeopleToEventMappingQueryParams{
		{
			EventID:  events[0].ID,
			PersonID: people[0].ID,
			EventDay: []int32{1, 2},
		},
		{
			EventID:  events[1].ID,
			PersonID: people[1].ID,
			EventDay: []int32{1},
		},
	}
	for _, mapping := range manualMappings {
		err := q.SeedPeopleToEventMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual people to event mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		person := people[i%len(people)]
		mapping := db.SeedPeopleToEventMappingQueryParams{
			EventID:  event.ID,
			PersonID: person.ID,
			EventDay: []int32{int32((i % 3) + 1)},
		}
		err := q.SeedPeopleToEventMappingQuery(context.Background(), conn, mapping)
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

	records, err := q.ViewEventTagMappingSeedQuery(context.Background(), conn)
	if len(records) > 0 {
		pkg.Log.Info("Event to Tag mappings already seeded, skipping...")
		return nil
	}
	if err != nil {
		pkg.Log.Error("failed to list event to tag mappings: %v", err)
	}

	events, err := q.ViewEventSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	tags, err := q.ViewTagSeedQuery(context.Background(), conn)
	if err != nil {
		pkg.Log.Error("Error listing tags: %v\n", err)
		return err
	}

	// Manually map the first two events to tags
	manualMappings := []db.SeedEventTagMappingQueryParams{
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
		err := q.SeedEventTagMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting manual event tag mapping: %v\n", err)
			return err
		}
	}

	// Seed additional mappings with random associations
	for i := 2; i < len(events); i++ {
		event := events[i]
		tag := tags[i%len(tags)]
		mapping := db.SeedEventTagMappingQueryParams{
			TagID:   tag.ID,
			EventID: event.ID,
		}
		err := q.SeedEventTagMappingQuery(context.Background(), conn, mapping)
		if err != nil {
			pkg.Log.Error("Error inserting event tag mapping: %v\n", err)
			return err
		}
	}

	pkg.Log.Info("Successfully seeded event tag mappings.")
	return nil
}

func SeedBookings(conn *pgx.Conn) error {
	ctx := context.Background()
	q := db.New()

	// 1. Check if bookings already exist
	bookings, _ := q.ViewBookingsSeedQuery(ctx, conn)
	if len(bookings) > 0 {
		pkg.Log.Info("Bookings already seeded, skipping...")
		return nil
	}

	// 2. Fetch Events
	events, err := q.ViewEventSeedQuery(ctx, conn)
	if err != nil {
		pkg.Log.Error("Error listing events: %v\n", err)
		return err
	}

	// 3. Fetch Students
	students, err := q.ViewStudentSeedQuery(ctx, conn)
	if err != nil {
		pkg.Log.Error("Error listing students: %v\n", err)
		return err
	}

	for _, student := range students {
		for _, event := range events {

			// to avoid a dense matrix where everyone booked everything.
			randomNumber := gofakeit.Float32Range(0, 1)

			if randomNumber < 0.3 {

				// A. Determine Status Logic
				statusRand := gofakeit.Float32Range(0, 1)
				txnStatus := "SUCCESS"
				seatsReleased := int32(0)

				if statusRand > 0.95 {
					txnStatus = "FAILED"
					seatsReleased = 1
				} else if statusRand > 0.90 {
					txnStatus = "PENDING"
				}

				// B. Prepare Parameters
				params := db.SeedBookingsQueryParams{
					TxnID:           "TXN-" + gofakeit.UUID(),
					StudentID:       student.ID,
					EventID:         event.ID,
					RegistrationFee: event.Price,
					ProductInfo:     event.Name,
					SeatsReleased:   seatsReleased,
					TxnStatus:       txnStatus,
				}

				err := q.SeedBookingsQuery(ctx, conn, params)
				if err != nil {
					pkg.Log.Error("Failed to seed booking", err)
					return err
				}
			}
		}
	}

	pkg.Log.Info(fmt.Sprintf("Successfully seeded %d bookings", len(students)))
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

	if err := SeedAdmin(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	if err := SeedStudents(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

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

	if err := SeedBookings(conn); err != nil {
		pkg.Log.Error("Seeding failed: %v\n", err)
		os.Exit(1)
	}

	pkg.Log.Info("Database seeding completed successfully.")
}
