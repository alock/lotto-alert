package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// email-template.json holds placeholder values only (no PII) and is embedded so
// the package always builds, even on a clean checkout or in CI where the real
// config/email.json (gitignored) is absent.
//
//go:embed email-template.json
var templateConfigJSON []byte

var EmailStruct EmailConfig

type EmailConfig struct {
	From         string         `json:"from"`
	FromOverride string         `json:"fromOverride"`
	Participants map[int]string `json:"participants"`
}

// configCandidates returns the paths searched for the real participant config,
// in priority order:
//  1. $LOTTO_EMAIL_CONFIG (explicit override)
//  2. ./config/email.json (running from the repo root, e.g. `go run .`)
//  3. email.json next to the executable (the Raspberry Pi deploy layout)
func configCandidates() []string {
	var paths []string
	if p := os.Getenv("LOTTO_EMAIL_CONFIG"); p != "" {
		paths = append(paths, p)
	}
	paths = append(paths, filepath.Join("config", "email.json"))
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "email.json"))
	}
	return paths
}

// LoadConfigs loads the participant config from the first candidate path that
// exists, falling back to the embedded placeholder template. The template has no
// real participants, so commands still run (they simply never find a winner) —
// which is exactly what we want for tests and dry runs.
func LoadConfigs() error {
	data := templateConfigJSON
	for _, p := range configCandidates() {
		b, err := os.ReadFile(p)
		if err == nil {
			data = b
			break
		}
		if !os.IsNotExist(err) {
			return fmt.Errorf("reading %s: %w", p, err)
		}
	}
	if err := json.Unmarshal(data, &EmailStruct); err != nil {
		return fmt.Errorf("parsing participant config: %w", err)
	}
	return nil
}
