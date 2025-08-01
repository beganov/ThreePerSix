package player

type Player struct { //Структура игрока, по итогу нужна только в JSON
	Id   int `json:"id"`
	Move int `json:"move,omitempty"`
}
