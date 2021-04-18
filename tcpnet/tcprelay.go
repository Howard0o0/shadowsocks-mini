package tcpnet

import (
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Howard0o0/shadowsocks-mini/encrypt"
	"github.com/Howard0o0/shadowsocks-mini/tinylog"
)

func Forward(srcConn net.Conn, destConn net.Conn) {
	defer srcConn.Close()
	defer destConn.Close()

	readbuf := make([]byte, 4096)

	for {
		readLen, err := srcConn.Read(readbuf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		if readLen == 0 {
			continue
		}
		tinylog.LogDebug("foward from %v to %v : %v \n", srcConn.RemoteAddr(), destConn.RemoteAddr(), readbuf[0:readLen])
		_, err = destConn.Write(readbuf[0:readLen])
		if err != nil {
			tinylog.LogError("write error:%v", err)
			break
		}

	}

}

const (
	ENCRYPT = iota
	DECRYPT
)

func BuildSSTunnel(left, right encrypt.CipherStreamer) error {

	var err, err1 error
	var wg sync.WaitGroup
	var wait = 5 * time.Second
	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err1 = io.Copy(right, left)
		right.SetReadDeadline(time.Now().Add(wait)) // unblock read on right
	}()
	_, err = io.Copy(left, right)
	left.SetReadDeadline(time.Now().Add(wait)) // unblock read on left
	wg.Wait()
	if err1 != nil && !errors.Is(err1, os.ErrDeadlineExceeded) && !errors.Is(err1, io.EOF) { // requires Go 1.15+
		return err1
	}
	if err != nil && !errors.Is(err, os.ErrDeadlineExceeded) && !errors.Is(err1, io.EOF) {
		return err
	}
	return nil

}
