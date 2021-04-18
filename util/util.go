package util

import (
	"bytes"
	"encoding/binary"
	"io"
)

//整形转换成字节,4字节
func IntToBytes(n int) []byte {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int32
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
