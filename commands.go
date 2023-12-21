package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "conspiracy",
		Description: "boydbot will tell you a conspiracy theory",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "quote",
		Description: "searches",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{{
			Name:        "word",
			Description: "The quote must contain this word",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    false,
		}},
	},
	{
		Name:        "regex",
		Description: "searches for a quote matching a regular expression",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{{
			Name:        "regex",
			Description: "please do not DOS the bot",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		}},
	},
	{
		Name:        "add",
		Description: "adds a quote to the database",
		Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{{
			Name:        "quote",
			Description: "Can not contain newlines...",
			Type:        discordgo.ApplicationCommandOptionString,
			Required:    true,
		}},
	},
	{
		Name:        "quotes",
		Description: "returns the number of quotes in the database",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "source",
		Description: "Link to my source code",
		Type:        discordgo.ChatApplicationCommand,
	},
}

var handlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"conspiracy": func(s *discordgo.Session, d *discordgo.InteractionCreate) {
		log.Println("\\conspiracy")
		sentence := buildSentence(5, 5)
		sendMessage(s, d, sentence)
	},
	"quote": func(s *discordgo.Session, d *discordgo.InteractionCreate) {

		var regex string
		if len(d.ApplicationCommandData().Options) == 0 {
			regex = "."
		} else {
			espaced := regexp.QuoteMeta(d.ApplicationCommandData().Options[0].StringValue())
			regex = fmt.Sprintf("(?i)\\b%s\\b", espaced)
		}
		log.Println("\\quote", regex)

		q, err := search(regex)
		if err != nil {
			if err.compileFailed != nil {
				message := fmt.Sprintf("You tried to use a regex to solve a problem and now you have 2 problems.\n```%s```\nhttps://regex101.com/", err.compileFailed)
				sendMessage(s, d, message)
				return
			}

			if err.searchFailed != nil {
				sendError(s, d, err.searchFailed)
				return
			}
		}

		sendMessage(s, d, quoteString(q))
	},
	"regex": func(s *discordgo.Session, d *discordgo.InteractionCreate) {
		input := d.ApplicationCommandData().Options[0].StringValue()
		log.Println("\\regex", input)

		q, err := search(input)

		if err != nil {
			if err.compileFailed != nil {
				message := fmt.Sprintf("You tried to use a regex to solve a problem and now you have 2 problems.\n```\n%s\n```", err.compileFailed)
				sendMessage(s, d, message)
				return
			}

			if err.searchFailed != nil {
				sendError(s, d, err.searchFailed)
				return
			}
		}

		sendMessage(s, d, quoteString(q))
	},
	"add": func(s *discordgo.Session, d *discordgo.InteractionCreate) {
		quote := d.ApplicationCommandData().Options[0].StringValue()

		log.Println("\add", quote)

		// The quote must not contain newlines
		if len(strings.Split(quote, "\n")) != 1 {
			sendError(s, d, fmt.Errorf("quote must not contain newlines"))
			return
		}

		err := writeQuote(quoteFile, quote)
		if err != nil {
			sendError(s, d, err)
			return
		}
		updateSearches(quote)
		quoteList = append(quoteList, quote)
		sendMessage(s, d, fmt.Sprintf("Added!\n%s", quoteString(quote)))
	},
	"quotes": func(s *discordgo.Session, d *discordgo.InteractionCreate) {
		log.Println("\\quotes")
		sendMessage(s, d, fmt.Sprintf("There are %d quotes in the database", len(quoteList)))
	},
	"source": func(s *discordgo.Session, d *discordgo.InteractionCreate) {
		log.Println("\\source")
		sendMessage(s, d, "https://github.com/bashawhm/boyd_cooper")
	},
}

// Adds a "> " to the beginning of a string
func quoteString(s string) string {
	return fmt.Sprintf("> %s", s)
}

func sendMessage(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})

	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}
}

func sendError(s *discordgo.Session, i *discordgo.InteractionCreate, e error) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: e.Error(),
		},
	})

	if err != nil {
		log.Println("Error responding to interaction: ", err)
	}
}
