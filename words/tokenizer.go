package words

import "fmt"
import "unicode"
import "unicode/utf8"
import "regexp"
import "strings"

// sentence type
const (
	SENTENCE    = 1
	NOSENTENCE  = 2
	SEPARATOR   = 3
	PARENTHESIS = 4
)

type TokenizeContext struct {
	dict    *Dictionary
	rtime   *regexp.Regexp
	rtemp   *regexp.Regexp
	rnumber *regexp.Regexp
	rroman  *regexp.Regexp
	rhex    *regexp.Regexp
	rdate   *regexp.Regexp
	rurl    *regexp.Regexp
}

type Token struct {
	Pos      []int
	Word     *Word
	IsNumber bool
	IsTime   bool
	IsDate   bool
	IsTemp   bool
	IsURL    bool
	IsUpper  bool
}

type TokenSentence struct {
	Tokens []Token
	Type   byte
}

func TokenizeNewContext() (context *TokenizeContext, err error) {
	context = new(TokenizeContext)

	// compile regexp
	context.rtime, err = regexp.Compile("[0-6]{1,2}(h|:)[0-6][0-9]")
	if err != nil {
		panic(err)
		return
	}

	context.rtemp, err = regexp.Compile("-?[0-9]+°(C|F|K)")
	if err != nil {
		panic(err)
		return
	}

	context.rnumber, err = regexp.Compile("[0-9]+e?")
	if err != nil {
		panic(err)
		return
	}

	context.rroman, err = regexp.Compile("[MDCLXVI]+e?")
	if err != nil {
		panic(err)
		return
	}

	context.rhex, err = regexp.Compile("(0x)?[0-9a-fA-F]+")
	if err != nil {
		panic(err)
		return
	}

	context.rdate, err = regexp.Compile("[0-9]{2,4}/[0-9]{2}/[0-9]{2,4}")
	if err != nil {
		panic(err)
		return
	}
	context.rurl, err = regexp.Compile(`((http|https|ftp)://)?([0-9a-z_-]+\x2E)+(aero|asia|biz|cat|com|coop|edu|gov|info|int|jobs|mil|mobi|museum|name|net|org|pro|tel|travel|ac|ad|ae|af|ag|ai|al|am|an|ao|aq|ar|as|at|au|aw|ax|az|ba|bb|bd|be|bf|bg|bh|bi|bj|bm|bn|bo|br|bs|bt|bv|bw|by|bz|ca|cc|cd|cf|cg|ch|ci|ck|cl|cm|cn|co|cr|cu|cv|cx|cy|cz|cz|de|dj|dk|dm|do|dz|ec|ee|eg|er|es|et|eu|fi|fj|fk|fm|fo|fr|ga|gb|gd|ge|gf|gg|gh|gi|gl|gm|gn|gp|gq|gr|gs|gt|gu|gw|gy|hk|hm|hn|hr|ht|hu|id|ie|il|im|in|io|iq|ir|is|it|je|jm|jo|jp|ke|kg|kh|ki|km|kn|kp|kr|kw|ky|kz|la|lb|lc|li|lk|lr|ls|lt|lu|lv|ly|ma|mc|md|me|mg|mh|mk|ml|mn|mn|mo|mp|mr|ms|mt|mu|mv|mw|mx|my|mz|na|nc|ne|nf|ng|ni|nl|no|np|nr|nu|nz|nom|pa|pe|pf|pg|ph|pk|pl|pm|pn|pr|ps|pt|pw|py|qa|re|ra|rs|ru|rw|sa|sb|sc|sd|se|sg|sh|si|sj|sj|sk|sl|sm|sn|so|sr|st|su|sv|sy|sz|tc|td|tf|tg|th|tj|tk|tl|tm|tn|to|tp|tr|tt|tv|tw|tz|ua|ug|uk|us|uy|uz|va|vc|ve|vg|vi|vn|vu|wf|ws|ye|yt|yu|za|zm|zw|arpa)(:[0-9]+)?(/[[0-9a-z\/\?\=\#\(\)_\-\.]+)?`)
	if err != nil {
		panic(err)
		return
	}

	// load dicts
	context.dict = new(Dictionary)
	err = context.dict.ReadBinary("lm/lm.bin")
	if err != nil {
		return
	}
	return
}

func(context *TokenizeContext) GetDictionary() *Dictionary {
	return context.dict
}

func (t *Token) Content(content string) string {
	return content[t.Pos[0]:t.Pos[1]]
}

func (t *Token) IsValid() bool {
	return t.Word != nil || t.IsNumber || t.IsTime || t.IsTemp || t.IsURL || t.IsDate
}

