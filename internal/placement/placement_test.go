package placement

import (
	"reflect"
	"testing"

	"github.com/beganov/gingonicserver/internal/card"
)

func TestOrderer(t *testing.T) {
	hands := [][]card.Card{
		{{Val: 5}, {Val: 6}, {Val: 7}},
		{{Val: 3}, {Val: 4}, {Val: 8}},
		{{Val: 9}, {Val: 10}, {Val: 11}},
	}
	openeds := [][]card.Card{
		{{Val: 1}}, {{Val: 2}}, {{Val: 3}},
	}
	idMap := map[int]int{2: 1, 3: 2}

	expectedHands := [][]card.Card{
		{{Val: 3}, {Val: 4}, {Val: 8}},
		{{Val: 5}, {Val: 6}, {Val: 7}},
		{{Val: 9}, {Val: 11}, {Val: 10}},
	}

	expectedOpeneds := [][]card.Card{
		{{Val: 2}}, {{Val: 1}}, {{Val: 3}},
	}
	expectedIdMap := map[int]int{2: 0, 3: 2}

	resHands, resOpeneds, resIdMap := Orderer(hands, openeds, idMap)

	if !reflect.DeepEqual(resHands, expectedHands) {
		t.Errorf("hands mismatch: got %v, want %v", resHands, expectedHands)
	}

	if !reflect.DeepEqual(resOpeneds, expectedOpeneds) {
		t.Errorf("openeds mismatch: got %v, want %v", resOpeneds, expectedOpeneds)
	}

	if !reflect.DeepEqual(resIdMap, expectedIdMap) {
		t.Errorf("idMap mismatch: got %v, want %v", resIdMap, expectedIdMap)
	}
}

func TestShufflePlayer(t *testing.T) {
	hands := [][]card.Card{
		{{Val: 5}, {Val: 6}, {Val: 7}},
		{{Val: 3}, {Val: 4}, {Val: 8}},
		{{Val: 9}, {Val: 10}, {Val: 11}},
	}
	openeds := [][]card.Card{
		{{Val: 1}}, {{Val: 2}}, {{Val: 3}},
	}
	idMap := map[int]int{1: 0, 2: 1}

	expectedHands := [][]card.Card{
		{{Val: 3}, {Val: 4}, {Val: 8}},
		{{Val: 9}, {Val: 10}, {Val: 11}},
		{{Val: 5}, {Val: 6}, {Val: 7}},
	}

	expectedOpeneds := [][]card.Card{
		{{Val: 2}}, {{Val: 3}}, {{Val: 1}},
	}

	resHands, resOpeneds := ShufflePlayer(hands, openeds, idMap)

	if !reflect.DeepEqual(resHands, expectedHands) {
		t.Errorf("hands mismatch: got %v, want %v", resHands, expectedHands)
	}

	if !reflect.DeepEqual(resOpeneds, expectedOpeneds) {
		t.Errorf("openeds mismatch: got %v, want %v", resOpeneds, expectedOpeneds)
	}
}

func TestNewPlacementArray(t *testing.T) {

	expectedArray := []int{0, 1, 2, 3, 4, 5}

	resArray := NewPlacementArray(6, 6)

	if !reflect.DeepEqual(resArray, expectedArray) {
		t.Errorf("hands mismatch: got %v, want %v", resArray, expectedArray)
	}
}
