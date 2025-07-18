package card

import (
	"sort"
)

type Card struct {
	Id  int `json:"id"`
	Val int `json:"val"`
}

func SortCard(Hands []Card) {
	sort.Slice(Hands, func(k, l int) bool {
		if isSpecial(Hands[k].Val) {
			return false
		}
		if isSpecial(Hands[l].Val) {
			return true
		}
		return Hands[k].Val < Hands[l].Val
	})
}

func cardDelete(arr []Card, index int) []Card {
	arr[index] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return arr
}

func DecksUpdate(arr, arr2 []Card, index int) ([]Card, []Card) {
	arr = append(arr, arr2[index])
	arr2 = cardDelete(arr2, index)
	return arr, arr2
}

func isSpecial(Val int) bool {
	return Val == 2 || Val == 10 || Val == 0
}
