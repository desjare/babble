package words

import (
	"bufio"
	"io"
	"math"
	"os"
	"strings"
	"unicode/utf8"
)

type Dictionary struct {
	root       WordLetter
	MaxWordLen int
	MaxTokens  int
}

func (dict *Dictionary) ReadLanguage(path string, lang byte) (err error) {
	words := WordList{}
	err = words.Read(path, lang, func(form string, word *Word)(error){
		dict.AddWord(form, word)
		return nil
	})
	if err != nil {
		return err
	}
	return
}

func (dict *Dictionary) ReadXML() error {
	err := dict.ReadLanguage("lmdata/dela-fr-public-u8.dic.xml", FRENCH)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-addon-fr-u8.dic.xml", FRENCH)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-abbr-fr-u8.dic.xml", FRENCH)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-en-public-u8.dic.xml", ENGLISH)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-proper-u8.dic.xml", 0)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-acron-u8.dic.xml", 0)
	if err != nil {
		panic(err)
		return err
	}
	err = dict.ReadLanguage("lmdata/dela-punc-u8.dic.xml", 0)
	if err != nil {
		panic(err)
		return err
	}

	word := new(Word)
	dict.AddWord("...", word)

	dict.AddBuiltin()

	return err
}

func (dict *Dictionary) AddWordWithTag(s string, tag byte) {

	word := Word{}
	variant := WordVariant{}
	variant.Tag = tag
	word.Variants = append(word.Variants, variant)

	dict.AddWord(s, &word)
}

func (dict *Dictionary) AddBuiltin() {

	dict.AddWordWithTag("\t", TAB)
	dict.AddWordWithTag("\n", EOL)
	dict.AddWordWithTag("\r", CR)
	dict.AddWordWithTag("/", SLASH)
	dict.AddWordWithTag(">", GREATHERTHAN)
	dict.AddWordWithTag("<", SMALLERTHAN)
	dict.AddWordWithTag("&", AND)
	dict.AddWordWithTag("|", OR)
	dict.AddWordWithTag("[", BEGINBRACKET)
	dict.AddWordWithTag("]", ENDBRACKET)
	dict.AddWordWithTag("#", NUMBERSIGN)
	dict.AddWordWithTag("$", DOLLARSIGN)
	dict.AddWordWithTag("@", ATSIGN)
	dict.AddWordWithTag("©", COPYRIGHTSIGN)
}

func (dict *Dictionary) AddWord(letters string, word *Word) {
	dict.root.AddWord(letters, word)
	dict.MaxWordLen = int(math.Max(float64(len(letters)), float64(dict.MaxWordLen)))
	dict.MaxTokens = int(math.Max(float64(strings.Count(letters, " ")), float64(dict.MaxTokens)))
}

func (dict *Dictionary) FindWord(letters string) (*Word, bool) {
	return dict.root.FindWord(letters)
}

func (dict *Dictionary) FindLonguestWord(letters string) *Word {
	return dict.root.FindLonguestWord(letters)
}

func (dict *Dictionary) FindPath(letters string) *WordLetter {
	return dict.root.FindPath(letters)
}

func (dict *Dictionary) FindAlternatives(word string, lang byte, maxerror int) []*Word {
	var words []*Word

	wordch := make(chan *Word)
	count := utf8.RuneCountInString(word)
	go dict.WalkOfSize(count, wordch)

	for {
		w, ok := <-wordch
		if ok {
			errors := CountError(word, w.String())
			if errors <= maxerror {
				if w.Language(lang) {
					words = append(words, w)
				}
			}
		} else {
			break
		}
	}
	return words
}

func (dict *Dictionary) AutoComplete(word string, filter *WordVariant) []*Word {
	var words []*Word

	wordch := make(chan *Word)
	go dict.WalkFromPath(word, wordch)

	for {
		w, ok := <-wordch
		if ok {
			if filter == nil || filter.Filter(w) {
				words = append(words, w)
			}
		} else {
			break
		}
	}
	return words
}

func (dict *Dictionary) Walk(wordch chan *Word) {
	dict.root.Walk(wordch)
	close(wordch)
	return
}

func (dict *Dictionary) WalkOfSize(size int, wordch chan *Word) {
	dict.root.WalkOfSize(size, wordch)
	close(wordch)
	return
}

func (dict *Dictionary) WalkFromPath(word string, wordch chan *Word) {
	wordletter := dict.FindPath(word)
	if wordletter != nil {
		wordletter.Walk(wordch)
	}
	close(wordch)
	return
}

func (dict *Dictionary) WriteBinary(path string) (err error) {
	wordch := make(chan *Word)

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	go dict.Walk(wordch)
	for {
		word, ok := <-wordch
		if ok {
			word.Write(w)
		} else {
			break
		}

	}
	return
}

func (dict *Dictionary) ReadBinary(path string) (err error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	r := bufio.NewReaderSize(f, 1024*1024)
	for {
		word := Word{}
		letters, err := word.Read(r)

		if err == nil {
			dict.AddWord(letters, &word)
		} else {
			if err != io.EOF {
				return err
			} else {
				err = nil
				break
			}
		}
	}
	return err
}
