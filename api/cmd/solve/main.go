package main

import (
	"flag"
	"fmt"
	"letter-boxed-solver/internal/letterboxed"
	"strings"
	"time"
)

func main() {
	// Define CLI flags
	maxWords := flag.Int("max-words", 2, "maximum number of words in a solution")
	flag.Parse()

	// Fetch today's game data
	gameData, err := letterboxed.GetTodaysGameData()
	if err != nil {
		fmt.Println("failed to get today's game data:", err)
		return
	}

	box := letterboxed.NewBox(gameData.Dictionary, gameData.Sides)
	startTime := time.Now()

	// Get every solution
	count := 0
	for sol := range box.Solutions(*maxWords) {
		fmt.Println(strings.Join(sol, " "))
		count++
	}

	fmt.Printf("Found %d solutions in %s\n", count, time.Since(startTime))
}
