package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func main() {

	threadCount := flag.Int("threads", 1, "Thread Count")
	rampUpTime := flag.Int("rampup", 60, "Ramp up time in seconds")
	url := flag.String("url", "https://www.google.com", "Destination URL")
	flag.Parse()

	log.Printf("spawning http listener")
	http.HandleFunc("/status", httpStatus)
	http.HandleFunc("/", httpHandler)
	go func() {
		http.ListenAndServe(":80", nil)
	}()

	waitTime := *rampUpTime / *threadCount

	log.Printf("Execution will happen with %d threads with a ramp up time of %d seconds\n", *threadCount, *rampUpTime)

	tchan := make(chan int)
	go func(c chan<- int) {
		for ti := 1; ti <= *threadCount; ti++ {
			c <- ti
			time.Sleep(time.Duration(waitTime) * time.Second)
		}
	}(tchan)

	for {
		select {
		case ts := <-tchan:
			log.Printf("Thread #%d started", ts)
			go func(t int) {
				for {
					client := http.Client{Timeout: 1 * time.Second}

					log.Printf("HTTP request inflight to %s", *url)
					r, err := client.Get(*url)
					if err != nil {
						log.Printf("ERROR: %s", err)
					} else {
						log.Printf("Succesfully received data from: %s, with HTTP: %d", *url, r.StatusCode)
					}
					time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
				}
			}(ts)
		}
	}
}

func httpStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "order-processor")
}
