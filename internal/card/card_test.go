package card

import (
	"reflect"
	"testing"
)

func TestIsSpecial(t *testing.T) {

	res := isSpecial(6)

	if !reflect.DeepEqual(res, false) {
		t.Errorf("isSpecial mismatch: got %v, want %v", res, false)
	}
}

func TestDecksUpdate(t *testing.T) {

	// resArray := isSpecial(6)

	// if !reflect.DeepEqual(resArray, false) {
	// 	t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	// }
}

func TestCardDelete(t *testing.T) {

	// resArray := isSpecial(6)

	// if !reflect.DeepEqual(resArray, false) {
	// 	t.Errorf("isSpecial mismatch: got %v, want %v", resArray, false)
	// }
}

func TestSortCard(t *testing.T) {

	// resArray := isSpecial(6)

	// if !reflect.DeepEqual(resArray, false) {
	// 	t.Errorf("hands mismatch: got %v, want %v", resArray, false)
	// }
}
