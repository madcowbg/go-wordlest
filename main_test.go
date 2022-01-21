package main

import (
	"fmt"
	"testing"
)

func TestToWord(t *testing.T) {
	str := "boobs"
	wordRepr := fmt.Sprintf("%v", toWord(str))
	if wordRepr != "[1 14 14 1 18]" {
		t.Fatalf("conversion failed!")
	}
}

func TestAskDaemon(t *testing.T) {
	checkCase("boobs", "boobs", "22222", t)
	checkCase("boobs", "potty", "02000", t)
	checkCase("boobs", "babes", "20102", t)
	checkCase("abbab", "baaaa", "11020", t)
	checkCase("baaaa", "abbab", "11020", t)
}

func checkCase(word string, guess string, expAns string, t *testing.T) {
	dm := Daemon{toWord(word)}
	ansStr := ansToStr(dm.ask(toWord(guess)))
	if ansStr != expAns {
		t.Fatalf("%s -> %s must give %s, gave [%s]", word, guess, expAns, ansStr)
	}
}

func TestAnsToString(t *testing.T) {
	ans := [5]byte{2, 0, 1, 1, 2}
	if ansToStr(ans) != "20112" {
		t.Fatalf("ans conversion failed for %v, expected %s but got %s", ans, "20112", ansToStr(ans))
	}
}

func TestWordToString(t *testing.T) {
	wrd := toWord("alzer")
	if wordToStr(wrd) != "alzer" {
		t.Fatalf("word to string failed to convert %v to %s", toWord("alzer"), "alzer")
	}
}
