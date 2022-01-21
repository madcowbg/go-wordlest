package game

import (
	"bytes"
	"fmt"
	"log"
)

type Word struct{ chars [5]byte } // encoded so 'a' is 0
func (word Word) String() string {
	var b bytes.Buffer
	for i := range word.chars {
		b.WriteByte(word.chars[i] + 'a')
	}
	return b.String()
}

func ToWord(word string) Word {
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
