package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Howard0o0/shadowsocks-mini/core"
)

func main() {

	if len(os.Args) != 2 {
		prtUsage()
		os.Exit(1)
	}

	cfg, err := parseConf(os.Args[1])
	if err != nil {
		fmt.Printf("parse config failed : %s \n", err)
	}
	fmt.Printf("config:\n %s \n", *cfg)
	if err := setLogDir(cfg.Logdir); err != nil {
		fmt.Println("set log dir error : ", err)
		os.Exit(1)
	}

	core.SSServer(cfg.ListenPort, cfg.Method, cfg.Password)
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
