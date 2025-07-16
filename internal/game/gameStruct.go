package game

import "sync"

const MaxValue = 15
const MinValue = 2
const cardQuantity = 54
const packSize = 3

type GameState struct {
	sync.Mutex
	Deck           []Card      `json:"deck,omitempty"`
	Out            []Card      `json:"out,omitempty"`
	Hands          [][]Card    `json:"hands,omitempty"`
	Openeds        [][]Card    `json:"openeds,omitempty"`
	Closeds        [][]Card    `json:"closeds,omitempty"`
	MaxPlayerCount int         `json:"maxPlayerCount,omitempty"`
	Iamind         map[int]int `json:"iamind,omitempty"`
	Alsoiamind     map[int]int `json:"alsoIamind,omitempty"`
	Iamindalso     map[int]int `json:"iamindalso,omitempty"`
	ch             map[int]chan int
	//Iamind int `json:"iamind,omitempty"`
}

type Card struct {
	Id  int `json:"id"`
	Val int `json:"val"`
}
