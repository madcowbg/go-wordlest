package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func readWordleList(filePath string) []string {
	loaded_csv := readCsvFile(filePath)
	wordlist := make([]string, len(loaded_csv))

	for i := range loaded_csv {
		wordlist[i] = loaded_csv[i][0]
	}
	return wordlist
}

type Ans = [5]byte  // 0, 1, 2
type Word = [5]byte // encoded so 'a' is 0

func toWord(word string) Word {
	err := checkValidWord(word)
	if err != nil {
		log.Fatal(err)
	}
	res := [5]byte{}
	for i := range word {
		res[i] = word[i] - 'a'
	}
	return res
}

func checkValidWord(word string) error {
	if len(word) != 5 {
		return fmt.Errorf("word must be 5 letters! Got [%s] which is [%d]", word, len(word))
	}
	for i := range word {
		if !('a' <= word[i] && word[i] <= 'z') {
			return fmt.Errorf("using forbidden character, must be 'a'...'z', got [%d]", word[i])
		}
	}
	return nil
}

func wordToStr(word Word) string {
	var b bytes.Buffer
	for i := range word {
		b.WriteByte(word[i] + 'a')
	}
	return b.String()
}

func ansToStr(ans Ans) string {
	var b bytes.Buffer
	for i := range ans {
		b.WriteByte(ans[i] + '0')
	}
	return b.String()
}

type Daemon struct {
	correctWord Word
}

func (dm *Daemon) ask(guess Word) Ans {
	letters := [26]byte{}
	for i := range dm.correctWord {
		letters[byte(dm.correctWord[i])]++
	}
	ans := [5]byte{}
	// mark all correct places and guesses
	for i := range guess {
		if guess[i] == dm.correctWord[i] {
			ans[i] = 2
			letters[byte(guess[i])]--
		}
	}
	// mark all guesses in incorrect places
	for i := range guess {
		if ans[i] == 0 && letters[byte(guess[i])] > 0 {
			ans[i] = 1
			letters[byte(guess[i])]--
		}
	}
	return ans
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
		ans := dm.ask(guessWrd)
		if ansToStr(ans) == "22222" {
			fmt.Printf("Word guessed!\n")
			return
		} else {
			fmt.Printf("Daemon says %s.\n", ansToStr(ans))
		}
		rnd += 1
	}
}

func main() {
	wordlist := readWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	dm := Daemon{toWord("atone")}
	dm.playInteractively()

	log.Println("Done!")

}
