package core

import (
	"errors"
	"net"

	"github.com/Howard0o0/shadowsocks-mini/cruiser"
	"github.com/Howard0o0/shadowsocks-mini/encrypt"
	"github.com/Howard0o0/shadowsocks-mini/socks5"
	"github.com/Howard0o0/shadowsocks-mini/tcpnet"
	"github.com/Howard0o0/shadowsocks-mini/util"
	"github.com/sirupsen/logrus"
)

func SSServer(port, method, passwd string) {

	listenAddr := ":" + port
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		logrus.Error("listen error : ", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			logrus.Error("accept error : ", err)
			continue
		}

		go func() {

			defer conn.Close()
			cypherStream, err := encrypt.NewCypherStream(method, passwd, conn)
			if err != nil {
				logrus.Error("new cypherstrea error : ", err)
				return
			}
			ssProxy(cypherStream)
		}()
	}
}

func ssProxy(localConn encrypt.CipherStreamer) {
	defer localConn.Close()

	proxyAddr, remoteConn, err := ssConectToRemote(localConn)
	if err != nil {
		logrus.Error("socks5 connect remote error : ", err)
		return
	}
	defer remoteConn.Close()

	logrus.Infof("build tunnel %v<->%v<->%v", localConn.RemoteAddr(), localConn.LocalAddr(), proxyAddr)
	if err := tcpnet.BuildSSTunnel(localConn, remoteConn); err != nil {
		logrus.Errorf("tunnel error %v<->%v<->%v", localConn.RemoteAddr(), localConn.LocalAddr(), proxyAddr)
		logrus.Error(err)
		if errors.Is(err, cruiser.ErrRepeatedSalt) {
			logrus.Warn("repeat salt from : ", localConn.RemoteAddr())
		}
	}
	logrus.Infof("release tunnel %v<->%v<->%v", localConn.RemoteAddr(), localConn.LocalAddr(), proxyAddr)

}

func ssConectToRemote(localConn encrypt.CipherStreamer) (string, net.Conn, error) {

	var proxyAddr string
	erw := util.ErrReadWrite{RW: localConn}

	buf := make([]byte, 1024)

	if n := erw.Read(buf); erw.Err == nil {
		proxyAddr, erw.Err = socks5.ParseAddr(buf[:n])
	}

	if erw.Err != nil {
		// refer to go-shadowsocks2
		// handle active probe, drain illegal localConn to avoid leaking server behavioral features
		// refer to https://gfw.report/blog/gfw_shadowsocks/
		logrus.Warn("suspected host : ", localConn.RemoteAddr())
		util.Suspend(localConn)
		return "", nil, erw.Err
	}

	var remoteConn net.Conn
	if erw.Err == nil {
		remoteConn, erw.Err = net.Dial("tcp", proxyAddr)
	}

	return proxyAddr, remoteConn, erw.Err

}
