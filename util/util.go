package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
)

type ErrReadWrite struct {
	RW  io.ReadWriter
	Err error
}

func (erw *ErrReadWrite) ReadFull(dst []byte) {
	if erw.Err != nil {
		return
	}
	_, erw.Err = io.ReadFull(erw.RW, dst)
}

func (erw *ErrReadWrite) Read(dst []byte) (n int) {
	if erw.Err != nil {
		return 0
	}
	n, erw.Err = erw.RW.Read(dst)
	return n

}

func (erw *ErrReadWrite) Write(data []byte) {
	if erw.Err != nil {
		return
	}
	_, erw.Err = erw.RW.Write(data)

}

//整形转换成字节,4字节
func IntToBytes(n int) []byte {
	x := int16(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int16
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}

func ReadAll(r io.Reader) ([]byte, error) {
	const bufLen = 1024
	buf := make([]byte, bufLen)
	byteBuffer := bytes.NewBuffer([]byte{})

	for {
		n, err := r.Read(buf)
		byteBuffer.Write(buf[:n])

		if err != nil && err != io.EOF {
			return []byte{}, err
		}

		if n < bufLen {
			break
		}
	}

	return byteBuffer.Bytes(), nil
}

func Suspend(r io.Reader) {
	io.Copy(ioutil.Discard, r)
}

// dir should be an absolute path
func CreateDir(dir string) error {
	_, err := os.Stat(dir)

	if err == nil {
		//directory exists
		return nil
	}

	err2 := os.MkdirAll(dir, 0755)
	if err2 != nil {
		return err2
	}

	return nil
}

func GetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("Can not find the ip address!")

}
