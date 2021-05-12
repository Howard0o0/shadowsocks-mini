package util

import (
	"bytes"
	"encoding/binary"
	"io"
	"io/ioutil"
	"net/http"
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
	responseClient, errClient := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if errClient != nil {
		os.Stderr.WriteString("获取外网 IP 失败，请检查网络\n")
		return "", errClient
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer responseClient.Body.Close()

	body, _ := ioutil.ReadAll(responseClient.Body)
	ip := string(body)
	return ip, nil

}
