package letterboxed

import (
	"context"
	"letter-boxed-solver/internal/sets"
	"runtime"
	"strings"
	"sync"
)

// PrefixDict maps the first character of a word to all words that start with it.
type PrefixDict map[rune]sets.Set[string]

type PartialTaskResult struct {
	lastWords          []string
	potentialNextWords []string
}

// LetterBoxed holds the dictionary and prefix dictionary for fast lookups.
type LetterBoxed struct {
	sides      []string
	dict       sets.Set[string]
	prefixDict PrefixDict
}

// NewBox constructs a [LetterBoxed] from a list of dictionary words and sides of the square.
func NewBox(dict []string, sides []string) *LetterBoxed {
	allowedWords := getAllowedWords(dict, sides)
	prefixes := buildPrefixDict(allowedWords.ToSlice())

	return &LetterBoxed{
		sides:      sides,
		dict:       allowedWords,
		prefixDict: prefixes,
	}
}

// Solutions finds all sequences of words that adhere to the rules of the game.
// It returns a read-only channel and launches a worker pool to send results.
func (box *LetterBoxed) Solutions(maxWords int) <-chan []string {
	if maxWords < 1 {
		return nil
	}

	out := make(chan []string)
	ctx, cancel := context.WithCancel(context.Background())
	wordTree := NewStringTree()
	// The counter of this waitgroup represents the number of tasks waiting to be completed
	var wg sync.WaitGroup

	// Start as many workers as there are CPU cores
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for {
				// Make this better. This could be a memory leak if the context is
				// cancelled right after the select (e.g. when maxWords = 1)
				select {
				case <-ctx.Done():
					return
				default:
				}

				previousWords := wordTree.WaitToPopSequence()
				res := box.subSolution(previousWords, maxWords)

				for _, last := range res.lastWords {
					out <- append(previousWords, last)
				}
				for _, base := range res.potentialNextWords {
					wordTree.PushSequence(append(previousWords, base))
					wg.Add(1)
				}

				wg.Done()
			}
		}()
	}

	initialRes := box.subSolution([]string{}, maxWords)
	for _, last := range initialRes.lastWords {
		out <- []string{last}
	}
	for _, word := range initialRes.potentialNextWords {
		wordTree.PushSequence([]string{word})
		wg.Add(1)
	}

	go func() {
		wg.Wait()
		cancel()
		close(out)
	}()

	return out
}

func (box *LetterBoxed) subSolution(
	previousWords []string,
	maxWords int,
) (res PartialTaskResult) {
	var wordSet sets.Set[string]

	if len(previousWords) == 0 {
		wordSet = box.dict
	} else {
		// The last letter of the last word is the starting letter for the current word
		lastWord := previousWords[len(previousWords)-1]
		var startingLetter rune
		for _, ch := range lastWord {
			startingLetter = ch
		}

		// Words that were already used in the previous words won't help
		wordSet = box.prefixDict[startingLetter].Diff(sets.New(previousWords...))
	}
	if len(wordSet) == 0 {
		return
	}

	for word := range wordSet {
		newWordSeq := append(previousWords, word)

		// If there are no more unused letters, it means we've found a solution
		if box.countUnusedLetters(newWordSeq) == 0 {
			res.lastWords = append(res.lastWords, word)
		} else if len(newWordSeq) < maxWords {
			res.potentialNextWords = append(res.potentialNextWords, word)
		}
	}

	return
}

// countUnusedLetters counts the number of letters that are not used in the given words.
func (box *LetterBoxed) countUnusedLetters(words []string) int {
	unusedLetters := sets.New[rune]()

	for _, side := range box.sides {
		for _, ch := range side {
			unusedLetters.Add(ch)
		}
	}

	for _, word := range words {
		for _, ch := range word {
			unusedLetters.Remove(ch)
		}
	}

	return len(unusedLetters)
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
