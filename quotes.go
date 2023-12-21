package main

import (
	"bufio"
	"io"
)

var quoteList []string

// Loads the quote list from an input stream
func loadQuotes(input io.Reader) {
	quoteList = make([]string, 0)

	scanner := bufio.NewScanner(bufio.NewReader(input))
	for scanner.Scan() {
		quoteList = append(quoteList, scanner.Text())
	}
}

// Writes a quote to an output stream
func writeQuote(output io.Writer, quote string) error {
	_, err := output.Write([]byte(quote + "\n"))
	return err
}
