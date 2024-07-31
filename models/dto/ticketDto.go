package dto

type CloseTicketRequestDto struct {
	UserId       string `json:"user_id"`
	GameType     string `json:"game_type"`
	PlayerAmount int    `json:"player_amount"`
	GameId       string `json:"game_id"`
}
