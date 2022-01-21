package main

import (
	"bytes"
	"fmt"
	"go-wordlyest/game"
	"log"
	"math"
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

func NaivePlayer(list game.WordList) Player {
	i := 0
	return func(history History) game.Word {
		defer func() { i++ }()
		return list[i]
	}
}

func NaiveNeedfulPlayer(list game.WordList) Player {
	return func(history History) game.Word {
		current := list
		for _, round := range history {
			current = filter(current, round)
		}
		return current[0]
	}
}

func GreedyNeedfulPlayer(list game.WordList, firstGuess game.Word) Player {
	return func(history History) game.Word {
		if len(history) == 0 {
			return firstGuess
		}

		currentPossibilities := filterHistory(list, history)

		maxSizeAfterGuess := make(map[game.Word]int)
		for _, guess := range currentPossibilities {
			maxSizeAfterGuess[guess] = 0
			for _, possibilitiesAfterGuess := range splitAfterGuess(guess, currentPossibilities) {
				if maxSizeAfterGuess[guess] < len(possibilitiesAfterGuess) {
					maxSizeAfterGuess[guess] = len(possibilitiesAfterGuess)
				}
			}
		}

		var minMaxWrd game.Word
		var minMaxSize = math.MaxInt
		for wrd, size := range maxSizeAfterGuess {
			if size < minMaxSize {
				minMaxSize = size
				minMaxWrd = wrd
			}
		}

		return minMaxWrd
	}
}

func filterHistory(list game.WordList, history History) game.WordList {
	current := list
	for _, round := range history {
		current = filter(current, round)
	}
	return current
}

func filter(list game.WordList, round Round) game.WordList {
	result := game.WordList{}
	for _, word := range list {
		dm := game.Daemon{CorrectWord: word}
		if dm.Ask(round.guess).String() == round.ans.String() {
			result = append(result, word)
		}
	}
	return result
}

func main() {
	wordlist := game.ReadWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	//dm := Daemon{toWord("atone")}
	//dm.playInteractively()

	//split := splitAfterGuess(wordlist[0], wordlist)
	//log.Printf("split list by %s:\n%v\n", wordlist[0], split)

	dm := &game.Daemon{CorrectWord: wordlist[10]}
	length, history := play(dm, GreedyNeedfulPlayer(wordlist, game.ToWord("atone")))

	log.Printf("Naive player guesses it in %d rounds by:\n%v.\n", length, history)

	for _, wrd := range wordlist {
		dm := &game.Daemon{CorrectWord: wrd}
		numRounds, _ := play(dm, GreedyNeedfulPlayer(wordlist, game.ToWord("atone")))
		log.Printf("%s\t%d\n", wrd, numRounds)
	}

	log.Println("Done!")
}
