package configs

import (
	"encoding/json"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/google/uuid"
)

const (
	bootstrap = "./configs/bootstrap.yml"
	security  = "./configs/security.json"
)

type Config struct {
	Application
	Database
	Server
	Cache
	Security
}

type Application struct {
	XLSXExportPath string `yaml:"xlsx-export-path"`
}

type applicationConfig struct {
	Application yaml.Node
}

type Server struct {
	Mode string `yaml:"mode"`
	Port int    `yaml:"port"`
}

type serverConfig struct {
	Server yaml.Node
}

type Cache struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type cacheConfig struct {
	Cache yaml.Node
}

type Security struct {
	Issuer             string
	JWTExpirationDelta int
	JWTSecret          uuid.UUID
	RefreshSecret      uuid.UUID
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
	var config Config
	err := initBootstrapSettings(&config)
	if err != nil {
		return nil, err
	}

	err = initSecuritySettings(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func initBootstrapSettings(config *Config) error {
	buff, err := ioutil.ReadFile(bootstrap)
	if err != nil {
		return err
	}

	var dc databaseConfig
	err = yaml.Unmarshal(buff, &dc)
	if err != nil {
		return err
	}

	d := Database{}
	err = dc.Database.Decode(&d)
	if err != nil {
		return err
	}

	var ac applicationConfig
	err = yaml.Unmarshal(buff, &ac)
	if err != nil {
		return err
	}

	a := Application{}
	err = ac.Application.Decode(&a)
	if err != nil {
		return err
	}

	var sc serverConfig
	err = yaml.Unmarshal(buff, &sc)
	if err != nil {
		return err
	}

	s := Server{}
	err = sc.Server.Decode(&s)
	if err != nil {
		return err
	}

	var cc cacheConfig
	err = yaml.Unmarshal(buff, &cc)
	if err != nil {
		return err
	}

	c := Cache{}
	err = cc.Cache.Decode(&c)

	config.Application = a
	config.Database = d
	config.Server = s
	config.Cache = c

	return nil
}

func initSecuritySettings(config *Config) error {
	buff, err := ioutil.ReadFile(security)
	if err != nil {
		return err
	}

	s := Security{}
	err = json.Unmarshal(buff, &s)
	if err != nil {
		return err
	}

	config.Security = s

	return nil
}
