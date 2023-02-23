package model

// Favorite 点赞表 /*
type Favorite struct {
	Model
	UserId  uint32 `gorm:"not null;uniqueIndex:user_video"`
	VideoId uint32 `gorm:"not null;uniqueIndex:user_video;index:video"`
}
