package models

import (
	"github.com/alexsuslov/ehttp"
	"go.mongodb.org/mongo-driver/v2/bson"

	"net/http"
	"snake-tournament/internal/config"
	"snake-tournament/models/dto"
	"sort"
	"time"

	"golang.org/x/exp/slices"
)

type Game struct {
	Id            bson.ObjectID `bson:"_id,omitempty"`
	PlayersAmount int           `bson:"players_amount"`
	Records       []Record      `bson:"records"`
	StartTime     time.Time     `bson:"start_time"`
	CloseTime     *time.Time    `bson:"close_time"`
}

func (g *Game) GetType() string {
	return "snake"
}

func (g *Game) GetPlayersAmount() int {
	return g.PlayersAmount
}

func NewGame(playersAmount int) *Game {
	currentTime := time.Now()
	gameStartTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), currentTime.Hour(), 0, 0, 0, time.UTC)

	game := Game{
		PlayersAmount: playersAmount,
		StartTime:     gameStartTime,
		Records:       []Record{},
	}

	return &game
}

func (g *Game) GetStartTime() time.Time {
	now := time.Now()

	return time.Date(g.StartTime.Year(), g.StartTime.Month(), g.StartTime.Day(), g.StartTime.Hour(), 0, 0, 0, now.Location())
}

func (g *Game) GetEndTime() time.Time {
	now := time.Now()

	if g.CloseTime != nil {
		return time.Date(g.CloseTime.Year(), g.CloseTime.Month(), g.CloseTime.Day(), g.CloseTime.Hour(), 0, 0, 0, now.Location())
	}

	endTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, 0, 0, 0, now.Location())

	if endTime.Before(g.GetStartTime().Add(time.Hour * 2)) {
		endTime = endTime.Add(time.Hour * 1)
	}

	return endTime
}

func (g *Game) GetId() string {
	return g.Id.Hex()
}

func (g *Game) SetId(id string) {
	g.Id, _ = bson.ObjectIDFromHex(id)
}

func (g *Game) AddPlayer(user dto.User) {
	currentTime := time.Now()
	startTime := currentTime.UTC()
	g.Records = append(g.Records, Record{
		UserId:        user.Id,
		CacheUsername: user.Username,
		UserScore:     nil,
		UserTime:      nil,
		StartTime:     &startTime,
		EndTime:       nil,
		Prize:         nil,
		UpdatesAmount: 0,
	})

	if len(g.Records) >= g.PlayersAmount && g.CloseTime == nil {
		endTime := currentTime.Add(time.Minute * 20)
		closeTime := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), endTime.Hour()+1, 0, 0, 0, time.UTC)

		if closeTime.Before(g.GetStartTime().Add(time.Hour * 2)) {
			startTime := g.GetStartTime()
			closeTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour()+2, 0, 0, 0, time.UTC)
		}

		g.CloseTime = &closeTime
	}
}

func (g *Game) FindPlayerResult(user dto.User) *Record {
	for i := range g.Records {
		if g.Records[i].UserId == user.Id {
			return &g.Records[i]
		}
	}

	return nil
}

func (g *Game) IsCloseToEnter(user dto.User) bool {
	if time.Now().After(g.GetEndTime()) {
		return true
	}

	for i := range g.Records {
		if g.Records[i].UserId == user.Id && g.Records[i].UpdatesAmount >= 3 {
			return true
		} else if g.Records[i].UserId == user.Id {
			return false
		}
	}

	return len(g.Records) >= g.PlayersAmount
}

func (g *Game) FullPlayers() bool {
	return len(g.Records) >= g.PlayersAmount
}
func (g *Game) FindPlayerPlace(user dto.User) int {
	return slices.IndexFunc(g.GetSortedRecords(), func(current Record) bool {
		return current.UserId == user.Id
	})
}

func (g *Game) GetSortedRecords() []Record {
	sort.Slice(g.Records, func(i, j int) bool {
		return g.compareRecords(g.Records[i], g.Records[j])
	})

	return g.Records
}

func (g *Game) compareRecords(r1, r2 Record) bool {
	if r1.UserScore == nil {
		return false
	}

	if r2.UserScore == nil {
		return true
	}

	if *r1.UserScore != *r2.UserScore {
		return *r1.UserScore > *r2.UserScore
	}

	if r1.UserTime == nil {
		return false
	}

	if r2.UserTime == nil {
		return true
	}

	return *r1.UserTime > *r2.UserTime
}

