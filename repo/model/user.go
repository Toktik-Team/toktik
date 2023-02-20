package model

// User 用户表 /*
type User struct {
	Model                   // 基础模型
	Username        string  `gorm:"not null;unique;size: 32"` // 用户名
	Password        *string `gorm:"not null;size: 32"`        // 密码
	FollowCount     uint32  `gorm:"default:0"`                // 关注总数
	FollowerCount   uint32  `gorm:"default:0"`                // 粉丝总数
	Avatar          *string // 用户头像
	BackgroundImage *string // 背景图片
	Signature       *string // 个人简介
	TotalFavorited  *uint32 `gorm:"default:0"`       // 获赞数量
	WorkCount       *uint32 `gorm:"default:0"`       // 作品数量
	FavoriteCount   *uint32 `gorm:"default:0"`       // 点赞数量
	Name            string  `gorm:"not null;unique"` // 用户名称
	Role            string  `gorm:"default:0"`       // 用户角色
}
