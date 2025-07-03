package mail

var MailQueue = make(chan EmailRequest)

func StartMailWorkerPool(n int) {
	for i:=0; i<n; i++ {
		go func() {
			for req:= range MailQueue {
				if err:=SendMail(req.To, req.Subject, req.Type, req.Data); err!=nil {
					
				}
			}
		}
	}
}