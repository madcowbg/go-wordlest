package main

import (
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
	reducedWordlist := wordlist[0:100]

	for _, wrd := range reducedWordlist {
		dm := &game.Daemon{CorrectWord: wrd}
		numRounds, _ := play(reducedWordlist, dm, MinMaxPlayer(reducedWordlist))
		log.Printf("%s\t%d\n", wrd, numRounds)
	}
}