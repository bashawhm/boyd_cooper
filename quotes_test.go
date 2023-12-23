package main

import (
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestLoadQuotes(t *testing.T) {
	log.Println(quoteList)

	// Define quotes
	quotes := map[int]string{
		0: "Be yourself; everyone else is already taken. - Oscar Wilde",
		1: "In the end, it's not the years in your life that count. It's the life in your years. - Abraham Lincoln",
		2: "Who is a greater reprobate than he, Who feels compassion at the doom divine? - Virgil",
		3: "Opinion is the medium between knowledge and ignorance. - Plato",
	}

	// Build quotes into quotefile format
	var quotestr string
	for i := 0; i < len(quotes); i++ {
		quotestr += quotes[i]
		if i+1 != len(quotes) {
			quotestr += "\n"
		}
	}

	// Parse quotefile
	r := strings.NewReader(quotestr)
	loadQuotes(r)

	// Compare parsed values against originals
	for key, val := range quotes {
		if quoteList[key] != val {
			t.Errorf(`error, expected "%s", got "%s"`, val, quoteList[key])
		}
	}
}

// writer satisfies io.Writer, stores whatever is written to it
type writer struct {
	Storage []byte
}

func (w *writer) Write(p []byte) (n int, err error) {
	w.Storage = p
	return len(p), nil
}

func TestWriteQuote(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 25; i++ {
		// Generate random quote
		quote := randomString(r.Intn(128))

		// Write quote to the fake quotefile
		qf := writer{}
		writeQuote(&qf, quote)

		// Compare the written quote to the one we sent
		if string(qf.Storage) != quote+"\n" {
			t.Errorf(`error, expected "%s", got "%s"`, quote, qf.Storage)
		}
	}
}

func randomString(n int) string {
	runes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890!@#$%^&*()_+")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
