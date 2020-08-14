package words

import (
	"fmt"
)

const (
	PHRASE_UNKNOWN = 0
	PHRASE_GN_NOUN_PHRASE = 1
	PHRASE_GN_NOUN = 2
	PHRASE_GN_DET_NOUN = 3
	PHRASE_GN_DET_ADN_NOUN = 4
	PHRASE_GN_DET_NOUN_ADJ = 5
)

type PhraseConstructDef struct {
	Tags []byte
	Name byte
}

type PhraseConstruct struct {
	Name byte
	Tokens []Token
	Childs []PhraseConstruct
}


type Grammar struct {
	Defs []PhraseConstructDef
}

func GrammarNew() *Grammar {
	grammar := new(Grammar)

	return grammar
}

func (g *Grammar) MatchDef(d PhraseConstructDef, tokens []Token) (c *PhraseConstruct, next int) {
	var tag int
	var tok int

	numtags := len(d.Tags)
	numtokens := len(tokens)
	c = new(PhraseConstruct)

	for tag < numtags && tok < numtokens {
		if TokenizeIsTag(tokens[tok],SPACE) {
			tok++
			continue
		}
		if !TokenizeIsTag(tokens[tok],d.Tags[tag]) {
			fmt.Printf("tag mismatch %s tag %d\n", tokens[tok].String(), d.Tags[tag])
			return
		}

		if tag == numtags -1 {
			c.Tokens = tokens[0:tok]
			c.Name = d.Name
			next = tok+1
			return
		}
		tag++
		tok++
	}
	return

}

func (g *Grammar) Match(tokens []Token) (constructs []PhraseConstruct) {
	for _, d := range(g.Defs) {
		c, _:= g.MatchDef(d, tokens)
		if c != nil {
			constructs = append(constructs, *c)
		}
	}
	return
}
