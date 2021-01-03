package tagger

import (
	"../config"
)

type ImaggaClient struct{
	config *config.Config
}

func NewImaggaClient(
	config *config.Config,
)*ImaggaClient{
	return 	&ImaggaClient{
		config	:config,
	}
}