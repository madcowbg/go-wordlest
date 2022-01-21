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
