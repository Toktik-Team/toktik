package model

// Comment 评论表 /*
type Comment struct {
	Model
	CommentId uint32 `json:"comment_id" column:"comment_id" gorm:"not null;index:comment_video"`
	VideoId   uint32 `json:"video_id" column:"video_id" gorm:"not null;index:comment_video"`
	UserId    uint32 `json:"user_id" column:"user_id" gorm:"not null"`
	Content   string `json:"content" column:"content"`
}
