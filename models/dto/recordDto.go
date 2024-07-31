package dto

type RecordDto struct {
	UserId        string `json:"user_id"`
	CacheUsername string `json:"username"`
	UserScore     int    `json:"user_score"`
	UserTime      int    `json:"user_time"`
}

type RecordCreateRequest struct {
	UserScore int `json:"user_score" bson:"user_score"`
	UserTime  int `json:"user_time" bson:"user_time"`
}
