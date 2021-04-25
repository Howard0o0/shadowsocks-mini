package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

type Config struct {
	ListenPort string `json:"listen_port"`
	Password   string `json:"password"`
	Method     string `json:"method"`
	Logdir     string `json:"logdir"`
}

func (cfg Config) String() string {
	str := "\n"
	str += "socks port\t" + cfg.ListenPort + "\n"
	str += "passwd\t" + cfg.Password + "\n"
	str += "method\t" + cfg.Method + "\n"
	str += "logdir\t" + cfg.Logdir + "\n"

	return str
}

func prtUsage() {
	logrus.Info("usage : ssmini path_to_config_file")
}

func parseConf(filename string) (*Config, error) {
	jsonFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var cfg Config
	if err := json.Unmarshal(byteValue, &cfg); err != nil {
		return nil, err
	}

	supportMethods := map[string]bool{
		"AEAD_CHACHA20_POLY1305": true,
	}
	if _, ok := supportMethods[cfg.Method]; !ok {
		return nil, fmt.Errorf("unsupported method : %s ", cfg.Method)
	}

	return &cfg, nil
}
