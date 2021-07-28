package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var Token string = ""

var quoteList []string
var quoteIndex map[string][]string
var quotesFile *os.File

var randGen *rand.Rand
var lastSearches map[string][]string
var searchOrder []string

const MAX_SEARCHES int = 5

func buildSentance(asideChance uint32, interjectionChance uint32) string {
	subjects := [40]string{"those little green cops", "the Milkman", "the military industrial complex", "the suits", "the analyticals, man,", "those Bermuda Triangle sharks", "all them haters", "Hernando", "that little fat kid, with the bunny,", "the doctors back at the clinic", "the pelicans", "the squirrels", "the manager of that boy band", "those eggheads in their ivory tower", "that guy with the eyepatch", "the Psycho-whatsits", "the freaky hunchback girl who loves brains so much", "the dairy industry", "the kid with the goggles", "the dogtrack regulators", "the tuna canneries", "the National Park system", "Big Oil", "organized labor", "the rodeo clown cartel", "the media", "the cows", "foreign toymakers", "the dairy industry", "the intelligentsia", "the fluoride producers", "a secret doomsday cult", "the president's brother", "my first cat, Seymour,", "oh! one of my nostril hairs", "the intelligence community", "the five richest families in the country", "all those stupid crows", "some sort of power, y'know?", "my good pal Vinny"}
	subjectConnector := [8]string{"and", "...or else maybe...", "...no, no, wait, I mean...", "in conjunction with", "with the full blessing of", "with the backing of", "who are merely the pawns of", "who are the puppet masters of"}
	transitiveVerb := [12]string{"went to the prom with", "ate a whole jar of olives with", "are working for", "are telling my location to", "made a deal, back in '68, with", "sold their soul to", "are controlled by", "bought votes to protect", "are doing the dirty work of", "got in bed with", "signed a secret treaty with", "has been officially linked with"}
	intransitiveVerb := [17]string{"know the truth", "won't stop visiting me", "keep sparring with me", "have been spitting on me all day", "do this horrible thing, but in conjunction with who? Or, whom?", "are crawling all over my skin", "bit me all night, so I couldn't sleep", "have everyone fooled", "were digging away at the plastics", "were dialing in through the optics", "stole my theories and reprinted them—incorrectly—to discredit them", "are not to be trusted", "have been living off the teat of the dairy industry", "have been fixing oil prices", "assassinated the one man in their way", "pretty much control everything", "pick who lives, and who dies, and what the football scores are going to be every week"}
	//    verbConnector := [7]string{"and they obviously", "I know they", "but they can't hide that they", "ha! Like I don't know that they", "and let's just say for now that they", "if I know anything, I know that they", "and sure as the nose on my face, I am sure they"}
	preposition := [7]string{"to get", "because they want", "in order to monopolize", "to keep down", "so the people never find out about", "and who wins? Them. Who loses?", "all in a big fight over"}
	object := [17]string{"the truth", "all of us", "the whole sack of lies", "the innocents", "the biggest conspiracy of all", "the infrastructure", "the lap belt man", "the water supply", "the rotundra", "the AM Tenderizer", "last specimen of the supervirus", "the witnesses", "my hooch", "the hanging udders", "a clean-burning perpetual energy source", "a religious artifact with supposedly unimaginable powers", "exactly what, nobody knows"}
	conclution := [9]string{"How long do they think they can hide that?", "I mean, who do they think they're fooling?", "Can I really be the only person who sees this?", "Someone has to get this information to the people.", "If they find out I know this stuff, I'm dead.", "Oh man, this stuff is hot.", "since the year \"dot\".", "right under peoples noses!", "and nobody seems to care!"}
	aside := [2]string{"Visiting hours are over!", "Why does that hydrant keep looking at me?"}
	interjection := [15]string{"*chuckles*", "(Ho ho!)", "(Wait...)", "(Uh...)", "(Um...)", "*cough*", "(Uh...)", "(Hmm...)", "(Ha!)", "(Yeah, yeah, yeah...)", "(What?)", "(No, no, nonono...)", "(Okay, okay but...)", "(Huh?)", "(Oh-hoh, RIGHT...)"}

	sentance := ""
	if randGen.Uint32()%asideChance == 0 {
		sentance = aside[rand.Uint32()%2]
		return sentance
	}
	if randGen.Uint32()%interjectionChance == 0 {
		sentance = interjection[rand.Uint32()%15]
		return sentance
	}

	if randGen.Int()%2 == 0 {
		sentance = subjects[rand.Int()%40] + " " + transitiveVerb[rand.Int()%12] + " " + object[rand.Int()%17] + " " + conclution[rand.Int()%9]
	} else {
		sentance = subjects[rand.Uint32()%40] + " " + subjectConnector[rand.Uint32()%8] + " " + subjects[rand.Uint32()%39] + " " + intransitiveVerb[rand.Uint32()%17] + " " + preposition[rand.Uint32()%7] + " " + object[rand.Uint32()%17]
	}
	return sentance
}

var punct = regexp.MustCompile("[[:punct:]]")

func buildIndex() {
	fmt.Println("Building quote index...")
	quoteIndex = make(map[string][]string)

	for _, quote := range quoteList {
		addToIndex(quote)
	}

	fmt.Println("Built index!")
}

func addToIndex(quote string) {
	// Removes all punctuation from string to make a more human friendly index
	filtered := punct.ReplaceAllString(quote, "")

	for _, word := range strings.Split(filtered, " ") {
		quoteIndex[word] = append(quoteIndex[word], quote)
	}
}

