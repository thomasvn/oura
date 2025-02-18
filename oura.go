package oura

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var pat string

func init() {
	if pat = os.Getenv("PAT"); pat == "" {
		panic("PAT environment variable not set")
	}

	functions.HTTP("Streaks", Streaks)
	functions.HTTP("Heatmap", Heatmap)
}

// -----------------------------------------------------------------------------
// STREAKS
// -----------------------------------------------------------------------------

type streaksResult struct {
	CurrentStreak int `json:"currentStreak"`
	LongestStreak int `json:"longestStreak"`
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

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days == 0 {
		days = 365
	}
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	switch r.URL.Path {
	case "/sleep":
		json.NewEncoder(w).Encode(fetchStreaks(fmt.Sprintf("usercollection/daily_sleep?start_date=%s&end_date=%s", startDate, endDate)))
	case "/activity":
		json.NewEncoder(w).Encode(fetchStreaks(fmt.Sprintf("usercollection/daily_activity?start_date=%s&end_date=%s", startDate, endDate)))
	case "/readiness":
		json.NewEncoder(w).Encode(fetchStreaks(fmt.Sprintf("usercollection/daily_readiness?start_date=%s&end_date=%s", startDate, endDate)))
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func fetchStreaks(endpoint string) streaksResult {
	responseBytes, _ := makeOuraAPIRequest("GET", endpoint)

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

	// Longest streak
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

	return streaksResult{
		CurrentStreak: currentStreak,
		LongestStreak: longestStreak,
	}
}

// -----------------------------------------------------------------------------
// HEATMAP
// -----------------------------------------------------------------------------

type heatmapDataPoint struct {
	Date  string `json:"date"`
	Score int    `json:"score"`
}

type heatmapResult []heatmapDataPoint

func Heatmap(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	if days == 0 {
		days = 365
	}
	startDate := time.Now().AddDate(0, 0, -days).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")

	switch r.URL.Path {
	case "/sleep":
		json.NewEncoder(w).Encode(fetchHeatmap(fmt.Sprintf("usercollection/daily_sleep?start_date=%s&end_date=%s", startDate, endDate)))
	case "/activity":
		json.NewEncoder(w).Encode(fetchHeatmap(fmt.Sprintf("usercollection/daily_activity?start_date=%s&end_date=%s", startDate, endDate)))
	case "/readiness":
		json.NewEncoder(w).Encode(fetchHeatmap(fmt.Sprintf("usercollection/daily_readiness?start_date=%s&end_date=%s", startDate, endDate)))
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func fetchHeatmap(endpoint string) heatmapResult {
	responseBytes, _ := makeOuraAPIRequest("GET", endpoint)
	var response struct {
		Data []struct {
			Score     int    `json:"score"`
			TimeStamp string `json:"timestamp"`
		} `json:"data"`
	}
	json.Unmarshal(responseBytes, &response)

	// Create combined result
	result := make(heatmapResult, 0, len(response.Data))

	// Convert scores to heatmap levels (0-4) and collect dates
	for _, entry := range response.Data {
		// Map scores to heatmap levels with more granularity
		var level int
		switch {
		case entry.Score >= 95:
			level = 4 // Excellent
		case entry.Score >= 85:
			level = 3 // Good
		case entry.Score >= 75:
			level = 2 // Average
		case entry.Score >= 60:
			level = 1 // Below average
		default:
			level = 0 // Poor
		}

		result = append(result, heatmapDataPoint{
			Date:  entry.TimeStamp,
			Score: level,
		})
	}

	return result
}

// -----------------------------------------------------------------------------
// HELPERS
// -----------------------------------------------------------------------------

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
