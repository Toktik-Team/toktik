package config

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
)

var EnvConfig = envConfigSchema{}

func (s *envConfigSchema) GetDSN() string {
	return dsn
}

var dsn string

func init() {
	envInit()
	envValidate()
	dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		EnvConfig.PGSQL_HOST,
		EnvConfig.PGSQL_USER,
		EnvConfig.PGSQL_PASSWORD,
		EnvConfig.PGSQL_DBNAME,
		EnvConfig.PGSQL_PORT)
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

	S3_ENDPOINT_URL: "http://localhost:9000",
	S3_PUBLIC_URL:   "http://localhost:9000",
	S3_BUCKET:       "bucket",
	S3_SECRET_ID:    "minio",
	S3_SECRET_KEY:   "12345678",
	S3_PATH_STYLE:   "true",
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

	S3_ENDPOINT_URL string
	S3_PUBLIC_URL   string
	S3_BUCKET       string
	S3_SECRET_ID    string
	S3_SECRET_KEY   string
	S3_PATH_STYLE   string
}

func (s *envConfigSchema) IsDev() bool {
	return s.ENV == "dev" || s.ENV == "TESTING"
}

func envValidate() {
	EnvConfig.CONSUL_ADDR = strings.TrimPrefix(EnvConfig.CONSUL_ADDR, "consul://")
}

// envInit Reads .env as environment variables and fill corresponding fields into EnvConfig.
// To use a value from EnvConfig , simply call EnvConfig.FIELD like EnvConfig.ENV
// Note: Please keep Env as the first field of envConfigSchema for better logging.
func envInit() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file, ignored")
	}
	v := reflect.ValueOf(defaultConfig)
	typeOfV := v.Type()

	for i := 0; i < v.NumField(); i++ {
		envNameAlt := make([]string, 0)
		fieldName := typeOfV.Field(i).Name
		fieldValue := v.Field(i).Interface()

		envNameAlt = append(envNameAlt, fieldName)
		if fieldTag, ok := typeOfV.Field(i).Tag.Lookup("env"); ok && len(fieldTag) > 0 {
			tags := strings.Split(fieldTag, ",")
			envNameAlt = append(envNameAlt, tags...)
		}

		configDefaultValue := fieldValue.(string)
		envValue := resolveEnv(envNameAlt, configDefaultValue)

		if EnvConfig.IsDev() {
			fmt.Printf("Reading field[ %s ] default: %v env: %s\n", fieldName, configDefaultValue, envValue)
		}
		if len(envValue) > 0 {
			reflect.ValueOf(&EnvConfig).Elem().Field(i).SetString(envValue)
		}
	}
}

func resolveEnv(configKeys []string, defaultValue string) string {
	for _, item := range configKeys {
		envValue := os.Getenv(item)
		if envValue != "" {
			return envValue
		}
	}
	return defaultValue
}
