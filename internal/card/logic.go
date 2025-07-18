package card

import "github.com/beganov/gingonicserver/internal/gameConst"

func GiveCardLogic(Hands, Out []Card, cardState, i int, iamindFlag, flag, istake bool, ch <-chan int) ([]Card, []Card, int, bool, bool) {
	if len(Hands) != 0 {
		if cardState > gameConst.TakedCardState {
			Out, Hands, flag = ReGiveCard(Out, Hands, cardState, iamindFlag, ch)
			istake = !flag
		}
		if cardState == gameConst.StartCardState {
			Out, Hands, istake, cardState = GiveCard(Out, Hands, iamindFlag, ch)
			flag = !istake
		}
		if len(Out) > 0 {
			if Out[len(Out)-1].Val == 0 || Out[len(Out)-1].Val == 10 {
				Out = Out[:0]
				flag = false
				cardState = gameConst.StartCardState
			}
		}
		if len(Out) >= 4 && Out[len(Out)-1].Val == Out[len(Out)-2].Val && Out[len(Out)-2].Val == Out[len(Out)-3].Val && Out[len(Out)-3].Val == Out[len(Out)-4].Val {
			Out = Out[:0]
			flag = false
			cardState = gameConst.StartCardState
		}
	}
	return Hands, Out, cardState, flag, istake
}

func TakeCard(Deck, Hands, Openeds, Closeds []Card, istake bool) ([]Card, []Card, []Card, []Card) {
	if len(Deck) == 0 && len(Openeds) == 0 && len(Hands) == 0 && len(Closeds) != 0 {
		Hands, Closeds = DecksUpdate(Hands, Closeds, 0)
	}
	if len(Deck) == 0 && len(Hands) == 0 {
		Hands = Openeds
		Openeds = Openeds[:0]
	}
	if len(Hands) < gameConst.PackSize && len(Deck) > 0 && istake {
		Hands, Deck = DecksUpdate(Hands, Deck, 0)
	}
	return Deck, Hands, Openeds, Closeds
}

func GiveCard(Out, Hands []Card, isAm bool, ch <-chan int) ([]Card, []Card, bool, int) {
	input := gameConst.MaxValue
	if isAm {
		input = <-ch
		if input == gameConst.LeaveGameCode {
			return Out, []Card{}, false, gameConst.TakedCardState //
		}
	}
	if len(Out) == 0 {
		if isAm {
			for j, i := range Hands {
				if i.Val == input {
					Out, Hands = DecksUpdate(Out, Hands, j)
					return Out, Hands, true, Out[0].Val
				}
			}
		} else {
			Out, Hands = DecksUpdate(Out, Hands, 0)
			return Out, Hands, true, Out[0].Val
		}
	} else {
		for j, i := range Hands {
			if !isAm || isAm && i.Val == input {
				if Out[len(Out)-1].Val == 7 {
					if isSpecial(i.Val) || i.Val <= Out[len(Out)-1].Val {
						Out, Hands = DecksUpdate(Out, Hands, j)
						return Out, Hands, true, i.Val
					}
				} else {
					if isSpecial(i.Val) || i.Val >= Out[len(Out)-1].Val {
						Out, Hands = DecksUpdate(Out, Hands, j)
						return Out, Hands, true, i.Val
					}
				}
			}
		}
	}
	Hands = append(Hands, Out...)
	Out = Out[:0]
	return Out, Hands, false, gameConst.TakedCardState
}

func ReGiveCard(Out, Hands []Card, Value int, isIam bool, ch <-chan int) ([]Card, []Card, bool) {
	if isIam {
		var input int
		input = <-ch
		if input == gameConst.LeaveGameCode {
			return Out, []Card{}, true //
		}
		if input != Value {
			return Out, Hands, true
		}
	}
	for j, i := range Hands {
		if i.Val == Value {
			Out, Hands = DecksUpdate(Out, Hands, j)
			return Out, Hands, false
		}
	}
	return Out, Hands, true

}
