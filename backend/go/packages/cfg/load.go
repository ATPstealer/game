package cfg

import (
	"log"
	"os"
	"strconv"
)

type Vars struct {
	AppPort       string `env:"APP_PORT" default:"8000"`
	Secure        bool   `env:"SECURE" default:"false"`
	GoogleAPI     string `env:"GOOGLE_API"`
	GoogleSheetID string `env:"GOOGLE_SHEET_ID"`
	Init          bool   `env:"INIT" default:"true"`
	MongoUser     string `env:"MONGO_USER" default:"root"`
	MongoPassword string `env:"MONGO_PASSWORD"`
	MongoHost     string `env:"MONGO_HOST" default:"localhost"`
	MongoDatabase string `env:"MONGO_DATABASE"`
	MongoPort     string `env:"MONGO_PORT"  default:"27017"`
}

var (
	Config = Vars{}
)

func LoadConfig() {
	if appPort := os.Getenv("APP_PORT"); appPort != "" {
		appPortInt, err := strconv.Atoi(appPort)
		if err != nil {
			log.Fatal("Port isn't Int")
		}
		Config.AppPort = strconv.Itoa(appPortInt)
	} else {
		Config.AppPort = "8000"
	}
	if googleApi := os.Getenv("GOOGLE_API"); googleApi != "" {
		Config.GoogleAPI = googleApi
	}
	if googleSheetID := os.Getenv("GOOGLE_SHEET_ID"); googleSheetID != "" {
		Config.GoogleSheetID = googleSheetID
	}
	if secureStr := os.Getenv("SECURE"); secureStr != "" {
		secure, err := strconv.ParseBool(secureStr)
		if err != nil {
			log.Fatal("secure Env isn't bool")
		}
		Config.Secure = secure
	} else {
		Config.Secure = false
	}
	if secureStr := os.Getenv("INIT"); secureStr != "" {
		init, err := strconv.ParseBool(secureStr)
		if err != nil {
			log.Fatal("Init Env isn't bool")
		}
		Config.Init = init
	} else {
		Config.Init = true
	}
	if mongoUser := os.Getenv("MONGO_USER"); mongoUser != "" {
		Config.MongoUser = mongoUser
	}
	if mongoPassword := os.Getenv("MONGO_PASSWORD"); mongoPassword != "" {
		Config.MongoPassword = mongoPassword
	}
	if mongoHost := os.Getenv("MONGO_HOST"); mongoHost != "" {
		Config.MongoHost = mongoHost
	}
	if mongoDatabase := os.Getenv("MONGO_DATABASE"); mongoDatabase != "" {
		Config.MongoDatabase = mongoDatabase
	}
	if mongoPort := os.Getenv("MONGO_PORT"); mongoPort != "" {
		Config.MongoPort = mongoPort
	} else {
		Config.MongoPort = "27017"
	}
}
