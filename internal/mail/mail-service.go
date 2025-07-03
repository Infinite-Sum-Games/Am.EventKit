package mail

import (
	"log"
)

// we can use this mail queue and the struct EmailRequest to send request to mail service
var MailQueue = make(chan EmailRequest)

func StartMailWorkerPool(n int) {
	for range n {
		go func() {
			for req:= range MailQueue {
				if err:=SendMail(req.To, req.Subject, req.Type, req.Data); err!=nil {
					//TODO: change this to logger
					log.Println("cannot send mail, error occured: %w", err)
				}
			}
		}()
	}
}