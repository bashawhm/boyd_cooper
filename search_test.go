package main

import (
	"regexp"
	"strings"
	"testing"
)

// TestDefaultSearch ensures that it contains all the quotes in quoteList
func TestDefaultSearch(t *testing.T) {
	randomQuotes(16)
	t.Cleanup(clean)

	// Default search
	if err := newSearch(".*", 0); err != nil {
		t.Error(err)
	}
	if len(searches[".*"].quotes) != len(quoteList) {
		t.Errorf("error, default search doesn't return all quotes")
	}
}

// TestSearch tests whether searching returns expected quotes
func TestSearch(t *testing.T) {
	setupQuotes([]string{
		"You can't tell where a program is going to spend its time.",
		"Measure. Don't tune for speed until you've measured, and even then don't unless one part of the code overwhelms the rest.",
		"Fancy algorithms are slow when n is small, and n is usually small.",
		"Fancy algorithms are buggier than simple ones, and they're much harder to implement.",
		"Data dominates. If you've chosen the right data structures and organized things well, the algorithms will almost always be self-evident.",
	})
	t.Cleanup(clean)
	indexQuotes()

	if err := newSearch("Fancy algorithms", 0); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(searches["Fancy algorithms"].quotes) != 2 {
		t.Error("error, search didn't return expected number of quotes")
	}
}

// TestCache tests whether new quotes are added to the end of pre-existing searches
func TestCache(t *testing.T) {
	randomQuotes(16)
	t.Cleanup(clean)
	indexQuotes()

	// Default search
	if err := newSearch(".*", 0); err != nil {
		t.Error(err)
	}

	newQuote := randomString(32)
	quoteList = append(quoteList, newQuote)
	updateSearches(newQuote)

	if len(searches[".*"].quotes) != 17 {
		t.Errorf("error, searches not updating with new quotes")
	}
}

// TestExhaust tests whether exhausted searches loop back to the beginning
func TestExhaust(t *testing.T) {
	setupQuotes([]string{
		"Test 1",
		"Test 2",
		"Test 3",
		"Uno reverso",
	})
	t.Cleanup(clean)

	result := Result{
		re: regexp.MustCompile("Test"),
		quotes: []int{
			0, 1, 2,
		},
		index: 0,
	}

	for i := 0; i < 16; i++ {
		if res := result.Next(); res == "" {
			t.Errorf("error, exhausted searches do not restart")
			t.FailNow()
		}
	}
}

// randomQuotes generates :num: amount of random quotes
func randomQuotes(num int) {
	// Define quotes
	var quotes []string
	for i := 0; i < num; i++ {
		quotes = append(quotes, randomString(64))
	}
	// Load quoteList
	setupQuotes(quotes)
}

// setupQuotes loads quoteList with quotes
func setupQuotes(quotes []string) {
	// Build quotes into quotefile format
	var quotestr string
	for i := 0; i < len(quotes); i++ {
		quotestr += quotes[i]
		if i+1 != len(quotes) {
			quotestr += "\n"
		}
	}
	// Parse quotefile
	reader := strings.NewReader(quotestr)
	loadQuotes(reader)
}

func clean() {
	quoteList = []string{}
	searches = make(map[string]*Result)
}
