package model

// Video 视频表 /*
type Video struct {
	Model
	UserId    uint32 `json:"user_id"`
	Title     string `json:"title"`
	FileName  string `json:"play_name"`
	CoverName string `json:"cover_name"`
}
