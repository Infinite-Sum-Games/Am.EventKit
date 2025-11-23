package mail

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/joncrlsn/dque"
)

func init() {
	gob.Register(OTPTemplateData{})
	gob.Register(WelcomeTemplateData{})
	gob.Register(RegistrationData{})
	gob.Register(EmailRequest{})
}

type MailerService struct {
	Queue   *dque.DQue
	workers int
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

var Mail *MailerService

func NewMailerService(path string, numWorkers int) (*MailerService, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("[MAIL-SERVICE]: queue path creation failed: %w", err)
	}

	queue, err := dque.NewOrOpen("mail-queue", path, 1000, func() any {
		return new(EmailRequest)
	})
	if err != nil {
		return nil, fmt.Errorf("[MAIL-SERVICE]: queue initialization failed: %w", err)
	}
	pkg.Log.Info("[MAIL-SERVICE]: mail queue created successfully!")
	ctx, cancel := context.WithCancel(context.Background())
	return &MailerService{
		Queue:   queue,
		workers: numWorkers,
		ctx:     ctx,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
	}, nil
}

func (m *MailerService) Start() {
	for i := range m.workers {
		go m.worker(i)
	}

	pkg.Log.Info(fmt.Sprintf("[OK]: Mail service initialized successfully with %d workers", m.workers))
}

func (m *MailerService) Enqueue(req *EmailRequest) error {
	return m.Queue.Enqueue(req)
}

func (m *MailerService) worker(id int) {
	sender := NewMailer()
	pkg.Log.Info(fmt.Sprintf("[MAIL-WORKER-%d]: started", id))

	for {
		select {
		case <-m.ctx.Done():
			if sender.sender != nil {
				_ = sender.sender.Close()
			}
			pkg.Log.Info(fmt.Sprintf("[MAIL-WORKER-%d]: shutting down", id))
			return

		default:
			item, err := m.Queue.DequeueBlock()
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to dequeue", id), err)
				continue
			}
			req, ok := item.(*EmailRequest)
			if !ok {
				pkg.Log.Error(fmt.Sprintf("type assertion failed for *EmailRequest, got: %#v", item), nil)
				continue
			}

			m.wg.Add(1)
			err = sender.Send(req.To, req.Subject, req.Type, req.Data)
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to send email on first attempt, retrying once...", id), err)
				// Retry once immediately
				err = sender.Send(req.To, req.Subject, req.Type, req.Data)
				if err != nil {
					pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to send email on second attempt", id), err)
					// After the second failure, the email is considered lost.
				}
			}
			m.wg.Done()
		}
	}
}

func (m *MailerService) Wait() {
	m.wg.Wait()
}

func (m *MailerService) Shutdown() {
	m.cancel()
	m.wg.Wait()
	if err := m.Queue.Close(); err != nil {
		pkg.Log.Error("[MAIL-SERVICE]: error in closing mail queue", err)
	}
}
