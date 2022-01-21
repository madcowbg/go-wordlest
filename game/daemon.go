package game

import (
	"fmt"
	"log"
)

type Daemon struct {
	CorrectWord Word
}

func (dm *Daemon) Ask(guess Word) Ans {
	letters := [26]byte{}
	for i := range dm.CorrectWord.chars {
		letters[byte(dm.CorrectWord.chars[i])]++
	}
	ans := [5]byte{}
	// mark all correct places and guesses
	for i := range guess.chars {
		if guess.chars[i] == dm.CorrectWord.chars[i] {
			ans[i] = 2
			letters[byte(guess.chars[i])]--
		}
	}
	// mark all guesses in incorrect places
	for i := range guess.chars {
		if ans[i] == 0 && letters[byte(guess.chars[i])] > 0 {
			ans[i] = 1
			letters[byte(guess.chars[i])]--
		}
	}
	return Ans{ans}
}

func (dm *Daemon) playInteractively() {
	fmt.Println("Starting interactive game...")
	rnd := 1
	for {
		fmt.Printf("Round %d: ", rnd)
		var guess string
		_, err := fmt.Scanln(&guess)
		if err != nil {
			log.Fatal(err)
		}

		err = checkValidWord(guess)
		if err != nil {
			fmt.Printf("Invalid word: [%s]\n", err)
			continue
		}

		guessWrd := toWord(guess)
		ans := dm.Ask(guessWrd)
		if ans.String() == "22222" {
			fmt.Printf("Word guessed!\n")
			return
		} else {
			fmt.Printf("Daemon says %s.\n", ans)
		}
		rnd += 1
	}
}
