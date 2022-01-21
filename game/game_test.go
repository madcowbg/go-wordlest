package game

import (
	"testing"
)

func TestToWord(t *testing.T) {
	str := "boobs"
	if ToWord(str).value != 1+26*(14+26*(14+26*(1+26*18))) {
		t.Fatalf("conversion failed! %d != %d", ToWord(str).value, 1+26*(14+26*(14+26*(1+26*18))))
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
	dm := Daemon{ToWord(word)}
	ans := dm.Ask(ToWord(guess))
	if ans.String() != expAns {
		t.Fatalf("%s -> %s must give %s, gave [%s]", word, guess, expAns, ans)
	}
}

func TestAnsToString(t *testing.T) {
	ans := FromBytes([5]byte{2, 0, 1, 1, 2})
	if ans.String() != "20112" {
		t.Fatalf("ans conversion failed for %v, expected %s but got %s", ans, "20112", ans)
	}
}

func TestWordToString(t *testing.T) {
	wrd := ToWord("alzer")
	if wrd.String() != "alzer" {
		t.Fatalf("word to string failed to convert %v to %s", ToWord("alzer"), "alzer")
	}
}
