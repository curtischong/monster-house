package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type AWSConfig struct {
	S3BucketName string `yaml:"s3BucketName"`
	S3Endpoint	string `yaml:"s3Endpoint"`
	Region string `yaml:"region"`
}

type ServerConfig struct {
	Port int `yaml:"port"`
}

type Config struct {
	ServerConfig ServerConfig `yaml:"server"`
	AWSConfig AWSConfig `yaml:"aws"`
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
