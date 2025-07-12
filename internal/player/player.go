package player

type Player struct {
	Id int `json:"id"`
}

type PlayerMove struct {
	Id   int `json:"id"`
	Move int `json:"move"`
}

func NewPlayer(id int) *Player {
	p := &Player{}
	p.Id = id
	return p
}
