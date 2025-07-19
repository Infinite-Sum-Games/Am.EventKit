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
		queue:   queue,
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
			pkg.Log.Info(fmt.Sprintf("[MAIL-WORKER-%d]: shutting down", id))
			return

		default:
			item, err := m.queue.DequeueBlock()
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to dequeue", id), err)
				continue
			}
			req := item.(*EmailRequest)

			m.wg.Add(1)
			err = sender.Send(req.To, req.Subject, req.Type, req.Data)
			if err != nil {
				pkg.Log.Error(fmt.Sprintf("[MAIL-WORKER-%d]: failed to send email", id), err)
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
		pkg.Log.Error("[MAIL-SERVICE]: error in closing mail queue", err)
	}
}
