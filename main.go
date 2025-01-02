package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var pat string

func init() {
	pat = os.Getenv("PAT")
	if pat == "" {
		panic("PAT environment variable not set")
	}
}

func main() {
	HeartBeat()
	StreaksSleep()
}

func makeOuraAPIRequest(method, endpoint string) ([]byte, error) {
	client := &http.Client{}
	url := fmt.Sprintf("https://api.ouraring.com/v2/%s", endpoint)
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+pat)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}

	return body, nil
}

func HeartBeat() {
	response, _ := makeOuraAPIRequest("GET", "usercollection/personal_info")
	fmt.Println(string(response))
}

func StreaksSleep() {
	startDate := time.Now().AddDate(0, 0, -90).Format("2006-01-02") // https://go.dev/src/time/format.go
	endDate := time.Now().Format("2006-01-02")

	fmt.Println(startDate, endDate)

	response, _ := makeOuraAPIRequest("GET", fmt.Sprintf("usercollection/daily_sleep?start_date=%s&end_date=%s", startDate, endDate))
	fmt.Println(string(response))
}
