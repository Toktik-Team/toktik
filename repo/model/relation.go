package model

// Relation 关系表 /*
type Relation struct {
	Model           // 基础模型
	UserId   uint32 `gorm:"not null;uniqueIndex:relation_user_target"` // 用户ID
	TargetId uint32 `gorm:"not null;uniqueIndex:relation_user_target"` // 目标ID
}
