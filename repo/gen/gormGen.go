package main

import (
	"toktik/constant/config"
	"toktik/repo/model"
	auth "toktik/service/auth/model"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// Querier Dynamic SQL
type Querier interface {
	// FilterWithNameAndRole SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "repo",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	gormdb, _ := gorm.Open(postgres.Open(config.EnvConfig.GetDSN()))
	g.UseDB(gormdb) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(auth.UserToken{}, model.Video{}, model.User{}, model.Comment{}, model.Relation{}, model.Favorite{})

	// Generate Type Safe API with Dynamic SQL defined on Querier interface
	g.ApplyInterface(func(Querier) {}, auth.UserToken{}, model.Video{}, model.User{}, model.Comment{})

	// Generate the code
	g.Execute()
}
