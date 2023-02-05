package main

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"toktik/config"
	auth "toktik/service/auth/model"
	publish "toktik/service/publish/model"
)

func main() {
	var err error
	db, err := gorm.Open(
		postgres.New(
			postgres.Config{
				DSN: config.DSN,
			}), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		})
	if err != nil {
		panic(fmt.Errorf("db connection failed: %v", err))
	}
	err = db.AutoMigrate(&auth.UserToken{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
	err = db.AutoMigrate(&publish.Publish{})
	if err != nil {
		panic(fmt.Errorf("db migrate failed: %v", err))
	}
}
