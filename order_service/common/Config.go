package common

import (
	"encoding/json"
	"log"
	"os"
)

type appConfig struct {
	PortSite       uint16
	SecretKey     string
	DatabaseConfig DbConfig
}

var Configs = appConfig{}

func (config *appConfig) Port() uint16 {
	return config.PortSite
}

func LoadConfiguration(filePath string, isProduction bool) {
	file, _ := os.Open(filePath)
	defer file.Close()

	decoder := json.NewDecoder(file)
	cf := map[string]appConfig{}

	err := decoder.Decode(&cf)
	if err != nil {
		log.Println("Load config error:", err)
	}

	if isProduction {
		Configs = cf["production"]
	} else {
		Configs = cf["development"]
	}
	MAIN_DB_CONSTRING = LoadTemplate(TEMPLATE_DB_CONSTRING, Configs.DatabaseConfig)
}