package main

import (
	"encoding/json"
	"fmt"
	"letter-boxed-solver/internal/letterboxed"
	"letter-boxed-solver/internal/middleware"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

var PORT = 3000

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/today", handler)
	fmt.Printf("Starting server at http://localhost:%d\n", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), middleware.GzipMiddleware(mux)))
}

func handler(w http.ResponseWriter, r *http.Request) {
	maxWords := 2 // Default value
	if maxWordsParam := r.URL.Query().Get("max-words"); maxWordsParam != "" {
		parsedMaxWords, err := strconv.Atoi(maxWordsParam)
		if err == nil && parsedMaxWords > 0 {
			maxWords = parsedMaxWords
		}
	}

	// Fetch today's game data
	gameData, err := letterboxed.GetTodaysGameData()
	if err != nil {
		fmt.Println("failed to get today's game data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	box := letterboxed.NewBox(gameData.Dictionary, gameData.Sides)
	startTime := time.Now()

	// Accumulate the solutions into a slice
	var solutions [][]string
	for sol := range box.Solutions(maxWords) {
		solutions = append(solutions, sol)
	}

	fmt.Printf("Found %d solutions in %s\n", len(solutions), time.Since(startTime))

	// Sort the solutions alphabetically
	slices.SortFunc(solutions, func(a, b []string) int {
		return strings.Compare(strings.Join(a, " "), strings.Join(b, " "))
	})

	solutionsJson, err := json.Marshal(struct {
		Date      string     `json:"date"`
		Sides     []string   `json:"sides"`
		Solutions [][]string `json:"solutions"`
	}{gameData.PrintDate, gameData.Sides, solutions})
	if err != nil {
		fmt.Println("failed to serialize solutions:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(solutionsJson)
}
