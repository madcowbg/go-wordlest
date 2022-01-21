package main

import (
	"go-wordlyest/game"
	"log"
	"testing"
)

var wordlist = game.ReadWordleList("data/wordle-answers-alphabetical.txt")

func TestNaivePlayer(t *testing.T) {
	dm := &game.Daemon{CorrectWord: wordlist[10]}
	length, history := play(dm, naivePlayer(wordlist))

	if length != 11 {
		log.Fatalf("Naive player should have guesses it in %d but did in %d rounds by:\n%v.\n", 11, length, history)
	}
}
