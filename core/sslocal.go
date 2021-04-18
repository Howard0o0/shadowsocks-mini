package core

import (
	"net"

	"github.com/Howard0o0/shadowsocks-mini/encrypt"
	"github.com/Howard0o0/shadowsocks-mini/tcpnet"
	"github.com/Howard0o0/shadowsocks-mini/tinylog"
)

func SSLocal(localPort, remotePort, remoteIP, method, passwd string) {

	listenAddr := ":" + localPort
	remoteAddr := remoteIP + ":" + remotePort
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		tinylog.LogFatal("Listen error : %v \n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			tinylog.LogError("accept error : %v \n", err)
			conn.Close()
			continue
		}

		go func() {
			defer conn.Close()
			remoteConn, err := net.Dial("tcp", remoteAddr)
			if err != nil {
				tinylog.LogError("accept error : %v \n", err)
				return
			}
			defer remoteConn.Close()

			cypherStream, err := encrypt.NewCypherStream(method, passwd, remoteConn, true)
			if err != nil {
				tinylog.LogError("accept error : %v \n", err)
				return
			}

			tinylog.LogInfo("%v <--tunnel built--> %v \n", conn.LocalAddr(), remoteAddr)
			if err := tcpnet.BuildSSTunnel(cypherStream, conn); err != nil {
				tinylog.LogError("%v <--tunnel error--> %v \n", conn.LocalAddr(), remoteAddr)
				tinylog.LogError("%s \n", err)
			}
			tinylog.LogInfo("%v <--tunnel destroyed--> %v \n", conn.LocalAddr(), remoteAddr)
		}()
	}

}
