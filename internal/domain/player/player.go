package player

type Player struct {
	Id   int `json:"id"`
	Move int `json:"move,omitempty"`
}
