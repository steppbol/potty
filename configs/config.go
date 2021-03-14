package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const (
	config = "./configs/bootstrap.yml"
)

type Config struct {
	Application
	Database
}

type Application struct {
	XLSXExportPath string `yaml:"xlsx-export-path"`
}

type applicationConfig struct {
	Application yaml.Node
}

type Database struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	User           string `yaml:"user"`
	Password       string `yaml:"password"`
	Name           string `yaml:"name"`
	Schema         string `yaml:"schema"`
	GenerateSchema bool   `yaml:"generate-schema"`
	SSLMode        string `yaml:"ssl-mode"`
}

type databaseConfig struct {
	Database yaml.Node
}

func Setup() (*Config, error) {
	buff, err := ioutil.ReadFile(config)
	if err != nil {
		return nil, err
	}

	var dc databaseConfig
	err = yaml.Unmarshal(buff, &dc)
	if err != nil {
		return nil, err
	}

	d := Database{}
	err = dc.Database.Decode(&d)

	var ac applicationConfig
	err = yaml.Unmarshal(buff, &ac)

	if err != nil {
		return nil, err
	}

	a := Application{}
	err = ac.Application.Decode(&a)

	return &Config{
		Application: a,
		Database:    d,
	}, nil
}
