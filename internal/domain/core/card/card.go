package card

import (
	"sort"
)

type Card struct {
	Id  int `json:"id"`
	Val int `json:"val"`
}

// Сортировка колоды (По возрастанию, но карты с спец. свойствами ценнее)
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

// Удаление карты из набора карт
func cardDelete(arr []Card, index int) []Card {
	arr[index] = arr[len(arr)-1]
	arr = arr[:len(arr)-1]
	return arr
}

// Перенос карты из одного набора карт в другой
func DecksUpdate(arr, arr2 []Card, index int) ([]Card, []Card) {
	arr = append(arr, arr2[index])
	arr2 = cardDelete(arr2, index)
	return arr, arr2
}

// Проверка на то, есть ли у карты особые свойства
func isSpecial(Val int) bool {
	return Val == 2 || Val == 10 || Val == 0
}
