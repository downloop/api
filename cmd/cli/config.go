package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Token        string `yaml:"token"`
	RefreshToken string `yaml:"refreshToken"`
	Endpoint     string `yaml:"endpoint"`
}

func ReadConfig() (*Config, error) {
	var c Config

	path := configPath()

	y, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(y, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (c Config) Save() error {
	path := configPath()

	y, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path, y, 0644)
}

func configPath() string {
	path := os.Getenv("DOWNLOOP_CONFIG")
	if path == "" {
		usr, _ := user.Current()
		dir := usr.HomeDir
		path = filepath.Join(dir, ".downloop.cfg")
	}
	return path
}
