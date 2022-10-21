package config

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"
)

type Server struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

type DB struct {
	Hosts    []string `yaml:"hosts"`
	Ports    []int    `yaml:"ports"`
	Username string   `yaml:"user"`
	Password string   `yaml:"password"`
	NameDB   string   `yaml:"dbname"`
}

type Config struct {
	Server   Server `yaml:"server"`
	Database DB     `yaml:"db"`
}

var (
	once sync.Once
	cfg  *Config
)

func GetConfigYML(configYML string) *Config {
	once.Do(func() {
		file, err := os.Open(configYML)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		if err = yaml.NewDecoder(file).Decode(&cfg); err != nil {
			log.Fatal(err)
		}
		cfg.Database.Password = os.Getenv("DB_PASSWORD")
		printConfig(cfg)
	})
	return cfg
}

func printConfig(cfg *Config) {
	jscfg, _ := json.MarshalIndent(cfg, "", "	")
	logrus.Info(string(jscfg))
	// fmt.Println(string(jscfg))
}
