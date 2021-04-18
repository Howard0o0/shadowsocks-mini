package core

import (
	"fmt"
	"net"

	"github.com/Howard0o0/ssmini/encrypt"
	"github.com/Howard0o0/ssmini/procotol"
	"github.com/Howard0o0/ssmini/tcpnet"
	"github.com/Howard0o0/ssmini/tinylog"
	"github.com/pkg/errors"
)

func SSServer(port, method, passwd string) {

	listenAddr := ":" + port
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		tinylog.LogFatal("listen error : %v \n", err)
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
			cypherStream, err := encrypt.NewCypherStream(method, passwd, conn, false)
			if err != nil {
				tinylog.LogError("new cypherstrea error : %v \n", err)
				return
			}
			ssProxy(cypherStream)
		}()
	}
}

func ssProxy(localConn encrypt.CipherStreamer) {
	defer localConn.Close()

	if err := ssHandshake(localConn); err != nil {
		tinylog.LogError("socks5 auth error : %v\n", err)
		return
	}

	proxyAddr, remoteConn, err := ssConectToRemote(localConn)
	if err != nil {
		tinylog.LogError("socks5 connect remote error : %v\n", err)
		return
	}
	defer remoteConn.Close()

	tinylog.LogInfo("%v <--tunnel built--> %v \n", localConn.LocalAddr(), proxyAddr)
	if err := tcpnet.BuildSSTunnel(localConn, remoteConn); err != nil {
		tinylog.LogError("%v <--tunnel error--> %v \n", localConn.LocalAddr(), proxyAddr)
		tinylog.LogError("%s \n", err)
	}
	tinylog.LogInfo("%v <--tunnel destroyed--> %v \n", localConn.LocalAddr(), proxyAddr)

}

func ssHandshake(localConn encrypt.CipherStreamer) error {

	var err error
	n := 0
	buf := make([]byte, 1024)

	if n, err = localConn.Read(buf); err != nil {
		return errors.New("read socks5 auth header error")
	}
	if !procotol.Socks5Auth(buf[:n]) {
		return fmt.Errorf("socks5 auth header illegal : %v", buf)
	}
	resp := procotol.BuildSocks5AuthOkRsp()
	if n, err := localConn.Write(resp); err != nil || n != len(resp) {
		return errors.New("socks5 auth response error")
	}

	return nil
}

func ssConectToRemote(localConn encrypt.CipherStreamer) (string, net.Conn, error) {

	var err error
	n := 0
	buf := make([]byte, 1024)

	if n, err = localConn.Read(buf); err != nil {
		return "", nil, errors.Wrap(err, "read socks5 remote addr error")
	}
	var proxyAddr string
	if proxyAddr, err = procotol.Socks5ReadProxyAddr(buf[:n]); err != nil {
		return "", nil, errors.Wrap(err, "socks5 proxy address illegal  ")
	}

	var remoteConn net.Conn
	for try := 0; try < 3; try++ {
		remoteConn, err = net.Dial("tcp", proxyAddr)
	}
	if err != nil {
		return "", nil, errors.Wrap(err, "connect to proxy address error ")
	}

	resp := procotol.BuildSocks5ConnectOkRsp()
	if n, err := localConn.Write(resp); err != nil || n != len(resp) {
		return "", nil, errors.New("socks5 connect response error")
	}

	return proxyAddr, remoteConn, nil

}
