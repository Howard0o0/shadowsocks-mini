package encrypt

import (
	"crypto/md5"
	"crypto/sha1"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/hkdf"
)

type MetaCipher interface {
	Seal(plaintext []byte) []byte
	Open(ciphertext []byte) ([]byte, error)
	NonceSize() int
	Overhead() int
}

// key-derivation function from original Shadowsocks
func kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}

func hkdfSHA1(secret, salt, info, outkey []byte) {
	r := hkdf.New(sha1.New, secret, salt, info)
	if _, err := io.ReadFull(r, outkey); err != nil {
		panic(err) // should never happen
	}
}

// increment little-endian encoded unsigned integer b. Wrap around on overflow.
func increment(b []byte) {
	for i := range b {
		b[i]++
		if b[i] != 0 {
			return
		}
	}
}

func pickKeysize(mthd int) int {
	switch mthd {
	case chacha20:
		return chacha20poly1305.KeySize
	default:
		return chacha20poly1305.KeySize
	}
}
