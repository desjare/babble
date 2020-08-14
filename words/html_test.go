package words

/*
import "testing"
import "io/ioutil"
import "sort"

func DoTestHtml(t *testing.T, context* TokenizeContext) {
	bytes, err := ioutil.ReadFile("techno.htm")
	if err != nil {
		t.Errorf("not found")
		return
	}

	contentBytes, _ := ParseHtml(bytes)

	content := string(contentBytes)
	tokens := Tokenize(content, context)

	sortedwords := make(sort.StringSlice,1)

	for _, tok := range tokens {
		if len(tok.Words) == 0 && !tok.IsValid() {
			i := sortedwords.Search(tok.Content(content))
			if i < len(sortedwords) && sortedwords[i] == tok.Content(content) {
				continue
			} else {
				sortedwords = append(sortedwords, tok.Content(content))
				sortedwords.Sort()
				t.Errorf("'%s' %s \n", tok.Content(content), tok.String())

			}
		}
	}

	wordList := WordList{}
	for _, word :=range(sortedwords) {
		entry := Entry{}
		entry.Lemma = []byte(word)

		inf := Inflected{}
		inf.Form = word

		entry.Inflections = append(entry.Inflections, inf)
		wordList.Entries = append(wordList.Entries, entry)
	}
	wordList.Write("dela-proper-u8.dic.xml")

}

*/
