package mail

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"
	"sync"

	"github.com/Infinite-Sum-Games/Am.EventKit/pkg"
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
		m.wg.Add(1)
		go m.worker(i)
	}

	msg := fmt.Sprintf("[OK]: Mail service initialized with %d workers", m.workers)
	pkg.Log.Info(msg)
}

func (m *MailerService) Enqueue(req *EmailRequest) error {
	return m.Queue.Enqueue(req)
}

func (m *MailerService) worker(id int) {
	defer m.wg.Done()
	sender := NewMailer()
	pkg.Log.Info(fmt.Sprintf("[MAIL-WORKER-%d]: started", id))

	for {
		select {
		case <-m.ctx.Done():
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

			err = sender.Send(req.To, req.Subject, req.Type, req.Data, req.Retries)
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to send email, re-enqueueing...", id), err)

				// Removed infinite retry to avoid blocking the worker forever on a bad email.
				// Instead, we re-enqueue the mail only for the Retries count specified in the request.
				req.Retries--
				if req.Retries > 0 {
					if err := m.Enqueue(req); err != nil {
						pkg.Log.Error("[MAILER-ERROR]: Failed to re-enqueue unsent mail", err)
					}
				}
			}
		}
	}
}

func (m *MailerService) Wait() {
	m.wg.Wait()
}

func (m *MailerService) Shutdown() {
	m.cancel()
	if err := m.Queue.Close(); err != nil {
		pkg.Log.Error("[MAIL-SERVICE]: error in closing mail queue", err)
	}
	m.wg.Wait()
}
