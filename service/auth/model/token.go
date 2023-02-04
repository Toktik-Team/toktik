package model

// UserToken 用户令牌儿表 /*
type UserToken struct {
	Token    string `gorm:"not null;primaryKey"`
	Username string `gorm:"not null;unique;size: 32"` // 用户名
	UserID   uint32 `gorm:"not null"`
	Role     string `gorm:"not null;default:0"` // 用户角色
}
