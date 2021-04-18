package encrypt

import (
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Howard0o0/shadowsocks-mini/tinylog"
	"github.com/Howard0o0/shadowsocks-mini/util"
	"github.com/pkg/errors"
)

type CypherStream struct {
	codec Codec
	conn  net.Conn
}

type CipherStreamer interface {
	io.Reader
	io.Writer
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetReadDeadline(t time.Time) error
}

// 工厂方法
func NewCypherStream(cypherType, passwd string, conn net.Conn, makeNonce bool) (*CypherStream, error) {
	var codec Codec
	var err error

	switch cypherType {
	case "aes-256-cfb":
		codec, err = NewAESCfbCodec(passwd)
	case "cypherbook":
		codec, err = NewCypherBookCodec(passwd)
	case "chacha20":
		codec, err = NewChacha20Codec(passwd)
	default:
		return nil, fmt.Errorf("unknown cypher type : %s", cypherType)
	}

	if makeNonce {
		tinylog.LogDebug("setting nonce\n")
		nonce := make([]byte, codec.NonceSize())
		if _, err := rand.Read(nonce); err != nil {
			return nil, err
		}
		if n, err := conn.Write(nonce); err != nil || n != len(nonce) {
			return nil, errors.Wrap(err, "write nonce failed")
		}
		tinylog.LogDebug("nonce : %v\n", nonce)
		codec.SetNonce(nonce)
	}

	if err != nil {
		return nil, err
	}

	return &CypherStream{codec: codec, conn: conn}, nil
}

func (stream *CypherStream) LocalAddr() net.Addr {
	return stream.conn.LocalAddr()
}

func (stream *CypherStream) RemoteAddr() net.Addr {
	return stream.conn.RemoteAddr()
}
func (stream *CypherStream) SetReadDeadline(t time.Time) error {
	return stream.conn.SetReadDeadline(t)
}

func (stream *CypherStream) Read(readBuf []byte) (int, error) {

	if stream.codec.NonceEmpty() {
		tinylog.LogDebug("receiving nonce\n")
		nonce := make([]byte, stream.codec.NonceSize())
		if _, err := io.ReadFull(stream.conn, nonce); err != nil {
			return 0, errors.Wrap(err, "can't read nonce from stream")
		}
		tinylog.LogDebug("nonce : %v \n", nonce)
		stream.codec.SetNonce(nonce)
	}

	//先读报头(msgLen,4B),再根据msgLen读取指定长度的msg
	header := make([]byte, 4)
	if _, err := io.ReadFull(stream.conn, header); err != nil {
		return 0, errors.Wrap(err, "read msgLen error ")
	}

	msgLen := util.BytesToInt(header)
	buf := make([]byte, msgLen)
	_, err := io.ReadFull(stream.conn, buf)
	if err != nil {
		return 0, err
	}

	decypherBuff := stream.codec.Decode(buf)
	copy(readBuf[:len(decypherBuff)], decypherBuff)
	tinylog.LogDebug("read : %s \n", string(readBuf[:len(decypherBuff)]))
	return len(decypherBuff), err
}

func (stream *CypherStream) Write(buf []byte) (int, error) {
	if _, err := stream.devideWrite(buf); err != nil {
		return 0, err
	}
	return len(buf), nil
}

func (stream *CypherStream) devideWrite(buf []byte) (int, error) {
	// todo
	// io.Copy()中，缓存区的长度为32k,buf是明文，加密后长度为是16B的整数倍
	// 所以密文的长度会比明文的长度要长
	// io.Copy的实现中，缓冲区的长度为32kB，如果此处buf长度刚好32kB，io.Copy的缓存区会溢出
	// 如果buf过大，需要递归的分为小包

	maxSize := 30 * 1024

	if len(buf) > maxSize {
		if _, err := stream.write(buf[:maxSize]); err != nil {
			return 0, errors.Wrap(err, "partial write failed")
		}
		return stream.Write(buf[maxSize:])
	}
	return stream.write(buf)
}

func (stream *CypherStream) write(buf []byte) (int, error) {
	tinylog.LogDebug("write : %s \n", string(buf))
	cypherText := stream.codec.Encode(buf)

	// 添加报头，解决tcp粘包问题
	msgLen := util.IntToBytes(len(cypherText))
	if _, err := stream.conn.Write(msgLen); err != nil {
		return 0, errors.Wrap(err, "write ss header error")
	}
	if _, err := stream.conn.Write(cypherText); err != nil {
		return 0, errors.Wrap(err, "write not complete")
	}
	return len(buf), nil
}

func (stream *CypherStream) Close() error {
	return stream.conn.Close()
}
