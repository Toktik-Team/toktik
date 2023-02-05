package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
)

var EnvConfig = envConfigSchema{}

const WebServiceName = "toktik-api-gateway"
const WebServiceAddr = ":40126"

const AuthServiceName = "toktik-auth"
const AuthServiceAddr = ":40127"

var DSN string

func init() {
	envInit()
	DSN = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		EnvConfig.PGSQL_HOST,
		EnvConfig.PGSQL_USER,
		EnvConfig.PGSQL_PASSWORD,
		EnvConfig.PGSQL_DBNAME,
		EnvConfig.PGSQL_PORT)
}

type envConfigSchema struct {
	ENV string

	CONSUL_ADDR string

	PGSQL_HOST     string
	PGSQL_PORT     string
	PGSQL_USER     string
	PGSQL_PASSWORD string
	PGSQL_DBNAME   string

	REDIS_ADDR     string
	REDIS_PASSWORD string
	REDIS_DB       string
}

var defaultConfig = envConfigSchema{
	ENV: "dev",

	CONSUL_ADDR: "127.0.0.1:8500",

	PGSQL_HOST:     "localhost",
	PGSQL_PORT:     "5432",
	PGSQL_USER:     "postgres",
	PGSQL_PASSWORD: "",
	PGSQL_DBNAME:   "postgres",

	REDIS_ADDR:     "localhost:6379",
	REDIS_PASSWORD: "",
	REDIS_DB:       "0",
}

// envInit Reads .env as environment variables and fill corresponding fields into EnvConfig.
// To use a value from EnvConfig , simply call EnvConfig.FIELD like EnvConfig.ENV
func envInit() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	v := reflect.ValueOf(defaultConfig)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldName := typeOfS.Field(i).Name
		fieldValue := v.Field(i).Interface()

		configKey := fieldName
		var configValue string
		configDefaultValue := fieldValue.(string)
		envValue := os.Getenv(configKey)
		if envValue != "" {
			configValue = envValue
		} else {
			configValue = configDefaultValue
		}
		if EnvConfig.ENV == "dev" {
			fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
		}
		reflect.ValueOf(&EnvConfig).Elem().Field(i).SetString(configValue)
	}
}
