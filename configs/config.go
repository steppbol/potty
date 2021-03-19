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

type Security struct {
	Issuer             string
	JWTExpirationDelta int
	Secret             uuid.UUID
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
	b, err := initBootstrapSettings()
	if err != nil {
		return nil, err
	}

	s, err := initSecuritySettings()
	if err != nil {
		return nil, err
	}

	return &Config{
		Application: b.Application,
		Database:    b.Database,
		Server:      b.Server,
		Security:    s.Security,
	}, nil
}

func initBootstrapSettings() (*Config, error) {
	buff, err := ioutil.ReadFile(bootstrap)
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

	var sc serverConfig
	err = yaml.Unmarshal(buff, &sc)

	if err != nil {
		return nil, err
	}

	s := Server{}
	err = sc.Server.Decode(&s)

	return &Config{
		Application: a,
		Database:    d,
		Server:      s,
	}, nil
}

func initSecuritySettings() (*Config, error) {
	buff, err := ioutil.ReadFile(security)
	if err != nil {
		return nil, err
	}

	s := Security{}
	err = json.Unmarshal(buff, &s)
	if err != nil {
		return nil, err
	}

	return &Config{
		Security: s,
	}, nil
}
