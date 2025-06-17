package linter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents the linter configuration
type Config struct {
	DisabledRules []string       `json:"disabled_rules"`
	RuleSettings  map[string]any `json:"rule_settings"`
}

// IsRuleEnabled returns whether a rule is enabled
func (c Config) IsRuleEnabled(ruleName string) bool {
	for _, disabled := range c.DisabledRules {
		if disabled == ruleName {
			return false
		}
	}
	return true
}

// GetRuleSetting returns a setting for a rule
func (c Config) GetRuleSetting(ruleName, settingName string, defaultValue any) any {
	if settings, ok := c.RuleSettings[ruleName]; ok {
		if settingsMap, ok := settings.(map[string]any); ok {
			if value, ok := settingsMap[settingName]; ok {
				return value
			}
		}
	}
	return defaultValue
}

// LoadConfig loads a configuration from a file
func LoadConfig(path string) (Config, error) {
	var config Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	return config, nil
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		DisabledRules: []string{},
		RuleSettings: map[string]any{
			"max-line-length":           100,
			"max-file-lines":            1000,
			"max-public-methods":        20,
			"max-returns":               6,
			"max-branches":              12,
			"max-statements":            50,
			"max-attributes":            10,
			"max-local-variables":       15,
			"function-arguments-number": 10,
		},
	}
}

// FindConfigFile searches for a config file in the current directory and parent directories
func FindConfigFile() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		configPath := filepath.Join(dir, "gdlintrc.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		configPath = filepath.Join(dir, ".gdlintrc.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		configPath = filepath.Join(dir, "gdlintrc")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		configPath = filepath.Join(dir, ".gdlintrc")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return ""
}
