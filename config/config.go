package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Config The config object from the config file
type Config struct {
	Spotify struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	} `yaml:"Spotify"`
	SessionSecret string `yaml:"SessionSecret"`
	DomainName    string `yaml:"DomainName"`
	DB            struct {
		Host     string `yaml:"Host"`
		Port     int    `yaml:"Port"`
		Name     string `yaml:"Name"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	} `yaml:"DB"`
}

// ReadConfig Function to read the config from the yaml config file
// It need to be in the config dir
func ReadConfig() (Config, error) {
	var config Config

	dat, errRead := ioutil.ReadFile("config/config.yaml")
	if errRead != nil {
		return config, errRead
	}

	errMarsh := yaml.UnmarshalStrict(dat, &config)
	if errMarsh != nil {
		return config, errMarsh
	}

	return config, nil
}
