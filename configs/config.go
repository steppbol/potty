package configs

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

const (
	config = "resources/bootstrap.yml"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Scheme   string `yaml:"scheme"`
}

type databaseConfig struct {
	Database yaml.Node
}

func GetDatabaseConfig() (*Database, error) {
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

	return &d, nil
}
