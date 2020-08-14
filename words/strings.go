package words

import (
	"unicode/utf8"
)

func ReplaceAtPos(word string, pos int, r rune, w int) string {
	var start string
	var middle string
	var end string

	if pos > 0 {
		start = word[0:pos]
	}

	buffer := []byte{0, 0, 0, 0}
	utf8.EncodeRune(buffer, r)
	middle = string(buffer[0:w])

	if pos+w < len(word) {
		end = word[pos+w:]
	}

	return start + middle + end
}

func CountError(correct string, alt string) int {
	var errors int
	for i := 0; i < len(correct); i++ {
		if correct[i] != alt[i] {
			errors++
		}
	}
	return errors
}
