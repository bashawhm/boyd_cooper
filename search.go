package main

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"strings"
	"sync"
)

var searchesLock = sync.Mutex{}
var searches map[string]*Result

// Result is a struct that holds the results of a search.
type Result struct {
	re     *regexp.Regexp
	quotes []int
	index  int
}

// Next returns the next quote in the search results.
func (r *Result) Next() string {
	if r.index >= len(r.quotes) {
		return ""
	}
	ret := quoteList[r.quotes[r.index]]
	r.index++
	return ret
}

func init() {
	searchesLock.Lock()
	searches = make(map[string]*Result)
	searchesLock.Unlock()
}

// If the search is already in progress return the next quote, otherwise start a new search.
func search(s string) (string, *searchError) {
	log.Println("Searching for", s)

	searchesLock.Lock()
	defer searchesLock.Unlock()

	r, ok := searches[s]
	if ok {
		log.Println("Returning cached result")

		// If the cached results have already been depleted, restart them from the beginning
		if r.index >= len(r.quotes) {
			r.index = 0
		}

		ret := r.Next()
		return ret, nil
	}

	log.Println("Starting new search")

	err := newSearchLocked(s, 0)
	if err != nil {
		return "", err
	}

	ret := searches[s].Next()
	return ret, nil
}

// Starts a new search for a given regex.
//
// `start` is used to optimize indexing by beginning the search at a specific index in the quote list.
// Passing a `start` value of 0 will search the entire quote list.
//
// Returns an error if the regex is invalid or if there are no results
func newSearch(s string, start int) *searchError {
	searchesLock.Lock()
	err := newSearchLocked(s, start)
	searchesLock.Unlock()

	return err
}

// This is `newSearch` but **assumes you are already holding searchesLock**
func newSearchLocked(s string, start int) *searchError {
	if _, ok := searches[s]; ok {
		return nil
	}

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
	searches[s] = &Result{re, quotes, 0}
	return nil
}

// Loop through all completed searches and consider adding the new quote to them
func updateSearches(q string) {
	searchesLock.Lock()
	for _, r := range searches {
		if r.re.FindStringSubmatch(q) != nil {
			r.quotes = append(r.quotes, len(quoteList)-1)
		}
	}
	searchesLock.Unlock()
}

// Starts a new search for every word found in the quote list
func indexQuotes() {
	log.Println("Indexing quotes...")

	searchesLock.Lock()
	for i, q := range quoteList {
		// Split the quote into words
		words := regexp.MustCompile(`\W+`).Split(q, -1)

		// Start a new search for each word
		for _, w := range words {
			w = strings.ToLower(w)

			// TODO: Multithread this call?
			newSearchLocked(fmt.Sprintf(`(?i)\b%s\b`, w), i)
		}
	}
	searchesLock.Unlock()

	log.Println("Indexed", len(searches), "words")
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
