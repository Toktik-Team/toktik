package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strings"
)

var EnvConfig = envConfigSchema{}

const WebServiceName = "toktik-api-gateway"
const WebServiceAddr = ":40126"

const AuthServiceName = "toktik-auth-api"
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
	ENV string `env:"ENV,DREAM_ENV"`

	CONSUL_ADDR string `env:"CONSUL_ADDR,DREAM_SERVICE_DISCOVERY_URI"`

	PGSQL_HOST     string
	PGSQL_PORT     string
	PGSQL_USER     string
	PGSQL_PASSWORD string
	PGSQL_DBNAME   string

	REDIS_ADDR     string
	REDIS_PASSWORD string
	REDIS_DB       string
}

func (s envConfigSchema) IsDev() bool {
	return s.ENV == "dev" || s.ENV == "TESTING"
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
		envNameAlt := make([]string, 0)
		fieldName := typeOfS.Field(i).Name
		fieldValue := v.Field(i).Interface()

		envNameAlt = append(envNameAlt, fieldName)
		if fieldTag, ok := typeOfS.Field(i).Tag.Lookup("env"); ok && len(fieldTag) > 0 {
			tags := strings.Split(fieldTag, ",")
			envNameAlt = append(envNameAlt, tags...)
		}

		configDefaultValue := fieldValue.(string)
		envValue := resolveEnv(envNameAlt)

		if EnvConfig.IsDev() {
			fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
		}
		if len(envValue) > 0 {
			reflect.ValueOf(&EnvConfig).Elem().Field(i).SetString(envValue)
		}
	}
}

func resolveEnv(configKeys []string) string {
	for _, item := range configKeys {
		envValue := os.Getenv(item)
		if envValue != "" {
			return envValue
		}
	}
	return ""
}
