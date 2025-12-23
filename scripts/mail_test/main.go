package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Thanus-Kumaar/anokha-2025-backend/cmd"
	"github.com/Thanus-Kumaar/anokha-2025-backend/mail"
	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
)

func main() {
	// 1. Load Configuration
	fmt.Println("Loading configuration...")
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cmd.Env = config // Important: Set the global variable used by mail package

	// 2. Initialize Logger
	fmt.Println("Initializing logger...")
	logger, err := pkg.InitLogger(cmd.Env.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	pkg.Log = logger // Important: Set the global logger used by mail package

	// 3. Initialize Mail Service
	// Using a temporary queue directory for the test to avoid messing with production/dev queue if any
	queueDir := filepath.Join(os.TempDir(), "anokha_mail_test_queue")
	fmt.Printf("Initializing mail service (queue: %s)...\n", queueDir)

	// Clean up previous test runs if needed
	os.RemoveAll(queueDir)

	mailerService, err := mail.NewMailerService(queueDir, 5) // 5 workers
	if err != nil {
		log.Fatalf("Failed to create mail service: %v", err)
	}

	// 4. Start the Service
	mailerService.Start()
	defer mailerService.Shutdown()

	// 5. Enqueue 100 Emails
	targetEmail := "kakshinaruto24@gmail.com" // You might want to change this
	fmt.Printf("Enqueueing 100 emails to %s...\n", targetEmail)

	startTime := time.Now()

	for i := 1; i <= 100; i++ {
		req := &mail.EmailRequest{
			To:      []string{targetEmail},
			Subject: fmt.Sprintf("Load Test Email #%d", i),
			Type:    "welcome", // Using 'welcome' template as it's simple
			Data: mail.WelcomeTemplateData{
				UserName: fmt.Sprintf("User %d", i),
			},
		}

		if err := mailerService.Enqueue(req); err != nil {
			log.Printf("Failed to enqueue email #%d: %v", i, err)
		}
	}

	fmt.Println("All emails enqueued. Waiting for workers to finish...")

	// 6. Wait for completion
	mailerService.Wait()

	duration := time.Since(startTime)
	fmt.Printf("Done! Sent 100 emails in %v\n", duration)

	// Clean up queue dir
	os.RemoveAll(queueDir)
}
