package main

import (
	"fmt"
	"io"
	"net"

	"github.com/Howard0o0/shadowsocks-mini/encrypt"
	"github.com/Howard0o0/shadowsocks-mini/tinylog"
	"github.com/Howard0o0/shadowsocks-mini/util"
)

var codec encrypt.Codec

func ssserverTest() {

	listenAddr := ":1998"
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		tinylog.LogFatal("Listen error : %v \n", err)
	}

	passwd := "2pKR1Kl2jayPXtbwT2GKd129GSOU64Cdv7ejfUZJ8+Ar/wgYc1hjla4fCsBrVB2gkKjf5tFQMbPjZ3tmTp+aZBZLLoS5hafGFzrb7c1F0P18QcriAmqcuyoe6cNlN38Q/iek9ldDOfLXR0BKAMjMsPWYjPHeBfccPHAT4e4t07Ul1aUzMIZTKAmOPRVVx3hCupv4iwf0D5ZfzkSq+cIEElrlO4hpEWKeA5myLGi4rzhR6ughogY1VuxxP4mDNNK82Xr7C8mmdA1tPkh+JimhGiB1xExcxS+rgVl5y02+Npdb725ggmxvG1LBItjntuQB/HIkDNyHrc/6FLEOk90ytA=="
	// codec, err = encrypt.NewCypherBookCodec(passwd)
	// passwd := "howard5279"
	codec, err = encrypt.NewChacha20Codec(passwd)

	if err != nil {
		tinylog.LogFatal("new codec error : %v \n", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			tinylog.LogError("accept error : %v \n", err)
			conn.Close()
			continue
		}

		cs, err := encrypt.NewCypherStream("chacha20", passwd, conn, false)
		if err != nil {
			tinylog.LogError("new cypherstream err : %s \n", err.Error())
			conn.Close()
		}

		go HandleConn(cs, true)
	}
}

func HandleConn(conn encrypt.CipherStreamer, needNonce bool) {

	// plainText := []byte("hello!")
	// cypherText := conn.codec.Encode(plainText)
	// tinylog.LogDebug("cypherText : %v\n", cypherText)
	// decypherText := conn.codec.Decode(cypherText)
	// tinylog.LogDebug("cypherText : %s\n", string(decypherText))

	defer conn.Close()
	// readbuf := make([]byte, 4096)
	var readbuf []byte
	var err error

	for {
		readbuf, err = util.ReadAll(conn)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		if len(readbuf) == 0 {
			continue
		}

		tinylog.LogInfo("receive : %v \n", string(readbuf))
		resp := []byte(fmt.Sprintf("receive : %v", string(readbuf)))
		tinylog.LogDebug("resp : %v\n", string(resp))

		_, err = conn.Write(resp)
		if err != nil {
			tinylog.LogError("write error:%v", err)
			break
		}

	}
}
