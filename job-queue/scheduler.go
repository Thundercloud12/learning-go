package main

import "time"


func retryScheduler(
	delayed<-chan string,
	queue chan<-string,
	shutdown<-chan struct{},
	store *Jobs,
)  {

	for{
		select{
		case <-shutdown:
			return
		case dj:=<-delayed:
			job, ok := store.Get(dj)
			if !ok {
				continue
			}

			wait := time.Until(job.NextRetryAt)
			if wait < 0 {
				wait = 0
			}

			go func (id string)  {
				timer:=time.NewTimer(time.Until(job.NextRetryAt))
				defer timer.Stop()

				select{
				case<-timer.C:
					queue<-dj
				case<-shutdown:
					return
				}
				
			}(dj)

		}
	}
	
}