package main

import (
	"bytes"
	"fmt"
	"go-wordlyest/game"
	"log"
)

func splitAfterGuess(guess game.Word, wordlist game.WordList) map[game.Ans]game.WordList {
	dm := game.Daemon{CorrectWord: guess}
	result := make(map[game.Ans]game.WordList)
	for i := range wordlist {
		word := wordlist[i]
		ans := dm.Ask(word)

		existantList, found := result[ans]
		if !found {
			existantList = []game.Word{}
		}

		result[ans] = append(existantList, word)
	}
	return result
}

type Round struct {
	guess game.Word
	ans   game.Ans
}
type History []Round

func (h History) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i := range h {
		buffer.WriteString(fmt.Sprintf("[Round %d: Guess %s Ans: %s]", i+1, h[i].guess, h[i].ans))
		if i < len(h) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

type Player func(history History) game.Word

func play(dm *game.Daemon, player Player) (int, History) {
	history := History{}
	round := 1
	for {
		guess := player(history)
		ans := dm.Ask(guess)
		history = append(history, Round{guess, ans})
		if ans.String() == "22222" {
			return round, history
		}
		round++
	}
}

func naivePlayer(list game.WordList) Player {
	i := 0
	return func(history History) game.Word {
		defer func() { i++ }()
		return list[i]
	}
}

func main() {
	wordlist := game.ReadWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	//dm := Daemon{toWord("atone")}
	//dm.playInteractively()

	//split := splitAfterGuess(wordlist[0], wordlist)
	//log.Printf("split list by %s:\n%v\n", wordlist[0], split)

	log.Println("Done!")
}
