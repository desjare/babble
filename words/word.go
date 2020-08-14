package words

import (
	"encoding/binary"
	"fmt"
	"io"
)

// tags
const (
	NOUN       = 1
	PREP       = 2
	ADVERB     = 3
	VERB       = 4
	ADJ        = 5
	NOMINALDET = 6
	PREFIX     = 7
	GNP        = 8
	GNPX       = 9
	CONJS      = 10
	CONJ       = 11
	GN         = 12
	CONJC      = 13
	PRONOUN    = 14
	PREPADJ    = 15
	PREPDET    = 16
	PREPPRO    = 17
	INTJ       = 18
	DET        = 19
	PRON       = 20
	VA         = 21
	NA         = 22
	CLARKN     = 23
	GHERYN     = 24
	GWELLSN    = 25
	MAYERN     = 26
	NES        = 27
	MAUGHAMN   = 28
	ADVA       = 29
	CFIELDSN   = 30
	X          = 31
	PCDN3      = 32
	XI         = 33
	PART       = 34
	PRED       = 35
	HARTN      = 36
	ABBR       = 37
	// punctuation
	BEGIN_PUNCT = 38
	TAB              = 38
	QUOTATIONMARK    = 39
	BEGINQUOTATION   = 40
	ENDQUOTATION     = 41
	APOS             = 42
	SLASH            = 43
	DASH             = 44
	GREATHERTHAN     = 45
	SMALLERTHAN      = 46
	SPACE            = 47
	COMMA            = 48
	SEMICOLON        = 49
	DOT              = 50
	COLON            = 51
	EXCLAMATIONMARK  = 52
	QUESTIONMARK     = 53
	BEGINPARENTHESIS = 54
	ENDPARENTHESIS   = 55
	BEGINBRACKET     = 56
	ENDBRACKET       = 57
	// symbols
	NUMBERSIGN    = 58
	DOLLARSIGN    = 59
	COPYRIGHTSIGN = 60
	ATSIGN        = 61
	NBNS          = 62
	AND           = 63
	OR            = 64
	EOL           = 65
	CR            = 66
)

// gender
const (
	NOGENDER = 0
	MALE     = 1
	FEMALE   = 2
)

// verb tense
const (
	IND      = 1
	GERONDIF = 2
	SUBJ     = 3
	PPAST    = 4
	IMP      = 5
	COND     = 6
	INF      = 7
)

// number
const (
	SINGULAR = 0x1
	PLURAL   = 0x2
)

// subcat
const (
	HUMAN         = 1
	ANIMAL        = 2
	CONCRET       = 3
	ABSTRACT      = 4
	UNIT          = 5
	INDEFINITE    = 6
	TEMPORAL      = 7
	DEMONSTRATIVE = 8
)

// languages
const (
	ENGLISH = 1
	FRENCH  = 2
)

// flags
const (
	PROPER              = 0x1
	SUBCAT              = 0x2
	COMPOUND            = 0x4
	COLL                = 0x8
	POSTPOS             = 0x10
	COLLECTIVE          = 0x20
	PROCATDEMONSTRATIVE = 0x40
)

type WordVariant struct {
	Tag      byte
	Language byte

	Flags  byte
	Subcat byte

	Gender byte // male or female
	Number byte // singular or plural

	Person byte // verb person 1 or 2 or 3
	Tense  byte // verb tense
}

type Word struct {
	LastLetter *WordLetter
	Variants   []WordVariant
}

func (word *Word) String() string {
	return word.LastLetter.GetWord()
}

func (word* Word) Description() string {
	s := "[" + word.String()

	for _, v := range(word.Variants) {
		s += fmt.Sprintf("[tag %d]", v.Tag)
	}

	s += "]"

	return s
}

func (word *Word) IsPunct() bool {
	if len(word.Variants) != 1 {
		return false
	}
	return word.Variants[0].Tag >= BEGIN_PUNCT
}

