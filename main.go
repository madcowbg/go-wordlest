package main

import (
	"bytes"
	"fmt"
	"go-wordlyest/game"
	"log"
	"math"
)

func splitAfterGuess(guess game.Word, wordlist game.WordList) map[game.Ans]game.WordList {
	result := make(map[game.Ans]game.WordList)
	for i := range wordlist {
		word := wordlist[i]
		dm := game.Daemon{CorrectWord: word}
		ans := dm.Ask(guess)

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

func MinMaxPlayer(verbose bool) Player {
	var calculator func(depth int, currentlyAllowed game.WordList, currentBestDepth int, currentBestDepthChan chan int) (game.Word, int)
	calculator = func(depth int, currentlyAllowed game.WordList, currentBestDepth int, currentBestDepthChan chan int) (game.Word, int) {
		if len(currentlyAllowed) == 1 {
			// base case - we are here!
			return currentlyAllowed[0], depth + 1
		}

		if depth >= currentBestDepth {
			return currentlyAllowed[0], math.MaxInt // prune case - we can't beat it, so may as well skip it and declare undesirable
		}

		var minWorstCaseGuess game.Word
		var minWorstCaseDepth = math.MaxInt
		maxDepth := make(map[game.Word]int)
		for _, guess := range currentlyAllowed {
			if verbose && depth == 0 {
				fmt.Printf("Trying %s...\n", guess)
			}
			afterGuess := splitAfterGuess(guess, currentlyAllowed)
			maxDepth[guess] = 0
			for _, afterGuessAllowed := range afterGuess {
				_, depthInThatHistory := calculator(
					depth+1,
					afterGuessAllowed,
					currentBestDepth,
					currentBestDepthChan)
				if maxDepth[guess] < depthInThatHistory {
					maxDepth[guess] = depthInThatHistory
				}
				if maxDepth[guess] >= currentBestDepth {
					maxDepth[guess] = math.MaxInt // prune when any other ans will pull us further into undesirable paths...
					break
				}
			}
			if verbose && depth <= 0 {
				fmt.Printf("Alternative checked %s, determined improvement in minmaxmax depth %d\n", guess, maxDepth[guess])
			}
			if minWorstCaseDepth > maxDepth[guess] {
				minWorstCaseDepth = maxDepth[guess]
				minWorstCaseGuess = guess
			}
			if currentBestDepth > minWorstCaseDepth {
				currentBestDepth = minWorstCaseDepth
			}
		}
		return minWorstCaseGuess, minWorstCaseDepth
	}
	return func(history History, currentlyAllowed game.WordList) game.Word {
		currentBestDepthChan := make(chan int)

		guess, worstCaseDepth := calculator(len(history), currentlyAllowed, math.MaxInt, currentBestDepthChan)
		if verbose {
			fmt.Printf("Guess: %s, Max Depth: %d\n", guess, worstCaseDepth)
		}
		return guess
	}
}

func FastFirstHand(firstGuess game.Word, fallback Player) Player {
	return func(history History, allowed game.WordList) game.Word {
		if len(history) == 0 {
			return firstGuess
		} else {
			return fallback(history, allowed)
		}
	}
}

func main() {
	wordlist := game.ReadWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	//dm := Daemon{toWord("atone")}
	//dm.playInteractively()

	//split := splitAfterGuess(wordlist[0], wordlist)
	//log.Printf("split list by %s:\n%v\n", wordlist[0], split)

	//dm := &game.Daemon{CorrectWord: wordlist[10]}
	//length, history := play(wordlist, dm, GreedyNeedfulPlayer(wordlist, game.ToWord("atone")))
	//
	//log.Printf("Naive player guesses it in %d rounds by:\n%v.\n", length, history)

	//maxSizeAfterGuess := make(map[game.Word]int)
	//for _, guess := range wordlist {
	//	maxSizeAfterGuess[guess] = 0
	//	for _, possibilitiesAfterGuess := range splitAfterGuess(guess, wordlist) {
	//		if maxSizeAfterGuess[guess] < len(possibilitiesAfterGuess) {
	//			maxSizeAfterGuess[guess] = len(possibilitiesAfterGuess)
	//		}
	//	}
	//	fmt.Printf("%s\t%d\n", guess, maxSizeAfterGuess[guess])
	//}

	reducedWordlist := wordlist
	for _, wrd := range reducedWordlist {
		dm := &game.Daemon{CorrectWord: wrd}
		numRounds, _ := play(reducedWordlist, dm, MinMaxPlayer(true))
		log.Printf("%s\t%d\n", wrd, numRounds)
		break
	}

	log.Println("Done!")
}