func (t *Token) String() string {
	s := fmt.Sprintf("[%d-%d] %s %s ", t.Pos[0], t.Pos[1], t.Word)
	if t.Word != nil {
		s += t.Word.Description()
	}
	if t.IsNumber {
		s += "IsNumber "
	}
	if t.IsTime {
		s += "IsTemp "
	}
	if t.IsDate {
		s += "IsDate "
	}
	if t.IsTemp {
		s += "IsTemp "
	}
	if t.IsURL {
		s += "IsURL "
	}
	if t.IsUpper {
		s += "IsUpper "
	}
	return s
}

func (t *TokenSentence) Content(content string) string {
	s := "["
	for _, tok := range t.Tokens {
		s += tok.Content(content)
	}
	s += "]"
	s += fmt.Sprintf("type %d", t.Type)
	return s
}

func TokenizeFindWord(s string, context *TokenizeContext) (word *Word, foundPath bool) {
	r, _ := utf8.DecodeRuneInString(s)
	isUpper := unicode.IsUpper(r)

	word, foundPath = context.dict.FindWord(s)

	// match lower case version of word
	if word == nil && isUpper {
		word, foundPath = context.dict.FindWord(TokenizeToLower(s))
	}
	return word, foundPath
}

func TokenizeToLower(s string) string {
	return strings.ToLower(s)
}

func TokenizeIsNumber(s string, context *TokenizeContext) bool {
	return TokenizeMatchOnly(s, context.rnumber) || TokenizeMatchOnly(s, context.rroman) || TokenizeMatchOnly(s, context.rhex)
}

func TokenizeIsTag(t Token, tag byte) bool {
	word := t.Word
	if word == nil {
		return false
	}
	return word.Tagged(tag)
}

func TokenizeIsDot(t Token) bool {
	word := t.Word
	if word == nil {
		return false
	}
	return word.Tagged(DOT) || word.Tagged(EXCLAMATIONMARK) || word.Tagged(QUESTIONMARK)
}

func TokenizeIsSpace(t Token) bool {
	return TokenizeIsTag(t, SPACE)
}

func TokenizeIsEnd(t Token) bool {
	return TokenizeIsDot(t) || TokenizeIsTab(t)
}

func TokenizeIsBeginParenthesis(t Token) bool {
	return TokenizeIsTag(t, BEGINPARENTHESIS)
}

func TokenizeIsEndParenthesis(t Token) bool {
	return TokenizeIsTag(t, ENDPARENTHESIS)
}

func TokenizeIsTab(t Token) bool {
	return TokenizeIsTag(t, TAB)
}

func TokenizeMatchOnly(s string, r *regexp.Regexp) bool {
	indexes := r.FindAllStringIndex(s, -1)
	if len(indexes) != 1 {
		return false
	}

	if indexes[0][0] != 0 || indexes[0][1] != len(s) {
		return false
	}
	return true
}

func TokenizeIsTemp(s string, context *TokenizeContext) bool {
	return TokenizeMatchOnly(s, context.rtemp)
}

func TokenizeIsTime(s string, context *TokenizeContext) bool {
	return TokenizeMatchOnly(s, context.rtime)
}

func TokenizeIsDate(s string, context *TokenizeContext) bool {
	return TokenizeMatchOnly(s, context.rdate)
}

func TokenizeIsURL(s string, context *TokenizeContext) bool {
	// twitter
	if strings.HasPrefix(s, "@") {
		return true
	}

	return TokenizeMatchOnly(s, context.rurl)
}

func TokenizeIsWord(s string) bool {
	var r rune
	for i, l := 0, 0; i < len(s); i += l {
		r, l = utf8.DecodeRuneInString(s[i:])
		if !unicode.IsLetter(r) && r != '-' {
			return false
		}
		if unicode.IsUpper(r) {
			return false
		}
	}
	return true
}

func TokenizeIsProperNoun(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	isUpper := unicode.IsUpper(r) && unicode.IsLetter(r)
	allUpper := true

	if !isUpper {
		return false
	}
	for i, l := 0, 0; i < len(s); i += l {
		r, l = utf8.DecodeRuneInString(s[i:])
		if !unicode.IsLetter(r) && r != '-' {
			return false
		}
		if !unicode.IsUpper(r) && r != '-' {
			allUpper = false
		}
	}
	return !allUpper
}

func TokenizeBuildToken(content string, searchPath *bool, start int, end int, context *TokenizeContext) *Token {
	var word *Word
	var foundPath bool

	if *searchPath {
		word, foundPath = TokenizeFindWord(content[start:end], context)
		if !foundPath {
			*searchPath = false
		}
	}

	r, _ := utf8.DecodeRuneInString(content[start:end])
	isUpper := unicode.IsUpper(r)
	isNumber := false
	isTime := false
	isDate := false
	isTemp := false
	isURL := false
	if word == nil {
		isNumber = TokenizeIsNumber(content[start:end], context)
		if !isNumber {
			isTime = TokenizeIsTime(content[start:end], context)
		}
		if !isNumber && !isTime {
			isDate = TokenizeIsDate(content[start:end], context)
		}
		if !isNumber && !isTime && !isDate {
			isTemp = TokenizeIsTemp(content[start:end], context)
		}
		if !isNumber && !isTime && !isDate && !isTemp {
			isURL = TokenizeIsURL(content[start:end], context)
		}
	}
	return &Token{[]int{start, end},
		word,
		isNumber,
		isTime,
		isDate,
		isTemp,
		isURL,
		isUpper}
}

