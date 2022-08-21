package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"sync"
)

var searchLock = sync.RWMutex{}
var searches map[string]*Result

type Result struct {
	re     *regexp.Regexp
	quotes []int
	index  int
}

func (r *Result) Next() string {
	if r.index >= len(r.quotes) {
		return ""
	}
	ret := quoteList[r.quotes[r.index]]
	r.index++
	return ret
}

func init() {
	searchLock.Lock()
	searches = make(map[string]*Result)
	searchLock.Unlock()
}

// If the search is already in progress return the next quote, otherwise start a new search.
func search(s string) (string, *searchError) {

	fmt.Println("Searching for", s)

	searchLock.Lock()
	r, ok := searches[s]
	if ok && r.index < len(r.quotes) {
		fmt.Println("Returning cached result")
		ret := r.Next()
		searchLock.Unlock()
		return ret, nil
	}
	searchLock.Unlock()

	fmt.Println("Starting new search")

	err := newSearch(s, 0)
	if err != nil {
		return "", err
	}

	searchLock.Lock()
	ret := searches[s].Next()
	searchLock.Unlock()

	return ret, nil
}

// Starts a new search for a given regex
// Pass in a starting index to start the search from a specific quote, helpful for indexing
func newSearch(s string, start int) *searchError {
	searchLock.RLock()
	if _, ok := searches[s]; ok {
		searchLock.RUnlock()
		return nil
	}
	searchLock.RUnlock()

	// WARNING: This is normally a DOS security risk.
	// However, since this bot is only accessible by our community and access is public we will trust our users.
	re, err := regexp.Compile(s)

	// If the regex is invalid, return an error
	if err != nil {
		return &searchError{compileFailed: err}
	}

	quotes := make([]int, 0)
	for i, q := range quoteList[start:] {
		if re.FindStringSubmatch(q) != nil {
			quotes = append(quotes, start+i)
		}
	}

	// If there are no results, return an error
	if len(quotes) == 0 {
		return &searchError{searchFailed: fmt.Errorf("no results found")}
	}

	// Randomize the order of the results
	rand.Shuffle(len(quotes), func(i, j int) {
		quotes[i], quotes[j] = quotes[j], quotes[i]
	})

	// Save the result
	searchLock.Lock()
	searches[s] = &Result{re, quotes, 0}
	searchLock.Unlock()

	return nil
}

// Loop through all completed searches and consider adding the new quote to them
func updateSearches(q string) {
	searchLock.Lock()

	for _, r := range searches {
		if r.re.FindStringSubmatch(q) != nil {
			r.quotes = append(r.quotes, len(quoteList)-1)
		}
	}

	searchLock.Unlock()
}

// Starts a new search for every word found in the quote list
func indexQuotes() {
	fmt.Println("Indexing quotes...")

	for i, q := range quoteList {
		// Split the quote into words
		words := regexp.MustCompile(`\W+`).Split(q, -1)

		// Start a new search for each word
		for _, w := range words {
			w = strings.ToLower(w)

			// TODO: Multithread this call
			newSearch(fmt.Sprintf(`(?i)\b%s\b`, w), i)
		}
	}

	fmt.Println("Indexed", len(searches), "words")
}

type searchError struct {
	compileFailed error
	searchFailed  error
}

func (e *searchError) Error() string {
	if e.compileFailed != nil {
		return e.compileFailed.Error()
	}
	return e.searchFailed.Error()
}
