package encrypt

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/Howard0o0/shadowsocks-mini/cruiser"
	"github.com/Howard0o0/shadowsocks-mini/util"
	"github.com/pkg/errors"
)

const (
	chacha20 = iota
)

const (
	enc = iota
	dec
)

type CipherStream struct {
	encryptor  MetaCipher
	decipherer MetaCipher
	conn       net.Conn
	passwd     string
	method     int
	salt       []byte
}

type CipherStreamer interface {
	io.ReadWriter
	Close() error
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
	SetReadDeadline(t time.Time) error
}

// reserve a factory interface
// so far only support method chacha20
func NewCypherStream(mthd, passwd string, conn net.Conn) (*CipherStream, error) {
	var method int

	switch mthd {
	case "AEAD_CHACHA20_POLY1305":
		method = chacha20
	default:
		method = chacha20
	}
	return &CipherStream{encryptor: nil, decipherer: nil, conn: conn, passwd: passwd, method: method}, nil
}

func (stream *CipherStream) initCipher(sel int, salt []byte) error {
	var cipher MetaCipher
	var err error

	switch stream.method {
	case chacha20:
		cipher, err = NewChacha20Codec(stream.passwd, salt)
	default:
		cipher, err = NewChacha20Codec(stream.passwd, salt)
	}

	if sel == enc {
		stream.encryptor = cipher
	} else if sel == dec {
		stream.decipherer = cipher
	} else {
		return fmt.Errorf("unknown sel : %d \n", sel)
	}

	return err

}

func (stream *CipherStream) LocalAddr() net.Addr {

	return stream.conn.LocalAddr()
}

func (stream *CipherStream) RemoteAddr() net.Addr {

	return stream.conn.RemoteAddr()
}

func (stream *CipherStream) SetReadDeadline(t time.Time) error {

	return stream.conn.SetReadDeadline(t)
}

func (stream *CipherStream) Read(buf []byte) (int, error) {

	erw := util.ErrReadWrite{RW: stream.conn}
	if stream.decipherer == nil {
		salt := make([]byte, pickKeysize(stream.method))
		erw.ReadFull(salt)
		if erw.Err == nil && cruiser.CheckSalt(salt) {
			erw.Err = cruiser.ErrRepeatedSalt
		}
		cruiser.AddSalt(salt)
		erw.Err = stream.initCipher(dec, salt)
	}

	header := make([]byte, 2+stream.decipherer.Overhead())
	erw.ReadFull(header)
	if erw.Err == nil {
		header, erw.Err = stream.decipherer.Open(header)
	}

	payloadlen := util.BytesToInt(header)
	payload := make([]byte, payloadlen+stream.decipherer.Overhead())
	erw.ReadFull((payload))

	if erw.Err == nil {
		payload, erw.Err = stream.decipherer.Open(payload)
	}
	if erw.Err == nil {
		copy(buf[:len(payload)], payload)
	}
	return len(payload), erw.Err

}

func (stream *CipherStream) Write(buf []byte) (int, error) {
	if _, err := stream.devideWrite(buf); err != nil {
		return 0, err
	}
	return len(buf), nil
}

// shadowsocks-whitepaper section 3.3
func (stream *CipherStream) devideWrite(buf []byte) (int, error) {

	maxSize := 16*1024 - 1

	if len(buf) > maxSize {
		if _, err := stream.write(buf[:maxSize]); err != nil {
			return 0, errors.Wrap(err, "partial write failed")
		}
		return stream.Write(buf[maxSize:])
	}
	return stream.write(buf)
}

// AEAD encrypted TCP stream as defined in shadowsocks whitepaper section 3.3
// [encrypted payload length][length tag][encrypted payload][payload tag]
func (stream *CipherStream) buildChunk(payload []byte) []byte {

	payloadLen := util.IntToBytes(len(payload))
	buffer := bytes.NewBuffer(stream.encryptor.Seal(payloadLen))
	buffer.Write(stream.encryptor.Seal(payload))

	return buffer.Bytes()

}

func (stream *CipherStream) write(buf []byte) (int, error) {
	erw := util.ErrReadWrite{RW: stream.conn}

	if stream.encryptor == nil {
		salt := make([]byte, pickKeysize(stream.method))
		io.ReadFull(rand.Reader, salt)
		erw.Err = stream.initCipher(enc, salt)
		erw.Write(salt)
		if erw.Err == nil {
			cruiser.AddSalt(salt)
		}
	}

	erw.Write(stream.buildChunk(buf))

	return len(buf), erw.Err
}

func (stream *CipherStream) Close() error {
	return stream.conn.Close()
}