func TokenizeAddToken(content string, start int, end int, intoks []Token, context *TokenizeContext) (tokens []Token) {
	if start == end {
		return intoks
	}
	searchPath := true
	token := TokenizeBuildToken(content, &searchPath, start, end, context)
	tokens = append(intoks, *token)

	return
}

func TokenizeCompoundToken(content string, intoks []Token, context *TokenizeContext) (tokens []Token) {
	for i := 0; i < len(intoks); {
		var startPos int
		var token *Token
		var endToken int

		startPos = intoks[i].Pos[0]
		r := rune(content[intoks[i].Pos[0]])
		isComp := !(unicode.IsSpace(r) || unicode.IsDigit(r))
		searchPath := true

		for j := i + 1; j < len(intoks) && isComp && searchPath; j++ {
			if j-i <= context.dict.MaxTokens {
				compoundToken := TokenizeBuildToken(content, &searchPath, startPos, intoks[j].Pos[1], context)
				if compoundToken.IsValid() {
					token = compoundToken
					endToken = j
				}
			}
		}
		if token == nil {
			tokens = append(tokens, intoks[i])
			i = i + 1
		} else {
			if len(tokens) == 0 {
				tokens = []Token{}
			}
			tokens = append(tokens, *token)
			i = endToken + 1
		}
	}
	return
}

func Tokenize(content string, context *TokenizeContext, compound bool) (tokens []Token) {
	var r rune
	for i, j, w := 0, 0, 0; i < len(content); i += w {
		r, w = utf8.DecodeRuneInString(content[i:])
		isAnd := r == '&'
		if (unicode.IsPunct(r) || unicode.IsSpace(r)) && !isAnd {
			if i > j {
				tokens = TokenizeAddToken(content, j, i, tokens, context)
			}
			tokens = TokenizeAddToken(content, i, i+w, tokens, context)
			j = i + w
		}
		if i == len(content)-1 {
			tokens = TokenizeAddToken(content, j, i+w, tokens, context)
		}
	}
	if compound {
		tokens = TokenizeCompoundToken(content, tokens, context)
	}
	return
}

func TokenizeGetSentenceType(i int, j int, tokens []Token) byte {
	if TokenizeIsBeginParenthesis(tokens[i]) && TokenizeIsEndParenthesis(tokens[j]) {
		return PARENTHESIS
	}
	return NOSENTENCE
}

func TokenizeSentence(content string, context *TokenizeContext) (s []TokenSentence) {

	var j int

	tokens := Tokenize(content, context, true)

	for i := 0; i < len(tokens); i++ {

		t := tokens[i]

		// dot separated sentence
		if TokenizeIsDot(t) && !TokenizeIsBeginParenthesis(tokens[j]) {
			s = append(s, TokenSentence{tokens[j : i+1], SENTENCE})
			j = i + 1
			continue
		}

		// separators [  A] [.  A] [\t  A]
		if TokenizeIsSpace(t) && (i == 0 || i > 0 && TokenizeIsEnd(tokens[i-1])) {
			var k int
			for k = i; k < len(tokens); k++ {
				if !TokenizeIsSpace(tokens[k]) {
					s = append(s, TokenSentence{tokens[j:k], SEPARATOR})
					i = k
					j = k
					break
				}
			}
			if k == len(tokens) {
				s = append(s, TokenSentence{tokens[j:len(tokens)], SEPARATOR})
				i = k // reach end
				j = k + 1
			}
			continue
		}

		// ()
		if TokenizeIsBeginParenthesis(t) {
			var k int
			for k = i; k < len(tokens); k++ {
				if TokenizeIsEndParenthesis(tokens[k]) {
					i = k
					continue
				}
			}
		}

		// non sentence separator
		if TokenizeIsTab(t) {
			if i-j > 0 {
				s = append(s, TokenSentence{tokens[j:i], NOSENTENCE})
			}
			j = i + 1 // skip sep
			continue
		}

		// end tokens
		if i == len(tokens)-1 && j < len(tokens) && i-j > 0 {
			t := TokenizeGetSentenceType(j, len(tokens)-1, tokens)
			s = append(s, TokenSentence{tokens[j:len(tokens)], t})
		}
	}
	return
}

func TokenizePrintTokens(content string, tokens []Token) {
	for _, t := range tokens {
		fmt.Printf("'%s' word %s\n", content[t.Pos[0]:t.Pos[1]], t.Word)
	}
	return
}