func (word *Word) Language(lang byte) bool {
	for _, v := range word.Variants {
		if v.Language == lang {
			return true
		}
	}
	return false
}

func (word *Word) Tagged(tag byte) bool {
	for _, v := range word.Variants {
		if v.Tag == tag {
			return true
		}
	}
	return false
}

func (rhs *WordVariant) Equals(lhs *WordVariant) bool {
	if rhs.Tag != lhs.Tag {
		return false
	}
	if rhs.Language != lhs.Language {
		return false
	}
	if rhs.Flags != lhs.Flags {
		return false
	}
	if rhs.Subcat != lhs.Subcat {
		return false
	}
	if rhs.Person != lhs.Person {
		return false
	}
	if rhs.Gender != lhs.Gender {
		return false
	}
	if rhs.Number != lhs.Number {
		return false
	}
	if rhs.Tense != lhs.Tense {
		return false
	}
	return true
}

func (v *WordVariant) Filter(w *Word) bool {

	variants := w.VariantsByTag(v.Tag, v.Language)

	for _, variant := range variants {
		match := true
		if v.Tag == VERB {
			if v.Person != 0 && v.Person != variant.Person {
				match = false
			}
			if v.Tense != 0 && v.Tense != variant.Tense {
				match = false
			}
		}

		if v.Number != 0 && v.Number != variant.Number {
			match = false
		}
		if v.Gender != 0 && v.Gender != variant.Gender {
			match = false
		}
		if v.Subcat == 0 && v.Subcat != variant.Subcat {
			match = false
		}

		if match {
			return true
		}
	}
	return false

}

func (word *Word) AddVariants(variants []WordVariant) {
	for _, v := range variants {
		found := false
		for _, w := range word.Variants {
			if w.Equals(&v) {
				found = true
				break
			}
		}
		if !found {
			word.Variants = append(word.Variants, v)
		}
	}
}

func (word *Word) VariantsByTag(tag byte, lang byte) []WordVariant {
	var variants []WordVariant

	for _, v := range word.Variants {
		if v.Tag == tag && v.Language == lang {
			variants = append(variants, v)
		}
	}
	return variants
}

func (word *Word) Write(buf io.Writer) {
	var l int32

	letters := word.LastLetter.GetWord()
	l = int32(len(letters))
	binary.Write(buf, binary.LittleEndian, l)
	binary.Write(buf, binary.LittleEndian, []byte(letters))

	v := int32(len(word.Variants))
	binary.Write(buf, binary.LittleEndian, v)
	for _, variant := range word.Variants {
		binary.Write(buf, binary.LittleEndian, variant.Tag)
		binary.Write(buf, binary.LittleEndian, variant.Language)
		binary.Write(buf, binary.LittleEndian, variant.Flags)
		binary.Write(buf, binary.LittleEndian, variant.Subcat)
		binary.Write(buf, binary.LittleEndian, variant.Person)
		binary.Write(buf, binary.LittleEndian, variant.Gender)
		binary.Write(buf, binary.LittleEndian, variant.Number)
		binary.Write(buf, binary.LittleEndian, variant.Tense)
	}
	return
}

func (word *Word) Read(buf io.Reader) (s string, err error) {
	var l int32
	err = binary.Read(buf, binary.LittleEndian, &l)
	if err != nil {
		return
	}
	letters := make([]byte, l)
	err = binary.Read(buf, binary.LittleEndian, &letters)
	if err != nil {
		return
	}
	s = string(letters)

	var v int32
	err = binary.Read(buf, binary.LittleEndian, &v)
	if err != nil {
		return
	}
	var i int32
	for i = 0; i < v; i++ {
		word.Variants = append(word.Variants, WordVariant{})
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Tag)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Language)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Flags)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Subcat)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Person)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Gender)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Number)
		if err != nil {
			return
		}
		err = binary.Read(buf, binary.LittleEndian, &word.Variants[i].Tense)
		if err != nil {
			return
		}
	}
	return
}