func (g *Game) GetAvailablePrizes(user dto.User) []int {
	place := g.FindPlayerPlace(user)

	if time.Now().Before(g.GetEndTime()) {
		return make([]int, 0)
	}

	if g.GetSortedRecords()[place].Prize != nil {
		return make([]int, 0)
	}

	return config.GamePrizeConfigConst.GetConfig(g.PlayersAmount, g.FindPlayerPlace(user))
}

func (g *Game) GetStatus(user dto.User) string {
	place := g.FindPlayerPlace(user)

	if time.Now().Before(g.GetEndTime()) {
		return "in_game"
	}

	if g.GetSortedRecords()[place].Prize != nil {
		return "prize_received"
	}

	if len(config.GamePrizeConfigConst.GetConfig(g.PlayersAmount, g.FindPlayerPlace(user))) != 0 {
		return "win"
	}

	return "lost"
}

func (g *Game) ToResultGameDto(user dto.User) dto.ResultGameDto {
	record := g.FindPlayerResult(user)

	var prizeId *int
	userScore := 0
	userTime := 0

	if record != nil && record.UserScore != nil {
		userScore = *record.UserScore
	}

	if record != nil && record.UserTime != nil {
		userTime = *record.UserTime
	}

	if record != nil && record.Prize != nil {
		prizeId = &record.Prize.SelectedPrize
	}

	return dto.ResultGameDto{
		Id:            g.GetId(),
		StartTime:     g.GetStartTime(),
		EndTime:       g.GetEndTime(),
		PlayersAmount: g.PlayersAmount,
		Position:      g.FindPlayerPlace(user) + 1,
		Gifts:         g.GetAvailablePrizes(user),
		Status:        g.GetStatus(user),
		UserScore:     userScore,
		UserTime:      userTime,
		PrizeId:       prizeId,
	}
}

func (g *Game) PasteResults(user dto.User, request dto.RecordCreateRequest) error {
	if time.Now().After(g.GetEndTime()) {
		return &ehttp.Error{
			Message: "This game is closed",
			Code:    http.StatusNotFound,
			Err:     nil,
		}
	}

	record := g.FindPlayerResult(user)

	if record == nil {
		return &ehttp.Error{
			Message: "Current user is not in game",
			Code:    http.StatusNotFound,
			Err:     nil,
		}
	}

	if record.UpdatesAmount >= 3 {
		return &ehttp.Error{
			Message: "Max updates amount exceeded",
			Code:    http.StatusNotFound,
			Err:     nil,
		}
	}

	newRecord := Record{
		UserScore: &request.UserScore,
		UserTime:  &request.UserTime,
	}

	if g.compareRecords(newRecord, *record) {
		record.UserScore = &request.UserScore
		record.UserTime = &request.UserTime
		currentTime := time.Now().UTC()
		record.EndTime = &currentTime
	}

	record.UpdatesAmount++

	return nil
}

func (g *Game) ToRecordsDto() []dto.RecordDto {
	result := make([]dto.RecordDto, 0, len(g.Records))

	records := g.GetSortedRecords()

	for i := 0; i < len(records); i++ {
		if records[i].IsFilled() {
			result = append(result, g.Records[i].ToRecordDto())
		}
	}

	return result
}

func (g *Game) ToGameDto() dto.GameDto {
	return dto.GameDto{
		Id:            g.GetId(),
		PlayersAmount: g.PlayersAmount,
		StartTime:     g.GetStartTime(),
		EndTime:       g.GetEndTime(),
	}
}

func (g *Game) SetPrizeForUser(user dto.User, prizeId int, email string) error {
	place := g.FindPlayerPlace(user)

	if place == -1 {
		return &ehttp.Error{
			Message: "User not in game",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	if g.FindPlayerResult(user).Prize != nil {
		return &ehttp.Error{
			Message: "Prize already selected",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	gifts := g.GetAvailablePrizes(user)

	giftIndex := slices.IndexFunc(gifts, func(current int) bool {
		return current == prizeId
	})

	if giftIndex == -1 {
		return &ehttp.Error{
			Message: "Prize not available",
			Code:    http.StatusForbidden,
			Err:     nil,
		}
	}

	prize := Prize{
		SelectedPrize: prizeId,
		PrizeEmail:    email,
		PrizeSended:   false,
	}

	result := g.FindPlayerResult(user)

	result.Prize = &prize

	return nil
}
