package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Howard0o0/shadowsocks-mini/util"
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
	str += "passwd\t\t" + cfg.Password + "\n"
	str += "method\t\t" + cfg.Method + "\n"
	str += "logdir\t\t" + cfg.Logdir + "\n"

	return str
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

func genURI(cfg *Config) (string, error) {

	// ss://AEAD_CHACHA20_POLY1305:your-password@:8488
	ip, err := util.GetIp()
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer([]byte{})
	buf.WriteString("chacha20-ietf-poly1305")
	buf.WriteString(":")
	buf.WriteString(cfg.Password)
	buf.WriteString("@")

	buf.WriteString(ip)
	buf.WriteString(":")
	buf.WriteString(cfg.ListenPort)
	logrus.Infoln(buf.String())
	return "ss://" + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
