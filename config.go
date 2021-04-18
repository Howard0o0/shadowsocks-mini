package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/Howard0o0/shadowsocks-mini/tinylog"
)

type Config struct {
	Server       string `json:"server"`
	ServerPort   string `json:"server_port"`
	LocalAddress string `json:"local_address"`
	LocalPort    string `json:"local_port"`
	Password     string `json:"password"`
	Method       string `json:"method"`
	Timeout      int    `json:"timeout"`
}

func (cfg Config) String() string {
	str := "\n"
	str += "server\t" + cfg.Server + "\n"
	str += "serverport\t" + cfg.ServerPort + "\n"
	str += "localaddr\t" + cfg.LocalAddress + "\n"
	str += "localport\t" + cfg.LocalPort + "\n"
	str += "passwd\t" + cfg.Password + "\n"
	str += "method\t" + cfg.Method + "\n"
	str += "timeout\t" + strconv.Itoa(cfg.Timeout) + "\n"

	return str
}

func prtUsage() {
	tinylog.LogError("usage : ssmini {-c|-s} -conf {ssmini.json}")
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
		"chacha20":    true,
		"cypherbook":  true,
		"aes-256-cfb": true,
	}
	if _, ok := supportMethods[cfg.Method]; !ok {
		return nil, fmt.Errorf("unsupported method : %s ", cfg.Method)
	}

	return &cfg, nil
}
