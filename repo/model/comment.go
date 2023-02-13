package model

// Comment 评论表 /*
type Comment struct {
	Model
	VideoId uint32 `json:"video_id" column:"video_id"`
	UserId  uint32 `json:"user_id" column:"user_id"`
	Content string `json:"content" column:"content"`
}
