package words

import (
	"testing"
	"unicode/utf8"
)

func GetDictionary() *Dictionary {
	context := GetTokenizeContext()
	return context.dict
}

func TestAddWord(t *testing.T) {
	dict := Dictionary{}
	word := Word{}
	dict.AddWord("test", &word)
	if dict.root.Children['t'].Children['e'].Children['s'].Children['t'].Word == nil {
		t.Errorf("test is not a word")
	}
	wordptr, _ := dict.FindWord("test")
	if wordptr == nil {
		t.Errorf("test is not a word")
	}
}

func TestWalk(t *testing.T) {
	dict := Dictionary{}

	word := Word{}
	dict.AddWord("test", &word)

	wordch := make(chan *Word)

	go dict.Walk(wordch)

	for {
		word, ok := <-wordch
		if ok {
			if word.String() != "test" {
				t.Errorf("word is not test: %s\n", word.String())
			}
		} else {
			break
		}
	}
	return
}

func TestWalkOfSize(t *testing.T) {
	dict := GetDictionary()
	wordch := make(chan *Word)

	go dict.WalkOfSize(4, wordch)

	for {
		word, ok := <-wordch
		if ok {
			wordlen := utf8.RuneCountInString(word.String())
			if wordlen != 4 {
				t.Errorf("bad word len %d %s", wordlen, word.String())
			}
		} else {
			break
		}
	}
	return
}

func TestBinary(t *testing.T) {
	writedict := Dictionary{}

	bananeword := new(Word)
	writedict.AddWord("banane", bananeword)
	testword := new(Word)
	writedict.AddWord("test", testword)

	word, _ := writedict.FindWord("test")
	if word == nil {
		t.Errorf("test was not added")
	}

	bananewords, _ := writedict.FindWord("banane")
	if bananewords == nil || bananewords.String() != "banane" {
		t.Errorf("banane was not added")
	}

	// write test
	err := writedict.WriteBinary("test.bin")
	if err != nil {
		t.Errorf("cannot write binary %s", err)
	}

	// read test
	readdict := Dictionary{}
	err = readdict.ReadBinary("test.bin")
	if err != nil {
		t.Errorf("cannot read binary %s", err)
	}

	word, _ = readdict.FindWord("test")
	if word == nil {
		t.Errorf("test not found")
	}
	word, _ = readdict.FindWord("banane")
	if word == nil {
		t.Errorf("banane not found")
	}

	return
}

func TestFindWord(t *testing.T) {
	dict := GetDictionary()
	word, _ := dict.FindWord("à")
	if word == nil {
		t.Errorf("à not found")
	}
	if word.LastLetter.GetWord() != "à" {
		t.Errorf("a get word mismatch %s", word.LastLetter.GetWord())
	}

	word, _ = dict.FindWord("depuis")
	if word == nil {
		t.Errorf("depuis not found")
	}

	if word.LastLetter.GetWord() != "depuis" {
		t.Errorf("depuis get word mismatch")
	}

	word, _ = dict.FindWord("belle.")
	if word != nil {
		t.Errorf("belle. found")
	}

	word, _ = dict.FindWord(" ")
	if word == nil {
		t.Errorf("space not found")
	}
}

func TestFindLonguestWord(t *testing.T) {
	dict := GetDictionary()
	word := dict.FindLonguestWord("à affirmer la vie")
	if word == nil {
		t.Errorf("à affirmer not found")
	}
	if word.LastLetter.GetWord() != "à affirmer" {
		t.Errorf("a affirmer get word mismatch %s", word.LastLetter.GetWord())
	}
}

func TestFindAlternative(t *testing.T) {
	dict := GetDictionary()

	missing, _ := dict.FindWord("tast")
	if missing != nil {
		t.Errorf("tast word is found")
	}
	alternatives := dict.FindAlternatives("tast", FRENCH, 1)

	found := false
	for _, words := range alternatives {
		if words.String() == "test" {
			found = true
		}
	}
	if !found {
		t.Errorf("could not find test")
	}
}

func TestAutoComplete(t *testing.T) {
	dict := GetDictionary()

	completed := dict.AutoComplete("lance", nil)

	found := false
	for _, word := range completed {
		if word.String() == "lancer" {
			found = true
		}
	}
	if !found {
		t.Errorf("lancer not found")
	}
}

func TestAutoCompleteFilter(t *testing.T) {
	dict := GetDictionary()
	filter := WordVariant{}

	filter.Tag = VERB
	filter.Language = FRENCH
	filter.Person = 1
	filter.Tense = IND
	filter.Number = SINGULAR

	completed := dict.AutoComplete("lanc", &filter)

	found := false
	for _, word := range completed {
		if word.String() == "lance" {
			found = true
		}
	}
	if !found {
		t.Errorf("lance not found")
	}
}
