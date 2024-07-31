package models

import (
	"snake-tournament/models/dto"
	"time"
)

type Record struct {
	UserId        string     `bson:"user_id"`
	CacheUsername string     `bson:"username"`
	UserScore     *int       `bson:"user_score,omitempty"`
	UserTime      *int       `bson:"user_time,omitempty"`
	StartTime     *time.Time `bson:"start_time,omitempty"`
	EndTime       *time.Time `bson:"end_time,omitempty"`
	Prize         *Prize     `bson:"prize"`
	UpdatesAmount int        `bson:"updates_amount"`
}

func (r *Record) ToRecordDto() dto.RecordDto {
	userScore := 0
	userTime := 0

	if r.UserScore != nil && r.UserTime != nil {
		userScore = *r.UserScore
		userTime = *r.UserTime
	}

	return dto.RecordDto{
		UserId:        r.UserId,
		CacheUsername: r.CacheUsername,
		UserScore:     userScore,
		UserTime:      userTime,
	}
}

func (r *Record) IsFilled() bool {
	return r.UserScore != nil && r.UserTime != nil
}
