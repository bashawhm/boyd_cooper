package main

import (
	"bufio"
	"log"
	"os"
)

var quoteList []string
var quoteFile *os.File

// Loads the quote list from a file
func loadQuotes() {
	log.Println("Loading quotes...")
	quoteList = make([]string, 0)

	scanner := bufio.NewScanner(bufio.NewReader(quoteFile))
	for scanner.Scan() {
		quoteList = append(quoteList, scanner.Text())
	}
	log.Println("Loaded", len(quoteList), "quotes")
}

// Saves a quote to the database
func writeQuote(quote string) error {
	_, err := quoteFile.WriteString(quote + "\n")
	return err
}
