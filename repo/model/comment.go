package model

// Comment 评论表 /*
type Comment struct {
	Model
	CommentId uint32 `json:"comment_id" column:"comment_id"`
	VideoId   uint32 `json:"video_id" column:"video_id"`
	UserId    uint32 `json:"user_id" column:"user_id"`
	Content   string `json:"content" column:"content"`
}
