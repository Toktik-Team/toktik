package main

import (
	"fmt"
	"toktik/constant/config"
	"toktik/repo/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var err error
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.EnvConfig.GetDSN(),
			}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	err = db.AutoMigrate(&model.UserToken{}, &model.User{}, &model.Video{}, &model.Comment{}, &model.Relation{}, &model.Favorite{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}
