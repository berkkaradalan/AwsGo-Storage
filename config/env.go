package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct{
	JWT_SECRET_KEY		string `mapstructure:"JWT_SECRET_KEY"`
}

func LoadEnv() (*Env){
	err := godotenv.Load(".env")

	if err != nil {
		log.Printf("Env file not loaded. Here's what happened : %v ", err)
		panic(err)
	}

	return &Env{
		JWT_SECRET_KEY: os.Getenv("JWT_SECRET_KEY"),
	}
}