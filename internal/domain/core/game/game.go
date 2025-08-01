package game

import (
	"context"
	"sync"

	"github.com/beganov/gingonicserver/internal/domain/core/card"
	"github.com/beganov/gingonicserver/internal/domain/core/gameConst"
	"github.com/beganov/gingonicserver/internal/domain/core/placement"
	"github.com/rs/zerolog"
)

func (g *GameState) PreInitialization(maxPlayerCount int, Players map[int]int, end GameEndHandler, ctx context.Context) {
	// Задание начальных значений, не требующих пользовательского ввода
	g.Turn = 0
	g.PlayerNow = 0
	g.handler = end                                                                        //Инициализация ручки для окончания игры переданной комнатой
	g.logger = zerolog.Ctx(ctx)                                                            //Инициализация логгера из контекста
	g.IdMap, g.ch = ChannelsInit(Players)                                                  //Инициализация мап для игроков - каналов и рассадки
	g.Deck = card.NewDeck()                                                                //Инициализация колоды
	g.Hands, g.Closeds, g.Deck = card.HandInitialization(maxPlayerCount, g.Deck)           //Инициализация рук
	g.Hands, g.Openeds = card.OpenedsInitialization(maxPlayerCount, len(g.IdMap), g.Hands) //Инициализация открытых карт (не игроками)
}

func (g *GameState) Initialization(maxPlayerCount int) {
	// Задание начальных значений, после пользовательского ввода
	orderMap := g.PlayerInitialization(maxPlayerCount) //Подготовка игроков к игре
	// orderMap -  вспомогательная мапа - где по ключу в виде настоящей позиции выдается желаемая позиция
	g.Hands, g.Openeds = placement.ShufflePlayer(g.Hands, g.Openeds, orderMap)   //Расставляет игроков по случайным позициям
	g.Hands, g.Openeds, g.IdMap = placement.Orderer(g.Hands, g.Openeds, g.IdMap) //Передает первый ход игроку с тройкой
	g.Closeds = placement.LeaveCheck(g.Hands, g.Closeds)                         //Если кто-то вышел до начала игры - удаляем его закрытые карты
	g.ReverceIdMap = keyValueReverse(g.IdMap)                                    //Инициализация мапы для получения playerId по номеру за столом
}

func (g *GameState) PlayerInitialization(maxPlayerCount int) map[int]int {
	// Подготовка игроков к игре

	shuffleArr := placement.NewPlacementArray(maxPlayerCount, len(g.IdMap))
	// В зависимости от числа игроков создаём массив возможных позиций для игроков

	EndOfArray := shuffleArr[len(shuffleArr)-len(g.IdMap):]
	// Изначально игроков помещаем в конец очередности
	g.IdMap = ArraytoMap(EndOfArray, g.IdMap)
	//Заполняем мапу значениями массива

	var wg sync.WaitGroup
	for i := range g.IdMap {
		wg.Add(1)
		go func(i int) { // для кажого игрока асинхронно ждем инициализации его открытых карт
			defer wg.Done()
			g.OpenedsPlayerInitialization(i, g.IdMap[i])
		}(i)

	}
	wg.Wait()

	var orderMap map[int]int
	g.IdMap, orderMap = placement.TakeRandomPlacement(shuffleArr, g.IdMap) // Выдаем игрокам случайные позиции, на которых они будут сидеть
	return orderMap
}

func (g *GameState) OpenedsPlayerInitialization(playerId, k int) {
	// Инициализация открытых карт (игроками)
	z := 0
	var Openedshoosen int
	tempSlice := make([]card.Card, 0, gameConst.PackSize) //Создаем временный слайс для выбранных карт
	for z != gameConst.PackSize {                         //Пока не будет выбрано PackSize карт
		g.logger.Info().Int("PlayerId", playerId).Interface("Player hand", g.Hands[k]).Msg("Opened Initialization")
		card.SortCard(g.Hands[k])                     //Сортируем карты (чтобы не прыгали по экрану туда-сюда)
		Openedshoosen = <-g.ch[playerId]              // Достаем из канала переданное значение карты
		if Openedshoosen == gameConst.LeaveGameCode { // Если переданное значение - константа для выхода из комнаты - то выходим
			z = gameConst.PackSize
			break
		}
		for i := range g.Hands[k] {
			if g.Hands[k][i].Val == Openedshoosen {
				z++
				tempSlice, g.Hands[k] = card.DecksUpdate(tempSlice, g.Hands[k], i)
				// Если находим в руке игрока нужную карту, то перемещаем ее из руки в временный слайс
				break
			}
		}
	}
	g.Openeds[k] = tempSlice //Инициализируем открытые игрока временным слайсом
}

func (g *GameState) Game(maxPlayerCount int) {
	// Основной игровой цикл
	// Счетчик ходов
	turnCounter := 0
	g.Out = make([]card.Card, 0, gameConst.DeckSize)

	//флаг, показывающий была ли выложена карта
	istake := false

	// Числовой флаг, показывающий на каком этапе хода игрок
	// Если > 0, то показывает значение выложенной игроком карты
	// Если равен gameConst StartCardState - игрок не клал еще карту
	// Если равен gameConst TakedCardState - игрок берет еще карту
	var cardState int

	//флаг, показывающий закончил ли игрок ход
	var endTurnFlag bool
	// Счетчик для вышедших игроков
	var outCounter int
	for outCounter < maxPlayerCount-1 {
		turnCounter++
		outCounter = 0
		for i := 0; i < maxPlayerCount; i++ {
			g.IsMoved = false
			endTurnFlag = false
			cardState = gameConst.StartCardState
			for !endTurnFlag { // Пока не сработали условия для окончания хода
				if len(g.Hands[i]) == 0 && len(g.Closeds[i]) == 0 { //Если у игрока кончилась рука и закрытые
					outCounter++ //Увеличиваем число вышедших
					break        // Пропускаем его
				}
				card.SortCard(g.Hands[i])
				_, ok := g.ReverceIdMap[i] //Смотрим реален ли нынешний игрок или бот
				g.Turn = turnCounter
				g.PlayerNow = i
				outer(turnCounter, i, g.Out, g.Openeds, g.Closeds, g.Hands, g.logger) //логгируем состояние игры
				// Пытаемся выложить карту
				g.Hands[i], g.Out, cardState, endTurnFlag, istake, g.IsMoved = card.GiveCardLogic(g.Hands[i], g.Out, cardState, i, ok, endTurnFlag, istake, g.ch[g.ReverceIdMap[i]])
				// Пытаемся взять карту
				g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i] = card.TakeCard(g.Deck, g.Hands[i], g.Openeds[i], g.Closeds[i], istake)

			}
		}

	}
	g.logger.Info().Msg("GameEnd") // Логгируем конец игры
	g.Hands = [][]card.Card{}
	g.Openeds = [][]card.Card{}
	g.Closeds = [][]card.Card{}
	g.handler.OnGameEnd() //Зануляем все и сбрасываем room.isStart
}
