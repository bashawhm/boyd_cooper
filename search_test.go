package main

import (
	"os"
	"testing"
)

func TestSearch(t *testing.T) {
	// open the quotes file
	var err error
	quoteFile, err = os.OpenFile("quotes.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}

	// Load the quotes
	loadQuotes()

	numQuotes := len(quoteList)

	// Prefer default search
	_, err = search(".")

	if len(searches["."].quotes) != numQuotes {
		t.Error("Default search returned wrong number of quotes")
	}

	// Verify that there are no repeats in quotes
	quotes := make(map[int]struct{})
	for _, q := range searches["."].quotes {
		quotes[q] = struct{}{}
	}
	if len(quotes) != numQuotes {
		t.Error("Default search returned duplicate quotes")
	}
}
