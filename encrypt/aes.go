package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type aesCfbCodec struct {
	password string
	noNeedNonce
}

func NewAESCfbCodec(passwd string) (*aesCfbCodec, error) {
	return &aesCfbCodec{password: string(md5Hash([]byte(passwd)))}, nil
}

func (codec aesCfbCodec) Encode(origData []byte) (encrypted []byte) {
	block, err := aes.NewCipher([]byte(codec.password))
	if err != nil {
		panic(err)
	}
	encrypted = make([]byte, aes.BlockSize+len(origData))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(encrypted[aes.BlockSize:], origData)
	return encrypted
}
func (codec aesCfbCodec) Decode(encrypted []byte) (decrypted []byte) {
	block, _ := aes.NewCipher([]byte(codec.password))
	if len(encrypted) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := encrypted[:aes.BlockSize]
	encrypted = encrypted[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(encrypted, encrypted)
	return encrypted
}
