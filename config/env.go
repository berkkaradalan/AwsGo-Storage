package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Env struct{
	JWT_SECRET_KEY		string `mapstructure:"JWT_SECRET_KEY"`
	JWT_EXPIRE_HOURS	int `mapstructure:"JWT_EXPIRE_HOURS"`
	S3_BUCKET_NAME		string `mapstructure:"S3_BUCKET_NAME"`
}

func LoadEnv() (*Env){
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Env file not loaded. Here's what happened : %v ", err)
		panic(err)
	}

	jwtExpireHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRE_HOURS"))

	if err != nil {
		log.Printf("Env file not loaded. Here's what happened : %v ", err)
		panic(err)
	}

	return &Env{
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
		JWT_EXPIRE_HOURS: jwtExpireHours,
		S3_BUCKET_NAME: os.Getenv("S3_BUCKET_NAME"),
	}
}