package main

import (
	"log"
	"net/http"
	"sync"
	
)

func main() {
	myApp := NewApp(10, 20, nil) // Add nil for customHandlers

	shutdown := make(chan struct{})
	 // For stress test
	var wg sync.WaitGroup

	go retryScheduler(myApp.delayed, myApp.queue, shutdown, myApp.store)

	// Workers
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go myApp.worker(i+1,&wg, shutdown)
	}

	
	

	// HTTP
	registerRoutes(myApp.queue, myApp.store)

	log.Println("HTTP server running on :8080")
	go http.ListenAndServe(":8080", nil)

	
	select {}   // block forever
	 // Close stop on exit (though blocked)
}
