package main

import (
	"log"
	"os"
	"runtime"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var quoteFile *os.File

func main() {
	log.Println("Starting BoydBot...")
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file")
	}

	// Get the bot token from the environment
	token := os.Getenv("DISCORD_TOKEN")
	if token == "" {
		log.Println("DISCORD_TOKEN missing from .env file")
		return
	}

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	// open the quotes file
	quoteFile, err = os.OpenFile("quotes.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer quoteFile.Close()

	// Load the quotes
	log.Println("Loading quotes...")
	loadQuotes(quoteFile)
	log.Println("Loaded", len(quoteList), "quotes")

	// Index quotes
	indexQuotes()
	err = newSearch(".*", 0)
	log.Println("Indexed", len(searches[".*"].quotes), "quotes")

	runtime.GC()

	// Setup discord event handlers
	ch := make(chan struct{})
	session.AddHandler(func(s *discordgo.Session, event *discordgo.Ready) {
		log.Println("Bot is ready.")
		ch <- struct{}{}
	})

	err = session.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
		return
	}

	// Handle application commands
	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := handlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}
	})

	// Wait for the bot to be ready.
	<-ch

	// Update the bot's interactions
	_, err = session.ApplicationCommandBulkOverwrite(session.State.User.ID, "586678438656475156", commands)
	if err != nil {
		log.Println("Error overwriting commands: ", err)
	}

	// Print memory usage
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Println("Memory usage:", m.Alloc/1024, "KB")

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	<-make(chan struct{})
}
