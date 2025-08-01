package card

import "github.com/beganov/gingonicserver/internal/domain/core/gameConst"

// GiveCardLogic управляет основной логикой хода игрока, определяет кладет ли он карту
func GiveCardLogic(Hands, Out []Card, cardState, i int, iamindFlag, flag, istake bool, ch <-chan int) ([]Card, []Card, int, bool, bool, bool) {
	isMoved := false
	if len(Hands) != 0 { //Если рука не пустая
		if cardState > gameConst.TakedCardState { //Если игрок положил карту
			Out, Hands, flag = ReGiveCard(Out, Hands, cardState, iamindFlag, ch) //Игрок проверяет, может ли он положить еще карту
			istake = !flag                                                       // В зависимости от того положил игрок карту или нет, смотрим надо ли ему брать карту
			isMoved = true
		}
		if cardState == gameConst.StartCardState { //Если игрок только начал ход
			Out, Hands, istake, cardState = GiveCard(Out, Hands, iamindFlag, ch) //Игрок кладет карту
			flag = !istake                                                       // В зависимости от того положил игрок карту или нет, смотрим надо ли ему брать карту
			isMoved = true
		}
		if len(Out) > 0 { //Если на столе есть карты
			if Out[len(Out)-1].Val == 0 || Out[len(Out)-1].Val == 10 { //Если игрок положил 0 или 10
				Out = Out[:0] //Сброс всего на столе, повтор хода
				flag = false
				cardState = gameConst.StartCardState
				isMoved = false
			}
		}
		if len(Out) >= 4 && Out[len(Out)-1].Val == Out[len(Out)-2].Val && Out[len(Out)-2].Val == Out[len(Out)-3].Val && Out[len(Out)-3].Val == Out[len(Out)-4].Val {
			//Если на столе лежат 4 карты одного номинала подряд
			Out = Out[:0] //Сброс всего на столе, повтор хода
			flag = false
			cardState = gameConst.StartCardState
			isMoved = false
		}
	}
	return Hands, Out, cardState, flag, istake, isMoved
}

//Добор карты
func TakeCard(Deck, Hands, Openeds, Closeds []Card, istake bool) ([]Card, []Card, []Card, []Card) {
	if len(Deck) == 0 && len(Openeds) == 0 && len(Hands) == 0 && len(Closeds) != 0 {
		//Если открытых уже нет, а закрытые есть - берем из них
		Hands, Closeds = DecksUpdate(Hands, Closeds, 0)
	}
	if len(Deck) == 0 && len(Hands) == 0 {
		//Если открытые есть, а колоды нет - берем из них
		Hands = Openeds
		Openeds = Openeds[:0]
	}
	if len(Hands) < gameConst.PackSize && len(Deck) > 0 && istake {
		//Если есть колода, рука меньше необходимого и надо добирать - берем из колоды
		Hands, Deck = DecksUpdate(Hands, Deck, 0)
	}
	return Deck, Hands, Openeds, Closeds
}

//Игрок впервые кладет карту
func GiveCard(Out, Hands []Card, isAm bool, ch <-chan int) ([]Card, []Card, bool, int) {
	input := gameConst.MaxValue
	if isAm { // Если игрок не бот
		input = <-ch
		if input == gameConst.LeaveGameCode { // Если игрок вышел
			return Out, []Card{}, false, gameConst.TakedCardState //Зануляем карты, конец хода
		}
	}
	if len(Out) == 0 { // Если на столе пусто
		if isAm { // Если игрок - ищем нужную карту и кладем ее на стол
			for j, i := range Hands {
				if i.Val == input {
					Out, Hands = DecksUpdate(Out, Hands, j)
					return Out, Hands, true, Out[0].Val
				}
			}
		} else { // Если бот - просто скидываем худщую карту с рук на стол
			Out, Hands = DecksUpdate(Out, Hands, 0)
			return Out, Hands, true, Out[0].Val
		}
	} else { // Если на столе не пусто
		for j, i := range Hands {
			if !isAm || isAm && i.Val == input { // Если бот или нашлась карта на руках, что игрок хочет выложить
				if Out[len(Out)-1].Val == 7 { // Если 7 - кладем либо спец карту, либо карту меньше
					if isSpecial(i.Val) || i.Val <= Out[len(Out)-1].Val {
						Out, Hands = DecksUpdate(Out, Hands, j)
						return Out, Hands, true, i.Val
					}
				} else { // Если не 7 - кладем либо спец карту, либо карту больше
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
	return Out, Hands, false, gameConst.TakedCardState // Если ничего на руках не нашли, берем все со стола, конец хода
}

//Если надо доложить n-ю карту какого либо номинала
func ReGiveCard(Out, Hands []Card, Value int, isIam bool, ch <-chan int) ([]Card, []Card, bool) {
	if isIam { // Если это не бот
		input := <-ch
		if input == gameConst.LeaveGameCode { // Если игрок вышел
			return Out, []Card{}, true // конец хода с дисквилификацией
		}
		if input != Value { // Если разнятся значения последней выложенной и той, что пытаются выложить
			return Out, Hands, true // конец хода
		}
	}
	for j, i := range Hands {
		if i.Val == Value { // Если находим нужную в руке
			Out, Hands = DecksUpdate(Out, Hands, j) // переносим ее на стол
			return Out, Hands, false                //Даем возможность сходить еще
		}
	}
	return Out, Hands, true // Если не нашли - конец хода

}
