package main

import (
	"bytes"
	"encoding/csv"
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

func calcAnswer(word Word) func(guess Word) Ans {
	return func(guess Word) Ans {
		letters := [26]byte{}
		for i := range word {
			letters[byte(word[i])]++
		}
		ans := [5]byte{}
		// mark all correct places and guesses
		for i := range guess {
			if guess[i] == word[i] {
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
}

func toWord(word string) Word {
	res := [5]byte{}
	if len(word) != 5 {
		log.Fatalf("word must be 5 letters! Got [%s] which is [%d]", word, len(word))
	}
	for i := range word {
		if !('a' <= word[i] && word[i] <= 'z') {
			log.Fatalf("using forbidden character, must be 'a'...'z', got [%d]", word[i])
		}
		res[i] = word[i] - 'a'
	}
	return res
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

func main() {
	wordlist := readWordleList("data/wordle-answers-alphabetical.txt")

	log.Printf("word list count: %d\n", len(wordlist))
	log.Println("Done!")
}
