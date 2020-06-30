package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func isWarm() bool {
	warm := os.Getenv("warm") == "true"
	os.Setenv("warm", "true")
	return warm
}

func runTest(sleepTime time.Duration) {
	time.Sleep(sleepTime)
}

func getSleepParameter(r *http.Request) (time.Duration, error) {
	userInput := r.URL.Query().Get("sleep")
	sleepTime, err := strconv.Atoi(userInput)
	if err != nil || sleepTime < 0 {
		return time.Nanosecond, errors.New("invalid sleep parameter")
	}
	return time.Duration(sleepTime) * time.Millisecond, nil
}

func getParameters(r *http.Request) (time.Duration, error) {
	return getSleepParameter(r)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	reused := isWarm()
	sleepTime, err := getParameters(r)

	if err != nil {
		fmt.Fprintf(w, `{"error": "%s"}`, err.Error())
		return
	}

	runTest(sleepTime)

	duration := time.Since(start).Nanoseconds()

	fmt.Fprintf(w, `{"reused": %t, "duration": %d}`, reused, duration)
}

func main() {
	log.Print("Starting server...")

	http.HandleFunc("/", handler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
