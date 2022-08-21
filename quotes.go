package main

import (
	"bufio"
	"fmt"
	"os"
)

var quoteList []string
var quoteFile *os.File

// Loads the quote list from a file
func loadQuotes() {
	fmt.Println("Loading quotes...")
	searchLock.Lock()
	quoteList = make([]string, 0)

	scanner := bufio.NewScanner(bufio.NewReader(quoteFile))
	for scanner.Scan() {
		quoteList = append(quoteList, scanner.Text())
	}
	searchLock.Unlock()
	fmt.Println("Loaded", len(quoteList), "quotes")
}

// Saves a quote to the database
func writeQuote(quote string) error {
	_, err := quoteFile.WriteString(quote)
	return err
}