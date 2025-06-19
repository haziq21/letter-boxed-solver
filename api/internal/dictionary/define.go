package dictionary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Define retrieves a summarized word definition from onelook.com. If the
// definition could not be found, it returns an empty string with a nil error.
func Define(word string) (string, error) {
	// Instead of scraping definitions from the definition pages, this function uses
	// OneLook's search suggestions API because the summarized definitions are sometimes
	// not shown on the definition pages even though they appear in search suggestions.
	reqURL := fmt.Sprintf("https://www.onelook.com/api/sug?v=ol_gte2&k=ol_home&s=%s", url.QueryEscape(word))
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// OneLook implements bot protection without these headers
	req.Header.Add("User-Agent", "LetterUnboxed/1.0.0")
	req.Header.Add("Referer", "https://onelook.com/")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %v: %w", req.URL, err)
	}
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received %s fetching %v", res.Status, req.URL)
	}

	defer res.Body.Close()
	dec := json.NewDecoder(res.Body)
	var data []struct {
		Word string   `json:"word"`
		Defs []string `json:"defs"`
	}
	if err := dec.Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	// Find the first matching definition in the response, naively ignoring case.
	// This is a loop because in some cases, when a plural word is searched and
	// OneLook doesn't have that plural word, the singular form is returned only
	// after the first suggestion. When there's an exact match for the searched
	// word, it is returned as the first suggestion.
	var definition string
	for _, d := range data {
		if !strings.EqualFold(d.Word, word) && !strings.EqualFold(d.Word+"s", word) {
			continue
		}

		if len(d.Defs) == 0 {
			return "", nil
		}

		// OneLook's generated summaries are prefixed with "u\t". The rest are
		// definitions from other sources. Also, the API seems to only ever return
		// one definition, so we assume the first one is always the only one.
		var prefixFound bool
		definition, prefixFound = strings.CutPrefix(d.Defs[0], "u\t")
		if !prefixFound {
			return "", nil
		}

		// Words that link to other words are prefixed with _
		definition = strings.ReplaceAll(definition, "_", "")
		// The definitions sometimes have trailing whitespace
		definition = strings.Trim(definition, " \n\t\r")
		break

	}

	return definition, nil
}
