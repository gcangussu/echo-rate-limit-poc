package main

import (
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	id := os.Args[1]
	ratePerSecondStr := os.Args[2]
	durationStr := os.Args[3]

	ratePerSecond, err := strconv.ParseFloat(ratePerSecondStr, 64)
	if err != nil {
		panic(err)
	}

	durationSeconds, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		panic(err)
	}

	periodInSeconds := 1 / ratePerSecond
	tickerPeriod := time.Duration(math.Round(periodInSeconds * float64(time.Second)))
	duration := time.Duration(math.Round(durationSeconds * float64(time.Second)))

	ticker := time.NewTicker(tickerPeriod)
	timer := time.NewTimer(duration)

	go getApi(id)

	for {
		select {
		case <-ticker.C:
			go getApi(id)
		case <-timer.C:
			ticker.Stop()
			return
		}
	}
}

func getApi(id string) {
	before := time.Now()

	req, _ := http.NewRequest("GET", "http://localhost:1323", nil)
	req.Header.Add("id", id)
	resp, err := http.DefaultClient.Do(req)

	after := time.Now()

	if err != nil {
		panic(err)
	}

	fmt.Printf(
		`{"status":%v,"latency":%v}`+"\n",
		resp.StatusCode,
		after.Sub(before).Nanoseconds(),
	)
}
