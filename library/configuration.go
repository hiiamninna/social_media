package library

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PrometheusAddr string
	JWTSecret      string
	BcryptSalt     int
	DB             Database
	S3Config       S3Config
}

// Setting up the environment to be used
func NewConfiguration() (Config, error) {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("load env : %w", &err)
	}

	config := Config{
		PrometheusAddr: EnvString("PROMETHEUS_ADDRESS"),
		JWTSecret:      EnvString("JWT_SECRET"),
		BcryptSalt:     EnvInt("BCRYPT_SALT"),
		DB: Database{
			Env:      EnvString("ENV"),
			Name:     EnvString("DB_NAME"),
			Host:     EnvString("DB_HOST"),
			Port:     EnvString("DB_PORT"),
			Username: EnvString("DB_USERNAME"),
			Password: EnvString("DB_PASSWORD"),
		},
		S3Config: S3Config{
			ID:         EnvString("S3_ID"),
			SecretKey:  EnvString("S3_SECRET_KEY"),
			BucketName: EnvString("S3_BUCKET_NAME"),
		},
	}

	return config, nil
}

func EnvString(key string) string {
	return os.Getenv(key)
}

func EnvInt(key string) int {
	val, _ := strconv.Atoi(os.Getenv(key))
	return val
}
