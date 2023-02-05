package model

import "gorm.io/gorm"

// Publish 视频表 /*
type Publish struct {
	gorm.Model
	UserId   int64  `json:"user_id" column:"user_id"`
	Title    string `json:"title" column:"title"`
	PlayUrl  string `json:"play_url" column:"play_url"`
	CoverUrl string `json:"cover_url" column:"cover_url"`
}
