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
	val := dm.CorrectWord.value
	for i := 0; i < 5; i++ {
		letters[byte(val%26)]++
		val /= 26
	}

	ans := [5]byte{}
	// mark all correct places and guesses
	gVal := guess.value
	cVal := dm.CorrectWord.value
	for i := 0; i < 5; i++ {
		if gVal%26 == cVal%26 {
			ans[i] = 2
			letters[byte(gVal%26)]--
		}
		gVal /= 26
		cVal /= 26
	}
	// mark all guesses in incorrect places
	gVal = guess.value          // reset
	cVal = dm.CorrectWord.value // reset
	for i := 0; i < 5; i++ {
		if ans[i] == 0 && letters[byte(gVal%26)] > 0 {
			ans[i] = 1
			letters[byte(gVal%26)]--
		}
		gVal /= 26
		cVal /= 26
	}
	return FromBytes(ans)
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

		guessWrd := ToWord(guess)
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
