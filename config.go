package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	TLS         TLSConfig     `yaml:"tls"`
	GitHub      GitHubConfig  `yaml:"github"`
	Discord     DiscordConfig `yaml:"discord"`
	Git         GitConfig     `yaml:"git"`
	ID          string        `yaml:"id"`
	Description string        `yaml:"description"`
	Projects    []string      `yaml:"projects"`
}

type GitHubConfig struct {
	Port   uint16 `yaml:"port"`
	URI    string `yaml:"uri"`
	Secret string `yaml:"secret"`
}

type GitConfig struct {
	Path string `yaml:"path"`
}

type DiscordConfig struct {
	Token         string `yaml:"token"`
	LogChannel    string `yaml:"log_channel"`
	EventChannel  string `yaml:"event_channel"`
	StatusChannel string `yaml:"status_channel"`
}

type TLSConfig struct {
	Cert string `yaml:"cert"`
	Key  string `yaml:"key"`
}

func (c *Config) Init(filename string) error {
	log.Infof("Reading configuration from %s", filename)
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Config not found: %s", err.Error())
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("Couldn't parse yaml: %s", err.Error())
	}

	return nil
}
