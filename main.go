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
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("Can't close file "+filePath, err)
		}
	}(f)

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

func readWordleList(filePath string) WordList {
	loaded_csv := readCsvFile(filePath)
	wordlist := make([]Word, len(loaded_csv))

	for i := range loaded_csv {
		wordlist[i] = toWord(loaded_csv[i][0])
	}
	return wordlist
}

type Ans struct{ bytes [5]byte } // 0, 1, 2
func (ans Ans) String() string {
	var b bytes.Buffer
	for i := range ans.bytes {
		b.WriteByte(ans.bytes[i] + '0')
	}
	return b.String()
}

type Word struct{ chars [5]byte } // encoded so 'a' is 0
func (word Word) String() string {
	var b bytes.Buffer
	for i := range word.chars {
		b.WriteByte(word.chars[i] + 'a')
	}
	return b.String()
}

func toWord(word string) Word {
	err := checkValidWord(word)
	if err != nil {
		log.Fatal(err)
	}
	res := [5]byte{}
	for i := range word {
		res[i] = word[i] - 'a'
	}
	return Word{res}
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

type Daemon struct {
	correctWord Word
}

func (dm *Daemon) ask(guess Word) Ans {
	letters := [26]byte{}
	for i := range dm.correctWord.chars {
		letters[byte(dm.correctWord.chars[i])]++
	}
	ans := [5]byte{}
	// mark all correct places and guesses
	for i := range guess.chars {
		if guess.chars[i] == dm.correctWord.chars[i] {
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
		ans := dm.ask(guessWrd)
		if ans.String() == "22222" {
			fmt.Printf("Word guessed!\n")
			return
		} else {
			fmt.Printf("Daemon says %s.\n", ans)
		}
		rnd += 1
	}
}

type WordList []Word

func splitAfterGuess(guess Word, wordlist WordList) map[Ans]WordList {
	dm := Daemon{guess}
	result := make(map[Ans]WordList)
	for i := range wordlist {
		word := wordlist[i]
		ans := dm.ask(word)

		existantList, found := result[ans]
		if !found {
			existantList = []Word{}
		}

		result[ans] = append(existantList, word)
	}
	return result
}

func main() {
	wordlist := readWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))

	//dm := Daemon{toWord("atone")}
	//dm.playInteractively()

	split := splitAfterGuess(wordlist[0], wordlist)
	log.Printf("split list by %s:\n%v\n", wordlist[0], split)
	log.Println("Done!")
}
