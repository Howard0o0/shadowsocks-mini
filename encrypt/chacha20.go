package encrypt

import (
	"crypto/cipher"
	"crypto/sha256"

	"golang.org/x/crypto/chacha20poly1305"
)

type Chacha20Codec struct {
	aead  cipher.AEAD
	nonce []byte
}

func NewChacha20Codec(passwd string) (*Chacha20Codec, error) {
	key := sha256.Sum256([]byte(passwd))
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}

	return &Chacha20Codec{aead: aead}, nil
}

func (codec *Chacha20Codec) SetNonce(nonce []byte) {
	codec.nonce = nonce
}

func (codec Chacha20Codec) NonceSize() int {
	return chacha20poly1305.NonceSizeX
}

func (codec Chacha20Codec) NonceEmpty() bool {
	return len(codec.nonce) != codec.NonceSize()
}

func (codec Chacha20Codec) Encode(origData []byte) []byte {
	if codec.NonceEmpty() {
		return []byte{}
	}

	return codec.aead.Seal(nil, codec.nonce, origData, nil)
}

func (codec Chacha20Codec) Decode(encrypted []byte) (decrypted []byte) {
	if codec.NonceEmpty() {
		return []byte{}
	}

	decypherText, _ := codec.aead.Open(nil, codec.nonce, encrypted, nil)
	return decypherText
}
