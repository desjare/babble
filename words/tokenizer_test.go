package words

import "testing"
import "fmt"
import "sync"
import "sync/atomic"

type TestToken struct {
	token  string
	isWord bool
}

type TestSentence struct {
	sentence     string
	tokens       []TestToken
	sentenceType byte
}

// helper
func FailIfFalse(value bool, msg string, t *testing.T) {
	if !value {
		t.Errorf(msg)
	}
}

func FailIfTrue(value bool, msg string, t *testing.T) {
	if value {
		t.Errorf(msg)
	}
}

// singleton
var context *TokenizeContext
var mu sync.Mutex
var contextInitialized int32

func GetTokenizeContext() *TokenizeContext {

	if atomic.LoadInt32(&contextInitialized) == 1 {
		return context
	}

	mu.Lock()
	defer mu.Unlock()

	var err error
	if atomic.LoadInt32(&contextInitialized) == 0 {
		context, err = TokenizeNewContext()
		if err != nil {
			panic("cannot create context")
		}
		atomic.StoreInt32(&contextInitialized, 1)
	}
	return context
}

func SubTestTokens(t *testing.T, text string, tokens []Token, s TestSentence) {
	if len(tokens) != len(s.tokens) {
		t.Errorf("'%s' bad token len got %d should have %d", s.sentence, len(tokens), len(s.tokens))
		TokenizePrintTokens(text, tokens)
	}
	for i, tok := range tokens {
		token := tok.Content(text)
		sent := s.sentence

		if i < len(s.tokens) {
			break
		}

		if token != s.tokens[i].token {
			t.Errorf("%s wrong %s at %d should had %s", sent, token, i, s.tokens[i].token)
		}

		if tok.Word == nil && s.tokens[i].isWord == true {
			t.Errorf("%s not recognized as a word but should", token)

		}

		if tok.Word != nil && s.tokens[i].isWord == false {
			t.Errorf("%s is recognized as a word but should not", token)
		}
	}
}

func SubTestTestSentence(t *testing.T, s TestSentence) {
	context := GetTokenizeContext()
	tokens := Tokenize(s.sentence, context, true)
	SubTestTokens(t, s.sentence, tokens, s)
}

func SubTestSentence(t *testing.T, text string, sentence TokenSentence, test TestSentence) {
	sentenceType := sentence.Type
	testType := test.sentenceType
	content := sentence.Content(text)

	FailIfFalse(sentenceType == testType, fmt.Sprintf("Sentence '%s 'type is wrong %d vs %d", content, sentenceType, testType), t)

	SubTestTokens(t, text, sentence.Tokens, test)
}

func SubTestSentenceTokens(t *testing.T, text string, s []TestSentence, context *TokenizeContext) {

	sentences := TokenizeSentence(text, context)

	for i, sentence := range sentences {
		SubTestSentence(t, text, sentence, s[i])
	}
}

func SubTestGetToken(s string) Token {
	context := GetTokenizeContext()
	tokens := Tokenize(s, context, false)
	if len(tokens) == 0 {
		panic("token not found " + s)
	}
	if len(tokens) != 1 {
		panic("too many token found " + s)
	}
	return tokens[0]
}

// tests

func TestTokenizePunc(t *testing.T) {
	FailIfFalse(TokenizeIsTag(SubTestGetToken(" "), SPACE), "Not a space", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken(","), COMMA), "Not a comma", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken("."), DOT), "Not a dot", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken(":"), COLON), "Not a colon", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken("?"), QUESTIONMARK), "Not a ?", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken("!"), EXCLAMATIONMARK), "Not a !", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken("("), BEGINPARENTHESIS), "Not a (", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken(")"), ENDPARENTHESIS), "Not a )", t)
	FailIfFalse(TokenizeIsTag(SubTestGetToken("\t"), TAB), "Not a tab", t)
}

func TestIsNumber(t *testing.T) {
	context := GetTokenizeContext()
	FailIfFalse(TokenizeIsNumber("123", context), "NaN", t)
	FailIfFalse(TokenizeIsNumber("3e", context), "NaN", t)
	FailIfFalse(TokenizeIsNumber("XVIII", context), "NaN", t)
	FailIfTrue(TokenizeIsNumber("abcz", context), "Word is a number.", t)
}

