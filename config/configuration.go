package config

import (
	"fmt"
	"os"
	"social_media/library"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	DB       library.Database
	S3Config library.S3Config
}

type AppConfig struct {
	Env            string
	PrometheusAddr string
	JWTSecret      string
	BcryptSalt     int
}

// Setting up the environment to be used
func NewConfiguration() (Config, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println(time.Now().Format("2006-01-02 15:01:02 "), "load env : "+err.Error())
	}

	config := Config{
		App: AppConfig{
			Env:            EnvString("ENV"),
			PrometheusAddr: EnvString("PROMETHEUS_ADDRESS"),
			JWTSecret:      EnvString("JWT_SECRET"),
			BcryptSalt:     EnvInt("BCRYPT_SALT"),
		},
		DB: library.Database{
			Name:         EnvString("DB_NAME"),
			Host:         EnvString("DB_HOST"),
			Port:         EnvString("DB_PORT"),
			Username:     EnvString("DB_USERNAME"),
			Password:     EnvString("DB_PASSWORD"),
			Params:       EnvString("DB_PARAMS"),
			MaxIdleTime:  EnvDurationTime("DB_MAX_IDLE_TIME"),
			MaxLifeTime:  EnvDurationTime("DB_MAX_LIFE_TIME"),
			MaxIdleConns: EnvInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: EnvInt("DB_MAX_OPEN_CONNS"),
		},
		S3Config: library.S3Config{
			ID:         EnvString("S3_ID"),
			SecretKey:  EnvString("S3_SECRET_KEY"),
			BucketName: EnvString("S3_BUCKET_NAME"),
			Region:     EnvString("S3_REGION"),
		},
	}

	return config, nil
}

// TODO : should we add the default value ?
// - yes, why?
// - no, why?
func EnvString(key string) string {
	return os.Getenv(key)
}

func EnvInt(key string) int {
	val, _ := strconv.Atoi(os.Getenv(key))
	return val
}

func EnvDurationTime(key string) time.Duration {
	val, _ := strconv.Atoi(os.Getenv(key))
	return time.Duration(val)
}
