package dto

import "time"

type GameCreateRequest struct {
	PlayersAmount int `json:"players_amount"`
}

type GameDto struct {
	Id            string    `json:"id"`
	PlayersAmount int       `json:"players_amount"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

type ResultGameDto struct {
	Id            string    `json:"id"`
	PlayersAmount int       `json:"players_amount"`
	Position      int       `json:"position"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Gifts         []int     `json:"gifts_ids"`
	Status        string    `json:"status"`
	UserScore     int       `json:"score"`
	UserTime      int       `json:"time"`
	PrizeId       *int      `json:"prize_id"`
}