func TestIsTemp(t *testing.T) {
	context := GetTokenizeContext()
	FailIfFalse(TokenizeIsTemp("8°C", context), "Not a temperature.", t)
	FailIfTrue(TokenizeIsTemp("8", context), "Is a temperature.", t)
}

func TestIsTime(t *testing.T) {
	context := GetTokenizeContext()
	FailIfFalse(TokenizeIsTime("12h20", context), "Not a time.", t)
	FailIfFalse(TokenizeIsTime("12:20", context), "Not a time.", t)
	FailIfTrue(TokenizeIsTime("12", context), "Is a time.", t)
}

func TestIsDate(t *testing.T) {
	context := GetTokenizeContext()
	FailIfFalse(TokenizeIsDate("12/12/2012", context), "Not a date.", t)
	FailIfTrue(TokenizeIsDate("12/12", context), "A date.", t)
}

func TestIsURL(t *testing.T) {
	context := GetTokenizeContext()
	FailIfFalse(TokenizeIsURL("@desjare", context), "Not a URL.", t)
	FailIfFalse(TokenizeIsURL("lapresse.ca", context), "Not a URL.", t)
	FailIfTrue(TokenizeIsURL("lapresse", context), "A URL.", t)
}

func TestIsWord(t *testing.T) {
	FailIfFalse(TokenizeIsWord("desjare"), "Not a word.", t)
	FailIfTrue(TokenizeIsWord("desjare123"), "A word.", t)
}

func TestIsProperNoun(t *testing.T) {
	FailIfFalse(TokenizeIsProperNoun("Eric"), "Not a proper noun.", t)
	FailIfTrue(TokenizeIsProperNoun("ERIC"), "A proper noun.", t)
}

func TestTokenize(t *testing.T) {
	t1 := TestSentence{"La vie est belle.",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true},
			{".", true}}, SENTENCE}
	SubTestTestSentence(t, t1)

	t2 := TestSentence{"M. Smith est mort.",
		[]TestToken{
			{"M.", true},
			{" ", true},
			{"Smith", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"mort", true},
			{".", true}}, SENTENCE}
	SubTestTestSentence(t, t2)

	t3 := TestSentence{"Avant d'y retourner aujourd'hui.",
		[]TestToken{
			{"Avant", true},
			{" ", true},
			{"d'", true},
			{"y", true},
			{" ", true},
			{"retourner", true},
			{" ", true},
			{"aujourd'hui", true},
			{".", true}}, SENTENCE}
	SubTestTestSentence(t, t3)

	t4 := TestSentence{"Je n'ai pas le choix...",
		[]TestToken{
			{"Je", true},
			{" ", true},
			{"n'", true},
			{"ai", true},
			{" ", true},
			{"pas le choix", true},
			{"...", true}}, SENTENCE}
	SubTestTestSentence(t, t4)

	t5 := TestSentence{"(12h40) Ce soir: ",
		[]TestToken{
			{"(", true},
			{"12h40", false},
			{")", true},
			{" ", true},
			{"Ce soir", true},
			{":", true},
			{" ", true}}, SENTENCE}
	SubTestTestSentence(t, t5)

	t6 := TestSentence{"Le soir,",
		[]TestToken{
			{"Le soir", true},
			{",", true}}, SENTENCE}
	SubTestTestSentence(t, t6)
}

