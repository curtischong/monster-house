package config

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

type Config struct {
	LimitValidatorConfig LimitValidatorConfig `yaml:"limitValidator"`
}

type LimitValidatorConfig struct {
	DailyLoadFrequencyLimit int   `yaml:"dailyLoadFrequencyLimit"`
	DailyLoadAmountLimit    int64 `yaml:"dailyLoadAmountLimit"`
	WeeklyLoadAmountLimit   int64 `yaml:"weeklyLoadAmountLimit"`
}

func NewConfig(
	filepath string,
) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
