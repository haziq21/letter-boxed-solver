package letterboxed

import (
	"letter-boxed-solver/pkg/sets"
	"slices"
	"strings"
	"sync"
)

// PrefixDict maps the first character of a word to all words that start with it.
type PrefixDict map[rune]sets.Set[string]

type SolveTask struct {
	previousWords []string
}

// LetterBoxed holds the dictionary and prefix dictionary for fast lookups.
type LetterBoxed struct {
	dict       sets.Set[string]
	prefixDict PrefixDict
}

// NewLetterBoxed constructs a [LetterBoxed] from a list of dictionary words and sides of the square.
func NewLetterBoxed(dict []string, sides []string) *LetterBoxed {
	allowedWords := getAllowedWords(dict, sides)
	prefixes := buildPrefixDict(allowedWords.ToSlice())

	return &LetterBoxed{
		dict:       allowedWords,
		prefixDict: prefixes,
	}
}

// Solutions finds all sequences of words that adhere to the rules of the game.
// It returns a read-only channel and launches a worker pool to send results.
func (s *LetterBoxed) Solutions(maxWords int) <-chan []string {
	out := make(chan []string)
	tasks := make(chan SolveTask, 100_000_000)
	var wg sync.WaitGroup

	workerCount := 100
	// Start all the workers
	for i := 0; i < workerCount; i++ {
		go func() {
			for t := range tasks {
				s.workSubSolutions(t.previousWords, maxWords, out, tasks, &wg)
			}
		}()
	}

	// Seed initial task
	wg.Add(1)
	tasks <- SolveTask{previousWords: nil}

	go func() {
		wg.Wait()
		close(tasks)
		close(out)
	}()

	return out
}

func (s *LetterBoxed) workSubSolutions(
	previousWords []string,
	maxWords int,
	out chan<- []string,
	tasks chan<- SolveTask,
	wg *sync.WaitGroup,
) {
	var wordSet sets.Set[string]

	if len(previousWords) == 0 {
		wordSet = s.dict
	} else {
		lastWord := previousWords[len(previousWords)-1]
		// The last letter of the last word is the starting letter for the current word
		var startingLetter rune
		for _, ch := range lastWord {
			startingLetter = ch
		}

		// Words that were already used in the previous words won't help
		wordSet = s.prefixDict[startingLetter].Diff(sets.New(previousWords...))
	}
	if len(wordSet) == 0 {
		return
	}

	sortedWords := wordSet.ToSlice()
	usedLetters := sets.New[rune]()
	for _, word := range previousWords {
		for _, ch := range word {
			usedLetters.Add(ch)
		}
	}
	unusedLetters := s.allowedLetters().Diff(usedLetters)
	s.sortWordsByLetterUsage(sortedWords, unusedLetters.ToSlice())

	for _, word := range sortedWords {
		newUnused := unusedLetters.Diff(sets.New([]rune(word)...))

		// If there are no more unused letters, it means we've found a solution
		if len(newUnused) == 0 {
			out <- append(previousWords, word)
		} else if len(previousWords)+1 < maxWords {
			wg.Add(1)
			tasks <- SolveTask{previousWords: append(previousWords, word)}
		}
	}

	wg.Done()
}

// getAllowedWords filters dictionary words that can be formed with side letters and game rules.
func getAllowedWords(dict []string, sides []string) sets.Set[string] {
	allowedLetters := sets.New([]rune(strings.Join(sides, ""))...)
	disallowedPairs := sets.New[string]()

	for _, side := range sides {
		for _, a := range side {
			for _, b := range side {
				disallowedPairs.Add(string([]rune{a, b}))
			}
		}
	}

	allowedWords := sets.New[string]()
	for _, word := range dict {
		if len(word) < 3 {
			continue
		}
		if !sets.New([]rune(word)...).IsSubsetOf(allowedLetters) {
			continue
		}

		valid := true
		for i := 0; i < len(word)-1; i++ {
			if disallowedPairs.Contains(word[i : i+2]) {
				valid = false
				break
			}
		}
		if valid {
			allowedWords.Add(word)
		}
	}
	return allowedWords
}

// sortWordsByLetterUsage sorts words by descending order of the number of distinct letters from countedLetters they use.
func (s *LetterBoxed) sortWordsByLetterUsage(words []string, countedLetters []rune) {
	letterUsages := make(map[string]int)
	for _, word := range words {
		letterUsages[word] = getLetterUsage(word, countedLetters)
	}

	slices.SortFunc(words, func(a, b string) int {
		return letterUsages[b] - letterUsages[a]
	})
}

// allowedLetters returns a set of letters that can be used to form words.
func (s *LetterBoxed) allowedLetters() sets.Set[rune] {
	allowed := sets.New[rune]()
	for prefix := range s.prefixDict {
		allowed.Add(prefix)
	}
	return allowed
}

// getLetterUsage counts distinct letters from countedLetters used in word.
func getLetterUsage(word string, countedLetters []rune) int {
	countedLettersSet := sets.New(countedLetters...)
	seen := sets.New[rune]()
	for _, ch := range word {
		if countedLettersSet.Contains(ch) {
			seen.Add(ch)
		}
	}
	return len(seen)
}

// buildPrefixDict constructs a map of starting letters to words.
func buildPrefixDict(words []string) PrefixDict {
	dict := make(PrefixDict)
	for _, word := range words {
		if len(word) == 0 {
			continue
		}

		// Get the first letter of the word as a rune
		var ch rune
		for _, c := range word {
			ch = c
			break
		}

		if _, ok := dict[ch]; !ok {
			dict[ch] = sets.New[string]()
		}
		dict[ch].Add(word)
	}
	return dict
}
