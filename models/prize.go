package models

type Prize struct {
	SelectedPrize int    `bson:"selected_prize"`
	PrizeEmail    string `bson:"prize_email"`
	PrizeSended   bool   `bson:"prize_sended"`
}
