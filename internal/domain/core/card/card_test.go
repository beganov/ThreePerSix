package card

import (
	"reflect"
	"testing"
)

func TestIsSpecial(t *testing.T) {
	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"Val is 2", 2, true},
		{"Val is 10", 10, true},
		{"Val is 0", 0, true},
		{"Val is 5", 5, false},
		{"Val is 13", 13, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSpecial(tt.input)
			if result != tt.expected {
				t.Errorf("isSpecial(%d) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDecksUpdate(t *testing.T) {

	tests := []struct {
		name                string
		firstArray          []Card
		secondArray         []Card
		index               int
		expectedFirstArray  []Card
		expectedSecondArray []Card
	}{
		{"Delete first Card", []Card{{Val: 2}, {Val: 4}}, []Card{{Val: 2}, {Val: 4}}, 0, []Card{{Val: 2}, {Val: 4}, {Val: 2}}, []Card{{Val: 4}}},
		{"Delete last Card", []Card{{Val: 2}, {Val: 4}}, []Card{{Val: 2}, {Val: 4}}, 1, []Card{{Val: 2}, {Val: 4}, {Val: 4}}, []Card{{Val: 2}}},
		{"Delete last Card", []Card{{Val: 2}, {Val: 4}}, []Card{{Val: 2}}, 0, []Card{{Val: 2}, {Val: 4}, {Val: 2}}, []Card{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input1 := tt.firstArray
			input2 := tt.secondArray
			tt.firstArray, tt.secondArray = DecksUpdate(tt.firstArray, tt.secondArray, tt.index)
			if !reflect.DeepEqual(tt.firstArray, tt.expectedFirstArray) {
				t.Errorf("DeckUpdate(%v,%v,%v) - Growing Array %v, want %v", input1, input2, tt.index, tt.firstArray, tt.expectedFirstArray)
			}
			if !reflect.DeepEqual(tt.secondArray, tt.expectedSecondArray) {
				t.Errorf("DeckUpdate(%v,%v,%v) - Reduced Array %v, want %v", input1, input2, tt.index, tt.secondArray, tt.expectedSecondArray)
			}
		})
	}
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