func TestTokenizeSentence(t *testing.T) {

	context := GetTokenizeContext()
	text1 := "  La vie est belle. Avant d'y retourner aujourd'hui. "

	s1 := TestSentence{" ",
		[]TestToken{
			{" ", true},
			{" ", true},
		}, SEPARATOR}
	s2 := TestSentence{"La vie est belle.",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true},
			{".", true}}, SENTENCE}

	s3 := TestSentence{" ",
		[]TestToken{
			{" ", true}}, SEPARATOR}

	s4 := TestSentence{"Avant d'y retourner aujourd'hui.",
		[]TestToken{
			{"Avant", true},
			{" ", true},
			{"d'", true},
			{"y", true},
			{" ", true},
			{"retourner", true},
			{" ", true},
			{"aujourd'hui", true},
			{".", true}}, SENTENCE}

	s5 := TestSentence{" ",
		[]TestToken{
			{" ", true},
		}, SEPARATOR}
	test1 := []TestSentence{s1, s2, s3, s4, s5}

	SubTestSentenceTokens(t, text1, test1, context)
}

func TestTokenizeSentenceParenthesis(t *testing.T) {

	context := GetTokenizeContext()
	text1 := "La vie est belle (C'est vrai.)."

	s1 := TestSentence{"La vie est belle (C'est vrai.).",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true},
			{" ", true},
			{"(", true},
			{"C'", true},
			{"est", true},
			{" ", true},
			{"vrai", true},
			{".", true},
			{")", true},
			{".", true}}, SENTENCE}

	test1 := []TestSentence{s1}

	SubTestSentenceTokens(t, text1, test1, context)

	text2 := "La vie est belle. (C'est vrai.)"

	s1 = TestSentence{"La vie est belle.",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true},
			{".", true}}, SENTENCE}

	s2 := TestSentence{" ",
		[]TestToken{
			{" ", true}}, SEPARATOR}

	s3 := TestSentence{"(C'est vrai.)",
		[]TestToken{
			{"(", true},
			{"C'", true},
			{"est", true},
			{" ", true},
			{"vrai", true},
			{".", true},
			{")", true}}, PARENTHESIS}

	test2 := []TestSentence{s1, s2, s3}

	SubTestSentenceTokens(t, text2, test2, context)
}

func TestTokenizeNonSentence(t *testing.T) {

	context := GetTokenizeContext()
	text1 := "  La vie est belle\t Avant d'y retourner aujourd'hui. "

	s1 := TestSentence{" ",
		[]TestToken{
			{" ", true},
			{" ", true},
		}, SEPARATOR}
	s2 := TestSentence{"La vie est belle.",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true}}, NOSENTENCE}

	s3 := TestSentence{" ",
		[]TestToken{
			{" ", true}}, SEPARATOR}

	s4 := TestSentence{"Avant d'y retourner aujourd'hui.",
		[]TestToken{
			{"Avant", true},
			{" ", true},
			{"d'", true},
			{"y", true},
			{" ", true},
			{"retourner", true},
			{" ", true},
			{"aujourd'hui", true},
			{".", true}}, SENTENCE}

	s5 := TestSentence{" ",
		[]TestToken{
			{" ", true},
		}, SEPARATOR}
	test1 := []TestSentence{s1, s2, s3, s4, s5}

	SubTestSentenceTokens(t, text1, test1, context)
}

func TestTokenizeNonSentenceWithURL(t *testing.T) {

	context := GetTokenizeContext()
	text1 := "  La vie est belle\t Avant d'aller sur lapresse.ca."

	s1 := TestSentence{" ",
		[]TestToken{
			{" ", true},
			{" ", true},
		}, SEPARATOR}
	s2 := TestSentence{"La vie est belle",
		[]TestToken{
			{"La", true},
			{" ", true},
			{"vie", true},
			{" ", true},
			{"est", true},
			{" ", true},
			{"belle", true}}, NOSENTENCE}

	s3 := TestSentence{" ",
		[]TestToken{
			{" ", true}}, SEPARATOR}

	s4 := TestSentence{"Avant d'aller sur lapresse.ca",
		[]TestToken{
			{"Avant", true},
			{" ", true},
			{"d'", true},
			{"aller", true},
			{" ", true},
			{"sur", true},
			{" ", true},
			{"lapresse.ca", true},
			{".", true}}, SENTENCE}

	s5 := TestSentence{" ",
		[]TestToken{
			{" ", true},
		}, SEPARATOR}
	test1 := []TestSentence{s1, s2, s3, s4, s5}

	SubTestSentenceTokens(t, text1, test1, context)
}
