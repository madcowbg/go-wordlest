package game

import (
	"bytes"
	"fmt"
	"log"
)

type Word struct{ value int } // encoded so 'a' is 0
func (word Word) String() string {
	var b bytes.Buffer
	val := word.value
	for i := 0; i < 5; i++ {
		b.WriteByte(byte(val%26) + 'a')
		val /= 26
	}
	return b.String()
}

func ToWord(word string) Word {
	err := checkValidWord(word)
	if err != nil {
		log.Fatal(err)
	}
	value := 0
	for i := range word {
		value = value*26 + int(word[4-i]-'a')
	}
	return Word{value}
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
