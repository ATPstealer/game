package cfg

import (
	"log"
	"os"
	"strconv"
)

type Vars struct {
	DBUser        string `env:"DB_USER" default:"root"`
	DBPassword    string `env:"DB_PASSWORD" default:"pass"`
	DBHost        string `env:"DB_HOST" default:"localhost"`
	DBPort        string `env:"DB_PORT" default:"3306"`
	DBDataBase    string `env:"DB_DATABASE" default:"db"`
	AppPort       string `env:"APP_PORT" default:"8000"`
	Secure        bool   `env:"SECURE" default:"false"`
	GoogleAPI     string `env:"GOOGLE_API"`
	GoogleSheetID string `env:"GOOGLE_SHEET_ID"`
	Init          bool   `env:"INIT" default:"true"` // change to separate run
}

var (
	Config = Vars{}
)

func LoadConfig() {
	if user := os.Getenv("DB_USER"); user != "" {
		Config.DBUser = user
	}
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		Config.DBPassword = password
	}
	if host := os.Getenv("DB_HOST"); host != "" {
		Config.DBHost = host
	}
	if port := os.Getenv("DB_PORT"); port != "" {
		Config.DBPort = port
	} else {
		Config.DBPort = "3306"
	}
	if db := os.Getenv("DB_DATABASE"); db != "" {
		Config.DBDataBase = db
	}
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

}
