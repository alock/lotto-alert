package config

import (
	_ "embed"
	"encoding/json"
	"log"
	"time"

	"github.com/alock/lotto-alert/util"
)

var (
	//go:embed email.json
	emailConfigJson []byte
	EmailStruct     emailConfig
	//go:embed special-dates.json
	specialDatesConfigJson []byte
	SpecialDates           map[string]PrizeInfo
)

type emailConfig struct {
	From         string         `json:"from"`
	FromOverride string         `json:"fromOverride"`
	Participants map[int]string `json:"participants"`
}

type PrizeInfo struct {
	Amount int    `json:"amount"`
	Reason string `json:"reason"`
}

func LoadConfigs() {
	err := json.Unmarshal(emailConfigJson, &EmailStruct)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(specialDatesConfigJson, &SpecialDates)
	if err != nil {
		log.Fatal(err)
	}
}

func GetDatesPrizeInfo(t time.Time) PrizeInfo {
	specialDateInfo, ok := SpecialDates[util.GetStringOfDate(t)]
	if ok {
		return specialDateInfo
	}
	// every saturday not special is $50
	if t.Weekday() == 6 {
		return PrizeInfo{
			Amount: 50,
			Reason: "a Saturday",
		}
	}
	return PrizeInfo{
		Amount: 30,
		Reason: "",
	}
}
