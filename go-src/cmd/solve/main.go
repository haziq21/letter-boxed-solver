package main

import (
	"bufio"
	"fmt"
	"letter-boxed-solver/pkg/letterboxed"
	"os"
	"strings"
)

func main() {
	dict, err := readWordsFromFile("../dictionary.txt")
	if err != nil {
		fmt.Println("Error reading dictionary:", err)
		return
	}

	letterboxed := letterboxed.NewLetterBoxed(dict, []string{"rxn", "aof", "htc", "epi"})
	numSols := 0

	for sol := range letterboxed.Solutions(3) {
		fmt.Println(strings.Join(sol, " -> "))
		numSols++
	}

	fmt.Printf("Found %d solutions\n", numSols)
}

func readWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word != "" {
			words = append(words, word)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return words, nil
}
