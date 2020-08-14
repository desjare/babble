package words

import (
	"testing"
)

func TestReplaceAtPos(t *testing.T) {

	corrected := ReplaceAtPos("tast", 1, 'e', 1)
	if corrected != "test" {
		t.Errorf("could not replace tast")
	}

	corrected = ReplaceAtPos("tess", 3, 't', 1)
	if corrected != "test" {
		t.Errorf("could not replace tess")
	}

	corrected = ReplaceAtPos("aest", 0, 't', 1)
	if corrected != "test" {
		t.Errorf("could not replace aest")
	}
}