func filter(array []string, f func(string) bool) []string {
	filteredArray := make([]string, 0)
	for _, str := range array {
		if f(str) {
			filteredArray = append(filteredArray, str)
		}
	}
	return filteredArray
}

func getSearchQuote(search string) string {
	if lastSearches == nil {
		lastSearches = make(map[string][]string)
	}

	sl, ok := lastSearches[search]
	if ok && len(sl) > 0 {
		fmt.Printf("Search found in lS, %d remaining\n", len(sl))
		ret := sl[0]
		if len(sl) == 1 {
			fmt.Println("lS emptied of search")
			delete(lastSearches, search)
			searchOrder = filter(searchOrder, func(s string) bool {
				return s != search
			})
		} else {
			lastSearches[search] = sl[1:]
			fmt.Printf("...now %d\n", len(sl)-1)
		}
		return ret
	}

	if len(quoteList) == 0 {
		return "No quotes found..."
	}

	re, err := regexp.Compile(search)
	if err != nil {
		return "Error compiling pattern: " + err.Error()
	}

	// get the first word in search so we can look it up in the index to speed up search time
	var word = search
	for i, v := range search {
		if v == ' ' {
			word = search[0:i]
			break
		}
	}

	filteredQuotes := quoteIndex[word]

	if len(filteredQuotes) == 0 {
		// if we didn't find anything in the index we might still have a regex
		// in this case search the whole database
		fmt.Printf("\"%s\" not found in index, filtering all %d quotes\n", word, len(quoteList))
		filteredQuotes = filter(quoteList, func(str string) bool {
			return re.FindStringSubmatch(str) != nil
		})
	} else {
		// if the "first word" is not the entire search string then the index probably
		// included results we were not looking for and we have to filter these off
		if word != search {
			fmt.Printf("\"%s\" found in index, filtering %d quotes\n", word, len(filteredQuotes))
			filteredQuotes = filter(filteredQuotes, func(str string) bool {
				return re.FindStringSubmatch(str) != nil
			})
		} else {
			fmt.Printf("\"%s\" found in index, %d perfect matches\n", word, len(filteredQuotes))
		}
	}

	fmt.Printf("Fresh search made %d matches\n", len(filteredQuotes))

	if len(filteredQuotes) == 0 {
		return "No quotes found with that query..."
	}

	shuffled := make([]string, len(filteredQuotes))
	for idx, perm := range randGen.Perm(len(filteredQuotes)) {
		shuffled[perm] = filteredQuotes[idx]
	}

	ret := shuffled[0]
	shuffled = shuffled[1:]
	if len(shuffled) > 0 {
		fmt.Printf("Reshuffling %d matches into lS\n", len(shuffled))
		lastSearches[search] = shuffled
		searchOrder = append(searchOrder, search)
		if len(searchOrder) > MAX_SEARCHES {
			wm := len(searchOrder) - MAX_SEARCHES
			for _, v := range searchOrder[:wm] {
				delete(lastSearches, v)
			}
			searchOrder = searchOrder[wm:]
		}
		fmt.Printf("Current sO is %#v\n", searchOrder)
	}

	return ret
}

func loadQuotes(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		return
	}
	fin := bufio.NewScanner(bufio.NewReader(file))
	fin.Split(bufio.ScanLines)
	for fin.Scan() {
		quoteList = append(quoteList, fin.Text())
	}
	file.Close()
}

func writeQuote(quote string) {
	quotesFile.WriteString(quote + "\n")
}

func stripPrefix(prefix, data string) string {
	var res string
	for i := 0; i < len(data); i++ {
		if strings.HasPrefix(data[:i], prefix) {
			res = data[i:]
			break
		}
	}
	return res
}

func bot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "Boyd") {
		s.ChannelMessageSend(m.ChannelID, buildSentance(5, 5))
	}
	if strings.HasPrefix(m.Content, "!quoteadd ") {
		res := stripPrefix("!quoteadd ", m.Content)
		fmt.Println("Adding quote: " + res)
		writeQuote(res)
		addToIndex(res)
		quoteList = append(quoteList, res)
		s.ChannelMessageSend(m.ChannelID, "Added!")
	} else if strings.HasPrefix(m.Content, "!quotes") {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There are %d quotes in my database", len(quoteList)))
	} else if strings.HasPrefix(m.Content, "!quote") {
		res := stripPrefix("!quote", m.Content)
		searchMsg := stripPrefix(" ", res)
		ret := getSearchQuote(searchMsg)
		s.ChannelMessageSend(m.ChannelID, ret)
	} else if strings.HasPrefix(m.Content, "!help") {
		s.ChannelMessageSend(m.ChannelID, "`!quote` shows a random quote\n`!quote PATTERN` finds a quote that has a substring matching PATTERN\n`!quoteadd` adds a new quote to the database\n`!help` shows this help message")
	}
}

func main() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	ds, err := discordgo.New("Bot " + Token)
	if err != nil {
		panic("failed to create session")
	}

	loadQuotes("./quotes.txt")
	buildIndex()

	quotesFile, err = os.OpenFile("./quotes.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer quotesFile.Close()

	ds.AddHandler(bot)
	err = ds.Open()
	if err != nil {
		panic("failed to connect")
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Boyd is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	ds.Close()

}
