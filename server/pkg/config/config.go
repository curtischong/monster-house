package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// This package is used to read the config and secrets file for the server

type AWSConfig struct {
	S3BucketName string `yaml:"s3BucketName"`
	S3Endpoint   string `yaml:"s3Endpoint"`
	Region       string `yaml:"region"`
}

type ServerConfig struct {
	Port                  int    `yaml:"port"`
	NgrokURL              string `yaml:"ngrokURL"`
	UseNBestGeneratedTags int    `yaml:"useNBestGeneratedTags"`
}
type DatabaseConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Config struct {
	ServerConfig   ServerConfig   `yaml:"server"`
	AWSConfig      AWSConfig      `yaml:"aws"`
	DatabaseConfig DatabaseConfig `yaml:"database"`
	SecretsConfig  SecretsConfig
}

type SecretsConfig struct {
	ImaggaUser   string `yaml:"imaggaUser"`
	ImaggaSecret string `yaml:"imaggaSecret"`
}

// NewConfig returns a config object for a given config and secrets path
func NewConfig(
	configPath, secretsPath string,
) (*Config, error) {
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}

	secretsFile, err := ioutil.ReadFile(secretsPath)
	var secretsConfig SecretsConfig
	err = yaml.Unmarshal(secretsFile, &secretsConfig)
	if err != nil {
		return nil, err
	}
	config.SecretsConfig = secretsConfig
	return &config, nil
}
