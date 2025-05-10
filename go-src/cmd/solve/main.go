package main

import (
	"encoding/json"
	"fmt"
	"letter-boxed-solver/pkg/letterboxed"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

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

	box := letterboxed.NewLetterBoxed(gameData.Dictionary, gameData.Sides)

	startTime := time.Now()

	// Accumulate the solutions into a slice
	solutions := [][]string{}
	for sol := range box.Solutions(maxWords) {
		solutions = append(solutions, sol)
	}

	fmt.Printf("Found %d solutions in %s\n", len(solutions), time.Since(startTime))

	// Sort the solutions alphabetically
	slices.SortFunc(solutions, func(a, b []string) int {
		return strings.Compare(strings.Join(a, " "), strings.Join(b, " "))
	})

	solutionsJson, err := json.Marshal(solutions)
	if err != nil {
		fmt.Println("failed to serialize solutions:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(solutionsJson)

}

func main() {
	http.HandleFunc("/todays-solutions", handler)
	fmt.Println("Starting server at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
