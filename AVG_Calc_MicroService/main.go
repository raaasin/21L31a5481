package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Response struct {
	Numbers         []int   `json:"numbers"`
	WindowPrevState []int   `json:"windowPrevState"`
	WindowCurrState []int   `json:"windowCurrState"`
	Avg             float64 `json:"avg"`
}

var (
	windowSize    = 10
	currentWindow []int
)

// Since authorization token is not workig, I have removed it from the code.
func fetchNumbers(numberID string) ([]int, error) {
	mockData := map[string][]int{
		"p": {2, 3, 5, 7, 11},    // Prime numbers
		"f": {1, 1, 2, 3, 5},     // Fibonaci numbers
		"e": {2, 4, 6, 8, 10},    // Even numbers
		"r": {6, 9, 12, 15, 18},  // Random numbers
	}

	numbers, exists := mockData[numberID]
	if !exists {
		return nil, fmt.Errorf("invalid number ID")
	}
	return numbers, nil
}

func updateWindow(newNumbers []int) ([]int, []int) {
	previousWindow := append([]int(nil), currentWindow...)
	for _, num := range newNumbers {
		if !contains(currentWindow, num) {
			if len(currentWindow) >= windowSize {
				currentWindow = currentWindow[1:]
			}
			currentWindow = append(currentWindow, num)
		}
	}
	return previousWindow, currentWindow
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}


func calculateAverage(numbers []int) float64 {
	if len(numbers) == 0 {
		return 0
	}
	total := 0
	for _, num := range numbers {
		total += num
	}
	return float64(total) / float64(len(numbers))
}


func handleRequest(w http.ResponseWriter, r *http.Request) {
	numberID := strings.TrimPrefix(r.URL.Path, "/numbers/")
	numberID = strings.ToLower(numberID)

	if !isValidNumberID(numberID) {
		http.Error(w, "Invalid number ID", http.StatusBadRequest)
		return
	}

	numbers, err := fetchNumbers(numberID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	prevWindow, currWindow := updateWindow(numbers)
	average := calculateAverage(currWindow)

	response := Response{
		Numbers:         numbers,
		WindowPrevState: prevWindow,
		WindowCurrState: currWindow,
		Avg:             average,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func isValidNumberID(id string) bool {
	validIDs := []string{"p", "f", "e", "r"}
	for _, validID := range validIDs {
		if id == validID {
			return true
		}
	}
	return false
}

func main() {
	http.HandleFunc("/numbers/", handleRequest)
	fmt.Println("Server running on http://localhost:9876")
	if err := http.ListenAndServe(":9876", nil); err != nil {
		fmt.Printf("Server failed: %v\n", err)
	}
}
