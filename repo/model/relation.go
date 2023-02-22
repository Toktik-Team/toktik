package model

// Relation 关系表 /*
type Relation struct {
	Model           // 基础模型
	UserId   uint32 `gorm:"not null;index:relation_user_target,unique"` // 用户ID
	TargetId uint32 `gorm:"not null;index:relation_user_target,unique"` // 目标ID
}
