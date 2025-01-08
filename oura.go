package oura

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var pat string

type streaksResult struct {
	CurrentStreak int `json:"currentStreak"`
	LongestStreak int `json:"longestStreak"`
}

func init() {
	if pat = os.Getenv("PAT"); pat == "" {
		panic("PAT environment variable not set")
	}

	functions.HTTP("Streaks", Streaks)
}

func Streaks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/sleep":
		json.NewEncoder(w).Encode(streaksSleep())
	case "/activity":
		json.NewEncoder(w).Encode(streaksActivity())
	case "/readiness":
		json.NewEncoder(w).Encode(streaksReadiness())
	case "/heartbeat":
		json.NewEncoder(w).Encode(json.RawMessage(heartBeat()))
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func heartBeat() []byte {
	response, _ := makeOuraAPIRequest("GET", "usercollection/personal_info")
	return response
}

func streaksSleep() streaksResult {
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

	return streaksResult{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}
}

func streaksActivity() streaksResult {
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

	return streaksResult{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}
}

func streaksReadiness() streaksResult {
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

	return streaksResult{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}
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
