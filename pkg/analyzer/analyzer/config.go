package analyzer

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Rules               RulesConfig `yaml:"rules"`
	SensitivePatterns   []string    `yaml:"sensitive_patterns"`
	AllowedSpecialChars string      `yaml:"allowed_special_chars"`
}

type RulesConfig struct {
	LowercaseStart bool `yaml:"lowercase_start"`
	EnglishOnly    bool `yaml:"english_only"`
	NoSpecialChars bool `yaml:"no_special_chars"`
	NoSensitive    bool `yaml:"no_sensitive"`
}

func DefaultConfig() Config {
	return Config{
		Rules: RulesConfig{
			LowercaseStart: true,
			EnglishOnly:    true,
			NoSpecialChars: true,
			NoSensitive:    true,
		},
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}
