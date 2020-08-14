package words

import (
	"encoding/xml"
	"io/ioutil"
	"strings"
)

type Pos struct {
	Name string `xml:"name,attr"`
}

type Feat struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Inflected struct {
	Form  string `xml:"form"`
	Feats []Feat `xml:"feat"`
}

type Entry struct {
	Lemma       string      `xml:"lemma"`
	Tag         Pos         `xml:"pos"`
	Feats       []Feat      `xml:"feat"`
	Inflections []Inflected `xml:"inflected"`
}

type WordList struct {
	XMLName xml.Name `xml:"dico"`
	Entries []Entry  `xml:"entry"`
}

func (words *WordList) Read(path string, lang byte, addWord func(form string, word *Word)(error)) (err error) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	err = xml.Unmarshal(bytes, words)
	if err != nil {
		return err
	}

	for _, entry := range words.Entries {
		for _, inflection := range entry.Inflections {
			word := Word{}
			v := inflection.GetVariant(&entry, lang)
			word.Variants = append(word.Variants, v)

			form := strings.Replace(inflection.Form, "\\-", "-", -1)
			err = addWord(form, &word)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (words *WordList) Write(path string) (err error) {
	bytes, err := xml.MarshalIndent(words, " ", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, bytes, 0666)
	return
}

func (inf *Inflected) GetVariant(entry *Entry, lang byte) WordVariant {
	variation := WordVariant{}
	variation.Gender = NOGENDER
	variation.Language = lang
	switch entry.Tag.Name {
	case "noun":
		variation.Tag = NOUN
	case "prep":
		variation.Tag = PREP
	case "adverb":
		variation.Tag = ADVERB
	case "verb":
		variation.Tag = VERB
	case "adj":
		variation.Tag = ADJ
	case "nominaldet":
		variation.Tag = NOMINALDET
	case "prefix":
		variation.Tag = PREFIX
	case "GNP":
		variation.Tag = GNP
	case "GNPX":
		variation.Tag = GNPX
	case "conjs":
		variation.Tag = CONJS
	case "X":
		variation.Tag = X
	case "PCDN3":
		variation.Tag = PCDN3
	case "intj":
		variation.Tag = INTJ
	case "GN":
		variation.Tag = GN
	case "conjc":
		variation.Tag = CONJC
	case "det":
		variation.Tag = DET
	case "pronoun":
		variation.Tag = PRONOUN
	case "prepadj":
		variation.Tag = PREPADJ
	case "prepdet":
		variation.Tag = PREPDET
	case "preppro":
		variation.Tag = PREPPRO
	case "PRON":
		variation.Tag = PRON
	case "XI":
		variation.Tag = XI
	case "PART":
		variation.Tag = PART
	case "PRED":
		variation.Tag = PRED
	case "conj":
		variation.Tag = CONJ
	case " Clark.N":
		variation.Tag = CLARKN
	case "VA":
		variation.Tag = VA
	case "NA":
		variation.Tag = NA
	case " Ghery.N":
		variation.Tag = GHERYN
	case " G\\. Wells.N":
		variation.Tag = GWELLSN
	case " Mayer.N":
		variation.Tag = MAYERN
	case "NES":
		variation.Tag = NES
	case " Maugham.N":
		variation.Tag = MAUGHAMN
	case "ADVA":
		variation.Tag = ADVA
	case " C\\. Fields.N":
		variation.Tag = CFIELDSN
	case " Hart.N":
		variation.Tag = HARTN
	case "abbr":
		variation.Tag = ABBR
	case "space":
		variation.Tag = SPACE
	case "comma":
		variation.Tag = COMMA
	case "dot":
		variation.Tag = DOT
	case "colon":
		variation.Tag = COLON
	case "questionmark":
		variation.Tag = QUESTIONMARK
	case "exclamationmark":
		variation.Tag = EXCLAMATIONMARK
	case "beginparenthesis":
		variation.Tag = BEGINPARENTHESIS
	case "endparenthesis":
		variation.Tag = ENDPARENTHESIS
	case "apos":
		variation.Tag = APOS
	case "quotationmark":
		variation.Tag = QUOTATIONMARK
	case "beginquotation":
		variation.Tag = BEGINQUOTATION
	case "endquotation":
		variation.Tag = ENDQUOTATION
	case "dash":
		variation.Tag = DASH
	case "tab":
		variation.Tag = TAB
	case "":
		variation.Tag = 0
	default:
		panic("'" + entry.Tag.Name + "'")
	}
	for _, feat := range entry.Feats {
		switch feat.Name {
		case "proper":
			switch feat.Value {
			case "true":
				variation.Flags |= PROPER
			default:
				panic(feat.Value)
			}
		case "subcat":
			switch feat.Value {
			case "human":
				variation.Subcat = HUMAN
			case "animal":
				variation.Subcat = ANIMAL
			case "concret":
				variation.Subcat = CONCRET
			case "abstract":
				variation.Subcat = ABSTRACT
			case "unit":
				variation.Subcat = UNIT
			case "indefinite":
				variation.Subcat = INDEFINITE
			case "temporal":
				variation.Subcat = TEMPORAL
			case "demonstrative":
				variation.Subcat = DEMONSTRATIVE
			default:
				panic(feat.Value)
			}
		case "compound":
			switch feat.Value {
			case "comp":
				variation.Flags |= COMPOUND
			default:
				panic(feat.Value)
			}
		case "coll":
			switch feat.Value {
			case "true":
				variation.Flags |= COLL
			default:
				panic(feat.Value)
			}
		case "postpos":
			switch feat.Value {
			case "true":
				variation.Flags |= POSTPOS
			default:
				panic(feat.Value)
			}
		case "collective":
			switch feat.Value {
			case "true":
				variation.Flags |= COLLECTIVE
			default:
				panic(feat.Value)
			}
		case "procat":
			switch feat.Value {
			case "demonstrative":
				variation.Flags |= PROCATDEMONSTRATIVE
			default:
				panic(feat.Value)
			}
		default:
			panic(feat.Name)
		}
	}
	for _, feat := range inf.Feats {
		switch feat.Name {
		case "gender":
			switch feat.Value {
			case "masculine":
				variation.Gender = MALE
			case "feminine":
				variation.Gender = FEMALE
			default:
				panic(feat.Value)
			}
		case "number":
			switch feat.Value {
			case "singular":
				variation.Number = SINGULAR
			case "plural":
				variation.Number = PLURAL
			default:
				panic(feat.Value)
			}
		case "person":
			switch feat.Value {
			case "1":
				variation.Person = 1
			case "2":
				variation.Person = 2
			case "3":
				variation.Person = 3
			default:
				panic(feat.Value)
			}
		case "tense":
			switch feat.Value {
			case "ind":
				variation.Tense = IND
			case "gerondif":
				variation.Tense = GERONDIF
			case "subj":
				variation.Tense = SUBJ
			case "ppast":
				variation.Tense = PPAST
			case "imp":
				variation.Tense = IMP
			case "cond":
				variation.Tense = COND
			case "inf":
				variation.Tense = INF
			default:
				panic(feat.Value)
			}
		default:
			panic(feat.Name)
		}
	}
	return variation
}
