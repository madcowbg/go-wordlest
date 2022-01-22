package main

import (
	"fmt"
	"go-wordlyest/game"
	"log"
	"testing"
)

func TestGreedyNeedfulPlayer(t *testing.T) {
	for _, wrd := range wordlist {
		dm := &game.Daemon{CorrectWord: wrd}
		numRounds, _ := play(wordlist, dm, GreedyNeedfulPlayer(wordlist, game.ToWord("atone")))
		log.Printf("%s\t%d\n", wrd, numRounds)
	}
}

func TestMinMaxPlayerHACK(t *testing.T) {
	var wordlist = game.ReadWordleList("data/wordle-answers-split-power.txt")

	reducedWordlist := wordlist

	for _, wrd := range reducedWordlist {
		dm := &game.Daemon{CorrectWord: wrd}
		numRounds, _ := play(reducedWordlist, dm, FastFirstHand(game.ToWord("learn"), MinMaxPlayer(false)))
		log.Printf("%s\t%d\n", wrd, numRounds)
	}
}

func TestSth(t *testing.T) {
	var wordlist = game.ReadWordleList("data/wordle-answers-split-power.txt")
	dm := &game.Daemon{CorrectWord: game.ToWord("taste")}
	numRounds, _ := play(wordlist, dm, FastFirstHand(game.ToWord("learn"), MinMaxPlayer(true)))
	fmt.Printf("Rounds: %d\n", numRounds)
}
