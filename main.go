package main

import (
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/buger/goterm"
)

var (
	url                string
	expectedStatusCode int
	waitTime           time.Duration
	numBodyBytes       int
)

func init() {
	flag.StringVar(&url, "url", "", "url to test")
	flag.IntVar(&expectedStatusCode, "expected-status-code", 200, "expected http status code")
	flag.DurationVar(&waitTime, "wait-time", 2*time.Second, "time to wait between requests")
	flag.IntVar(&numBodyBytes, "num-body-bytes", 1000, "number of bytes to show from body")
}

func main() {
	flag.Parse()

	if url == "" {
		log.Fatalln("Please specify a URL")
	}

	numRequests := 0
	numSuccessful := 0
	numFailed := 0

	start := time.Now()

	for {
		numRequests++

		resp, err := http.Get(url)
		if err != nil {
			log.Fatalln("An error occurred getting the URL:", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode == expectedStatusCode {
			numSuccessful++
		} else {
			numFailed++
		}

		testDuration := time.Since(start)
		uptime := (float64(numSuccessful) / float64(numRequests)) * float64(100)
		buf := make([]byte, numBodyBytes)
		_, err = resp.Body.Read(buf)
		if err != nil {
			log.Fatalln("An error occurred reading the body:", err)
		}

		goterm.Clear()
		goterm.MoveCursor(1, 1)
		goterm.Println("URL:", url)
		goterm.Println("Expected status code:", expectedStatusCode)
		goterm.Println("Wait time:", waitTime)
		goterm.Println("Requests:", numRequests)
		goterm.Println("Successful:", numSuccessful)
		goterm.Println("Failed:", numFailed)
		goterm.Println("Test duration:", testDuration)
		goterm.Printf("Uptime: %v%%\n", uptime)
		goterm.Println("Last status code:", resp.StatusCode)
		goterm.Printf("Last response body (%v bytes): \n%v", numBodyBytes, string(buf))
		goterm.Flush()

		time.Sleep(waitTime)
	}
}
