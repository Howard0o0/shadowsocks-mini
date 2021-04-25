package encrypt

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
)

type Chacha20Codec struct {
	aead  cipher.AEAD
	nonce []byte
}

func NewChacha20Codec(passwd string, salt []byte) (*Chacha20Codec, error) {

	subkey := make([]byte, chacha20poly1305.KeySize)
	key := kdf(passwd, chacha20poly1305.KeySize)
	hkdfSHA1(key, salt, []byte("ss-subkey"), subkey)
	aead, err := chacha20poly1305.New(subkey)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, aead.NonceSize())

	return &Chacha20Codec{aead: aead, nonce: nonce}, nil
}

func (codec Chacha20Codec) NonceSize() int {
	return codec.aead.NonceSize()
}

func (codec Chacha20Codec) Overhead() int {
	return codec.aead.Overhead()
}

func (codec Chacha20Codec) Seal(plaintext []byte) []byte {

	defer increment(codec.nonce)
	return codec.aead.Seal(nil, codec.nonce, plaintext, nil)
}

func (codec Chacha20Codec) Open(ciphertext []byte) ([]byte, error) {

	defer increment(codec.nonce)
	return codec.aead.Open(nil, codec.nonce, ciphertext, nil)
}
