package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Infinite-Sum-Games/Am.EventKit/cmd"
	"github.com/Infinite-Sum-Games/Am.EventKit/mail"
	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
)

func main() {
	// 1. Load Configuration
	log.Println("Loading configuration...")
	config, err := cmd.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cmd.Env = config // Important: Set the global variable used by mail package

	// 2. Initialize Logger
	log.Println("Initializing logger...")
	logger, err := pkg.InitLogger(cmd.Env.Environment)
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	pkg.Log = logger // Important: Set the global logger used by mail package

	// 3. Initialize Mail Service
	// Using a temporary queue directory for the test to avoid messing with production/dev queue if any
	queueDir := filepath.Join(os.TempDir(), "anokha_mail_test_queue")
	log.Printf("Initializing mail service (queue: %s)...\n", queueDir)

	// Clean up previous test runs if needed
	if err := os.RemoveAll(queueDir); err != nil {
		log.Printf("warning: failed to cleanup queue dir %s: %v", queueDir, err)
	}

	mailerService, err := mail.NewMailerService(queueDir, 5) // 5 workers
	if err != nil {
		log.Fatalf("Failed to create mail service: %v", err)
	}

	// 4. Start the Service
	mailerService.Start()
	defer mailerService.Shutdown()

	// 5. Enqueue 100 Emails
	targetEmail := "kakshinaruto24@gmail.com" // You might want to change this
	log.Printf("Enqueueing 100 emails to %s...\n", targetEmail)

	startTime := time.Now()

	for i := 1; i <= 100; i++ {
		req := &mail.EmailRequest{
			To:      []string{targetEmail},
			Subject: fmt.Sprintf("Load Test Email #%d", i),
			Type:    "welcome", // Using 'welcome' template as it's simple
			Data: mail.WelcomeTemplateData{
				UserName: fmt.Sprintf("User %d", i),
			},
			Retries: mail.MaxRetryCount,
		}

		if err := mailerService.Enqueue(req); err != nil {
			log.Printf("Failed to enqueue email #%d: %v", i, err)
		}
	}

	log.Println("All emails enqueued. Waiting for workers to finish...")

	// 6. Wait for completion
	mailerService.Wait()

	duration := time.Since(startTime)
	log.Printf("Done! Sent 100 emails in %v\n", duration)

	// Clean up queue dir

	if err := os.RemoveAll(queueDir); err != nil {
		log.Printf("warning: failed to cleanup queue dir %s: %v", queueDir, err)
	}
}
