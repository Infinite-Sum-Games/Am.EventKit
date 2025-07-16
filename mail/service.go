package mail

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/joncrlsn/dque"
)

type MailerService struct {
	queue   *dque.DQue
	workers int
	ctx     context.Context
	cancel  context.CancelFunc
	wg      *sync.WaitGroup
}

func NewMailerService(path string, numWorkers int) (*MailerService, error) {
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("queue path creation failed: %w", err)
	}

	queue, err := dque.NewOrOpen("mail-queue", path, 1000, func() interface{} {
		return new(EmailRequest)
	})
	if err != nil {
		return nil, fmt.Errorf("queue init failed: %w", err)
	}
	pkg.Log.Info("mail queue created successfully!")
	ctx, cancel := context.WithCancel(context.Background())
	return &MailerService{
		queue:   queue,
		workers: numWorkers,
		ctx:     ctx,
		cancel:  cancel,
		wg:      &sync.WaitGroup{},
	}, nil
}

func (m *MailerService) Start() {
	for i := 0; i < m.workers; i++ {
		go m.worker(i)
	}
}

func (m *MailerService) Enqueue(req EmailRequest) error {
	return m.queue.Enqueue(req)
}

func (m *MailerService) worker(id int) {
	sender := NewMailer()

	for {
		select {
		case <-m.ctx.Done():
			if sender.sender != nil {
				_ = sender.sender.Close()
			}
			pkg.Log.Info(fmt.Sprintf("worker %d: shutting down", id))
			return

		default:
			item, err := m.queue.DequeueBlock()
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: Failed to dequeue", id), err)
				continue
			}
			req := item.(*EmailRequest)

			m.wg.Add(1)
			err = sender.Send(req.To, req.Subject, req.Type, req.Data)
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: Failed to send email:", id), err)
				// TODO: Retry queue or dead-letter (if critical)
			}
			m.wg.Done()
		}
	}
}

func (m *MailerService) Shutdown() {
	m.cancel()
	m.wg.Wait()
	if err := m.queue.Close(); err != nil {
		pkg.Log.Error("error in closing mail queue: ", err)
	}
}
