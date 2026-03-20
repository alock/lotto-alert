package config

import (
	_ "embed"
	"encoding/json"
	"log"
)

var (
	//go:embed email.json
	emailConfigJson []byte
	EmailStruct     EmailConfig
)

type EmailConfig struct {
	From         string         `json:"from"`
	FromOverride string         `json:"fromOverride"`
	Participants map[int]string `json:"participants"`
}

func LoadConfigs() {
	err := json.Unmarshal(emailConfigJson, &EmailStruct)
	if err != nil {
		log.Fatal(err)
	}
}
