package model

// Favorite 点赞表 /*
type Favorite struct {
	Model
	UserId  uint32 `gorm:"not null;index:user_video;unique"`
	VideoId uint32 `gorm:"not null;index:user_video;index:video;unique"`
}
