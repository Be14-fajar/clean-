package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

var (
	JWT_KEY string = ""
)

type AppConfig struct {
	DBUser string
	DBPass string
	DBHost string
	DBPort int
	DBName string
	jwtKey string
}

func InitConfig() *AppConfig {
	return ReadEnv()
}

func ReadEnv() *AppConfig {
	app := AppConfig{}
	isRead := true

	if val, found := os.LookupEnv("JWT_KEY"); found {
		app.jwtKey = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBUser"); found {
		app.DBUser = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPass"); found {
		app.DBPass = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBHost"); found {
		app.DBHost = val
		isRead = false
	}
	if val, found := os.LookupEnv("DBPort"); found {
		cnv, _ := strconv.Atoi(val)
		app.DBPort = cnv
		isRead = false
	}
	if val, found := os.LookupEnv("DBName"); found {
		app.DBName = val
		isRead = false
	}

	if isRead {
		viper.AddConfigPath(".")
		viper.SetConfigName("local")
		viper.SetConfigType("env")

		err := viper.ReadInConfig()
		if err != nil {
			log.Println("error read config : ", err.Error())
			return nil
		}
		err = viper.Unmarshal(&app)
		if err != nil {
			log.Println("error parse config : ", err.Error())
			return nil
		}
	}

	JWT_KEY = app.jwtKey
	return &app
}
