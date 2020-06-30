package kubeless

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kubeless/kubeless/pkg/functions"
)

func isWarm() bool {
	warm := os.Getenv("warm") == "true"
	os.Setenv("warm", "true")
	return warm
}

func runTest(sleepTime time.Duration) {
	time.Sleep(sleepTime)
}

func getSleepParameter(event functions.Event) (time.Duration, error) {
	userInput := event.Extensions.Request.URL.Query().Get("sleep")
	sleepTime, err := strconv.Atoi(userInput)
	if err != nil || sleepTime < 0 {
		return time.Nanosecond, errors.New("invalid sleep parameter")
	}
	return time.Duration(sleepTime) * time.Millisecond, nil
}

func getParameters(event functions.Event) (time.Duration, error) {
	return getSleepParameter(event)
}

func Sleep(event functions.Event, context functions.Context) (string, error) {
	start := time.Now()
	reused := isWarm()
	sleepTime, err := getParameters(event)

	if err != nil {
		return fmt.Sprintf(`{"error": "%s"}`, err.Error()), nil
	}

	runTest(sleepTime)
	duration := time.Since(start).Nanoseconds()

	return fmt.Sprintf(`{"reused": %t, "duration": %d}`, reused, duration), nil
}
