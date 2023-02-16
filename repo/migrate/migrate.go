package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"toktik/constant/config"
	"toktik/repo/model"
	auth "toktik/service/auth/model"
)

func main() {
	var err error
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.EnvConfig.GetDSN(),
			}), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	err = db.AutoMigrate(&auth.UserToken{}, &model.User{}, &model.Video{}, &model.Comment{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}
