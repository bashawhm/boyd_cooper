package main

import (
	"math/rand"
	"time"
)

var randGen *rand.Rand

func init() {
	randGen = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func buildSentence(asideChance uint32, interjectionChance uint32) string {
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

	if randGen.Uint32()%asideChance == 0 {
		return aside[rand.Uint32()%2]
	} else if randGen.Uint32()%interjectionChance == 0 {
		return interjection[rand.Uint32()%15]
	} else if randGen.Int()%2 == 0 {
		return subjects[rand.Int()%40] + " " + transitiveVerb[rand.Int()%12] + " " + object[rand.Int()%17] + " " + conclution[rand.Int()%9]
	} else {
		return subjects[rand.Uint32()%40] + " " + subjectConnector[rand.Uint32()%8] + " " + subjects[rand.Uint32()%39] + " " + intransitiveVerb[rand.Uint32()%17] + " " + preposition[rand.Uint32()%7] + " " + object[rand.Uint32()%17]
	}
}
