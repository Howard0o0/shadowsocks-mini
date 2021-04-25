package tcpnet

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/Howard0o0/shadowsocks-mini/encrypt"
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
