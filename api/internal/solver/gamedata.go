package solver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

type GameData struct {
	Dictionary []string `json:"dictionary"`
	Sides      []string `json:"sides"`
	PrintDate  string   `json:"printDate"`
}

// GetTodaysGameData fetches the Letter Boxed game data for today from the NYT website.
func GetTodaysGameData() (*GameData, error) {
	// Fetch the Letter Boxed webpage
	resp, err := http.Get("https://www.nytimes.com/puzzles/letter-boxed")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch Letter Boxed webpage")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Extract the dictionary JSON using regex
	re := regexp.MustCompile(`window\.gameData\s*=\s*(.+)`)
	matches := re.FindSubmatch(body)
	if len(matches) < 2 {
		return nil, errors.New("failed to find gameData in response body")
	}

	// Parse the JSON
	var gameData GameData
	decoder := json.NewDecoder(bytes.NewReader(matches[1]))
	if err := decoder.Decode(&gameData); err != nil {
		return nil, err
	}

	return &gameData, nil
}
