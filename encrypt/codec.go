package encrypt

import (
	"crypto/md5"
)

type Codec interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
	SetNonce([]byte)
	NonceEmpty() bool
	NonceSize() int
}

func md5Hash(data []byte) []byte {
	hash := md5.New()
	hash.Write(data)
	return hash.Sum(nil)
}

type noNeedNonce struct {
}

func (codec noNeedNonce) SetNonce([]byte) {

}

func (codec noNeedNonce) NonceSize() int {
	return 0
}

func (codec noNeedNonce) NonceEmpty() bool {
	return false
}
