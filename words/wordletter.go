package words

import (
	"sync"
	"unicode/utf8"
)

type WordLetter struct {
	Letter   rune
	Word     *Word
	Parent   *WordLetter
	Children map[rune]*WordLetter
}

func (root *WordLetter) AddWord(letters string, word *Word) {

	if root.Children == nil {
		root.Children = make(map[rune]*WordLetter)
	}
	letter, width := utf8.DecodeRuneInString(letters)
	if root.Children[letter] == nil && width > 0 {
		node := WordLetter{}
		node.Parent = root
		node.Letter = letter
		node.Children = make(map[rune]*WordLetter)
		root.Children[letter] = &node
	}

	if len(letters) > width {
		root.Children[letter].AddWord(letters[width:], word)
	} else {
		// do not add same word twice
		if root.Children[letter].Word != nil {
			root.Children[letter].Word.AddVariants(word.Variants)
		} else {
			root.Children[letter].Word = word
			word.LastLetter = root.Children[letter]
		}
	}
}

func (root *WordLetter) FindWord(word string) (*Word, bool) {
	var wordptr *Word
	var letter rune
	last := root

	for i, w := 0, 0; i < len(word); i += w {
		letter, w = utf8.DecodeRuneInString(word[i:])
		if last.Children[letter] != nil {
			last = last.Children[letter]
			wordptr = last.Word
		} else {
			return nil, false
		}
	}
	return wordptr, true
}

func (root *WordLetter) FindPath(word string) *WordLetter {
	var letter rune
	var wordletter *WordLetter
	last := root

	for i, w := 0, 0; i < len(word); i += w {
		letter, w = utf8.DecodeRuneInString(word[i:])
		if last.Children[letter] != nil {
			last = last.Children[letter]
			wordletter = last
		} else {
			return nil
		}
	}
	return wordletter
}

func (root *WordLetter) FindLonguestWord(word string) *Word {
	var wordptr *Word
	var letter rune
	last := root

	for i, w := 0, 0; i < len(word); i += w {
		letter, w = utf8.DecodeRuneInString(word[i:])
		if last.Children[letter] != nil {
			last = last.Children[letter]
			if last.Word != nil {
				wordptr = last.Word
			}
		} else {
			break
		}
	}

	return wordptr
}

func (last *WordLetter) GetWord() string {
	var word []byte
	for l := last; l.Letter != 0; l = l.Parent {
		tmp := []byte{0, 0, 0, 0}
		width := utf8.EncodeRune(tmp, l.Letter)
		word = append(tmp[0:width], word...)
	}
	return string(word)
}

func (letter *WordLetter) Walk(wordch chan *Word) {
	if letter.Word != nil {
		wordch <- letter.Word
	}

	var wg sync.WaitGroup
	for _, subletter := range letter.Children {
		wg.Add(1)
		go func(sub *WordLetter, wordch chan *Word) {
			sub.Walk(wordch)
			wg.Done()
		}(subletter, wordch)
	}
	wg.Wait()

	return
}

func (letter *WordLetter) WalkOfSize(runesize int, wordch chan *Word) {
	if runesize == 0 && letter.Word != nil {
		wordch <- letter.Word
	}

	if runesize > 0 {
		var wg sync.WaitGroup
		for _, subletter := range letter.Children {
			wg.Add(1)
			go func(size int, sub *WordLetter, wordch chan *Word) {
				sub.WalkOfSize(size, wordch)
				wg.Done()
			}(runesize-1, subletter, wordch)
		}
		wg.Wait()
	}
	return
}
