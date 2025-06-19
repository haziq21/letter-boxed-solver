package main

import (
	"encoding/json"
	"fmt"
	"letter-unboxed/internal/dictionary"
	"letter-unboxed/internal/middleware"
	"letter-unboxed/internal/solver"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

var PORT = 3000

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/today", handler)
	mux.HandleFunc("/definition", definition)
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
	gameData, err := solver.GetTodaysGameData()
	if err != nil {
		fmt.Println("Failed to get today's game data:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	box := solver.NewBox(gameData.Dictionary, gameData.Sides)
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

	// Initialize the definitions map with empty strings
	definitions := make(map[string]string)
	for _, sol := range solutions {
		for _, word := range sol {
			definitions[word] = ""
		}
	}

	// Get the definitions for every word
	var wg sync.WaitGroup
	var mu sync.Mutex
	for word := range definitions {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			definition, err := dictionary.Define(word)
			if err != nil {
				fmt.Printf("Failed to get definition for %s: %v\n", word, err)
				delete(definitions, word)
				return
			}
			if definition == "" {
				fmt.Printf("No definition found for %s\n", word)
				delete(definitions, word)
				return
			}

			mu.Lock()
			definitions[word] = definition
			mu.Unlock()
		}(word)
	}
	wg.Wait()

	solutionsJson, err := json.Marshal(struct {
		Date        string            `json:"date"`
		Sides       []string          `json:"sides"`
		Solutions   [][]string        `json:"solutions"`
		Definitions map[string]string `json:"definitions"`
	}{gameData.PrintDate, gameData.Sides, solutions, definitions})
	if err != nil {
		fmt.Println("Failed to serialize solutions:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(solutionsJson)
}

func definition(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("word")
	if word == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error":"missing word parameter"}`))
		return
	}

	def, err := dictionary.Define(word)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if def == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Definition string `json:"definition"`
	}{def})
}
