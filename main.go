package main

import (
	"encoding/json"
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
	// HeartBeat()
	StreaksSleep()
	StreaksActivity()
	StreaksReadiness()
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
	startDate := time.Now().AddDate(0, 0, -365).Format("2006-01-02") // https://go.dev/src/time/format.go
	endDate := time.Now().Format("2006-01-02")

	responseBytes, _ := makeOuraAPIRequest("GET", fmt.Sprintf("usercollection/daily_sleep?start_date=%s&end_date=%s", startDate, endDate))

	var response struct {
		Data []struct {
			Score     int    `json:"score"`
			TimeStamp string `json:"timestamp"`
		} `json:"data"`
	}
	json.Unmarshal(responseBytes, &response)

	// Current streak
	currentStreak := 0
	for i := len(response.Data) - 1; i >= 0; i-- {
		if response.Data[i].Score >= 75 {
			currentStreak++
		} else {
			break
		}
	}
	fmt.Printf("Current streak of days with sleep score >= 75: %d\n", currentStreak)

	// Longest streak this year
	longestStreak := 0
	currentLongStreak := 0
	for i := 0; i < len(response.Data); i++ {
		if response.Data[i].Score >= 75 {
			currentLongStreak++
			if currentLongStreak > longestStreak {
				longestStreak = currentLongStreak
			}
		} else {
			currentLongStreak = 0
		}
	}
	fmt.Printf("Longest streak of days with sleep score >= 75 in the past year: %d\n", longestStreak)
}

func StreaksActivity() {
	startDate := time.Now().AddDate(0, 0, -365).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	responseBytes, _ := makeOuraAPIRequest("GET", fmt.Sprintf("usercollection/daily_activity?start_date=%s&end_date=%s", startDate, endDate))

	var response struct {
		Data []struct {
			Score         int    `json:"score"`
			TimeStamp     string `json:"timestamp"`
			Steps         int    `json:"steps"`
			TotalCalories int    `json:"total_calories"`
		} `json:"data"`
	}
	json.Unmarshal(responseBytes, &response)

	// Current streak
	currentStreak := 0
	for i := len(response.Data) - 1; i >= 0; i-- {
		if response.Data[i].Score >= 75 {
			currentStreak++
		} else {
			break
		}
	}
	fmt.Printf("Current streak of days with activity score >= 75: %d\n", currentStreak)

	// Longest streak this year
	longestStreak := 0
	currentLongStreak := 0
	for i := 0; i < len(response.Data); i++ {
		if response.Data[i].Score >= 75 {
			currentLongStreak++
			if currentLongStreak > longestStreak {
				longestStreak = currentLongStreak
			}
		} else {
			currentLongStreak = 0
		}
	}
	fmt.Printf("Longest streak of days with sleep activity >= 75 in the past year: %d\n", longestStreak)
}

func StreaksReadiness() {
	startDate := time.Now().AddDate(0, 0, -365).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	responseBytes, _ := makeOuraAPIRequest("GET", fmt.Sprintf("usercollection/daily_readiness?start_date=%s&end_date=%s", startDate, endDate))

	var response struct {
		Data []struct {
			Score     int    `json:"score"`
			TimeStamp string `json:"timestamp"`
		} `json:"data"`
	}
	json.Unmarshal(responseBytes, &response)

	// Current streak
	currentStreak := 0
	for i := len(response.Data) - 1; i >= 0; i-- {
		if response.Data[i].Score >= 75 {
			currentStreak++
		} else {
			break
		}
	}
	fmt.Printf("Current streak of days with readiness score >= 75: %d\n", currentStreak)

	// Longest streak this year
	longestStreak := 0
	currentLongStreak := 0
	for i := 0; i < len(response.Data); i++ {
		if response.Data[i].Score >= 75 {
			currentLongStreak++
			if currentLongStreak > longestStreak {
				longestStreak = currentLongStreak
			}
		} else {
			currentLongStreak = 0
		}
	}
	fmt.Printf("Longest streak of days with readiness score >= 75 in the past year: %d\n", longestStreak)
}
