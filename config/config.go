package config

import (
	_ "embed"
	"encoding/json"
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

func LoadConfigs() error {
	return json.Unmarshal(emailConfigJson, &EmailStruct)
}
