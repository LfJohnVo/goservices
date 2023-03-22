package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func GetEnvValue(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return os.Getenv(key)
}

// funcion to get completely server route
func PWD() string {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	return GetProtocolVar() + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT")
}

func GetProtocolVar() string {
	isHTTPS, err := strconv.ParseBool(os.Getenv("SERVER_HTTPS"))

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	if isHTTPS {
		return "https://"
	} else {
		return "http://"
	}
}
