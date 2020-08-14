package words

import (
	"testing"
)

func TestDetNoun(t *testing.T) {
	grammar := GrammarNew()

	context := GetTokenizeContext()
	tokens := Tokenize(" Un chien", context, true)

	constructs := grammar.Match(tokens)

	if len(constructs) == 1 {
		if constructs[0].Name != DET_NOUN {
			t.Errorf("Un chien not a det noun")
		}

	} else {
		t.Errorf("Un chien no det noun found: %d", len(constructs))
	}
}
