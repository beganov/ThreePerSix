package card

import (
	"reflect"
	"testing"
)

func TestIsSpecial(t *testing.T) {

	resArray := isSpecial(6)

	if !reflect.DeepEqual(resArray, false) {
		t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	}
}

func TestDecksUpdate(t *testing.T) {

	resArray := isSpecial(6)

	if !reflect.DeepEqual(resArray, false) {
		t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	}
}

func TestCardDelete(t *testing.T) {

	resArray := isSpecial(6)

	if !reflect.DeepEqual(resArray, false) {
		t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	}
}

func TestSortCard(t *testing.T) {

	resArray := isSpecial(6)

	if !reflect.DeepEqual(resArray, false) {
		t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	}
}
