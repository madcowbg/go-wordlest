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
	Guess game.Word
	Ans   game.Ans
}
type History []Round

func (h History) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i := range h {
		buffer.WriteString(fmt.Sprintf("[Round %d: Guess %s Ans: %s]", i+1, h[i].Guess, h[i].Ans))
		if i < len(h) {
			buffer.WriteString(", ")
		}
	}
	buffer.WriteString("}")
	return buffer.String()
}

type Player func(history History, allowed game.WordList) game.Word

func play(allowedWords game.WordList, dm *game.Daemon, player Player) (int, History) {
	history := History{}
	currentlyAllowed := allowedWords
	roundNumber := 1
	for {
		guess := player(history, currentlyAllowed)
		ans := dm.Ask(guess)
		round := Round{guess, ans}
		history = append(history, round)
		if ans.String() == "22222" {
			return roundNumber, history
		}
		currentlyAllowed = filter(currentlyAllowed, round)
		roundNumber++
	}
}

func NaivePlayer(list game.WordList) Player {
	i := 0
	return func(history History, _ game.WordList) game.Word {
		defer func() { i++ }()
		return list[i]
	}
}

func NaiveNeedfulPlayer(_ History, currentlyAllowed game.WordList) game.Word {
	return currentlyAllowed[0]
}

func GreedyNeedfulPlayer(_ game.WordList, firstGuess game.Word) Player {
	return func(history History, currentlyAllowed game.WordList) game.Word {
		if len(history) == 0 {
			return firstGuess
		}

		maxSizeAfterGuess := make(map[game.Word]int)
		for _, guess := range currentlyAllowed {
			maxSizeAfterGuess[guess] = 0
			for _, possibilitiesAfterGuess := range splitAfterGuess(guess, currentlyAllowed) {
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

func filter(list game.WordList, round Round) game.WordList {
	result := game.WordList{}
	for _, word := range list {
		dm := game.Daemon{CorrectWord: word}
		if dm.Ask(round.Guess).Equals(round.Ans) {
			result = append(result, word)
		}
	}
	return result
}

func MinMaxPlayer(_ game.WordList) Player {
	var calculator func(history History, currentlyAllowed game.WordList, currentBestDepth int) (game.Word, int)
	calculator = func(history History, currentlyAllowed game.WordList, currentBestDepth int) (game.Word, int) {
		if len(currentlyAllowed) == 1 {
			// base case - we are here!
			return currentlyAllowed[0], len(history) + 1
		}

		if len(history) >= currentBestDepth {
			return currentlyAllowed[0], math.MaxInt // prune case - we can't beat it, so may as well skip it and declare undesirable
		}

		var minWorstCaseGuess game.Word
		var minWorstCaseDepth = math.MaxInt
		for _, guess := range currentlyAllowed {
			afterGuess := splitAfterGuess(guess, currentlyAllowed)
			maxDepth := 0
			for ans, afterGuessAllowed := range afterGuess {
				_, depthInThatHistory := calculator(
					append(history, Round{Guess: guess, Ans: ans}),
					afterGuessAllowed,
					currentBestDepth)
				if maxDepth < depthInThatHistory {
					maxDepth = depthInThatHistory
				}
				if maxDepth > currentBestDepth {
					break // prune when any other ans will pull us further into undesirable paths...
				}
			}
			if minWorstCaseDepth > maxDepth {
				minWorstCaseDepth = maxDepth
				minWorstCaseGuess = guess
			}
			if currentBestDepth > minWorstCaseDepth {
				currentBestDepth = minWorstCaseDepth
			}
		}
		return minWorstCaseGuess, minWorstCaseDepth
	}
	return func(history History, currentlyAllowed game.WordList) game.Word {
		guess, worstCaseDepth := calculator(history, currentlyAllowed, math.MaxInt)
		fmt.Printf("Guess: %s, Max Depth: %d\n", guess, worstCaseDepth)
		return guess
	}
}

func main() {
	wordlist := game.ReadWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	//dm := Daemon{toWord("atone")}
	//dm.playInteractively()

	//split := splitAfterGuess(wordlist[0], wordlist)
	//log.Printf("split list by %s:\n%v\n", wordlist[0], split)

	dm := &game.Daemon{CorrectWord: wordlist[10]}
	length, history := play(wordlist, dm, GreedyNeedfulPlayer(wordlist, game.ToWord("atone")))

	log.Printf("Naive player guesses it in %d rounds by:\n%v.\n", length, history)

	log.Println("Done!")
}
