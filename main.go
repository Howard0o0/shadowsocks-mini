package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Howard0o0/shadowsocks-mini/core"
	"github.com/Howard0o0/shadowsocks-mini/tinylog"
)

func main() {

	// ssserverTest()

	tinylog.SetLevel(tinylog.InfoLevel)

	identity, confFile, status := parseFlag()
	if !status {
		prtUsage()
		os.Exit(1)
	}

	cfg, err := parseConf(confFile)
	if err != nil {
		tinylog.LogError("parse config failed : %s \n", err)
	}
	fmt.Printf("config:\n %s \n", *cfg)

	if identity == "server" {
		fmt.Println("server mode")
		core.SSServer(cfg.ServerPort, cfg.Method, cfg.Password)
	} else {
		tinylog.LogInfo("client mode\n")
		core.SSLocal(cfg.LocalPort, cfg.ServerPort, cfg.Server, cfg.Method, cfg.Password)
	}
}

func parseFlag() (identity, confFile string, status bool) {
	var local, server bool

	flag.BoolVar(&local, "c", false, "ss local mode")
	flag.BoolVar(&server, "s", false, "ss server mode")
	flag.StringVar(&confFile, "conf", "", "config file absolute path")
	flag.Parse()

	if len(confFile) == 0 || (!local && !server) {
		return "", "", false
	}

	if local {
		identity = "local"
	} else {
		identity = "server"
	}

	return identity, confFile, true

}
