package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Howard0o0/shadowsocks-mini/core"
)

func main() {

	uri := flag.Bool("uri", false, "generate URI")
	confFile := flag.String("conf", "/etc/ssmini/conf.json", "config file")
	flag.Parse()

	cfg, err := parseConf(*confFile)
	if err != nil {
		fmt.Printf("parse config failed : %s \n", err)
	}
	fmt.Printf("config:\n %s \n", *cfg)
	if err := setLogDir(cfg.Logdir); err != nil {
		fmt.Println("set log dir error : ", err)
		os.Exit(1)
	}

	if *uri {
		base64URI, err := genURI(cfg)
		if err != nil {
			os.Stderr.WriteString(err.Error())
		} else {
			fmt.Println("URI:")
			fmt.Println(base64URI)
		}
		os.Exit(0)
	}

	core.SSServer(cfg.ListenPort, cfg.Method, cfg.Password)
}
