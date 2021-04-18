package encrypt

import (
	"encoding/base64"
	"math/rand"
	"time"

	"github.com/pkg/errors"
)

type CypherBookCodec struct {
	password     string
	cypherDict   *CypherDict // 编码本
	decypherDict *CypherDict // 解码本
	noNeedNonce
}

const cypherDictSize = 256

type CypherDict [cypherDictSize]byte

// 用base64编码把CypherDict翻译成字符串(密码)
func (cypherDict *CypherDict) String() string {
	return base64.StdEncoding.EncodeToString(cypherDict[:])
}

func generateCypherDict() *CypherDict {
	rand.Seed(time.Now().Unix())
	intArr := rand.Perm(cypherDictSize)
	cypherDict := &CypherDict{}
	for i, v := range intArr {
		cypherDict[i] = byte(v)
	}
	return cypherDict
}

func isValidCypherDict(cypherDict *CypherDict) bool {
	for i, v := range cypherDict {
		if int(v) == i {
			return false
		}
	}
	return true
}

// 生成随机密码
func GeneratePasswd() string {
	cypherDict := generateCypherDict()
	for !isValidCypherDict(cypherDict) {
		cypherDict = generateCypherDict()
	}
	return cypherDict.String()
}

// 将密码转译为密码本
func passwdToCypherDict(passwd *string) (*CypherDict, error) {
	byteArr, err := base64.StdEncoding.DecodeString(*passwd)
	if err != nil || len(byteArr) != cypherDictSize {
		return nil, errors.New("invalid password")
	}

	cypherDict := CypherDict{}
	copy(cypherDict[:], byteArr)
	byteArr = nil
	return &cypherDict, nil
}

func NewCypherBookCodec(password string) (*CypherBookCodec, error) {
	cypherDict, err := passwdToCypherDict(&password)
	if err != nil {
		return nil, err
	}

	decypherDict := &CypherDict{}
	for i, v := range cypherDict {
		decypherDict[v] = byte(i)
	}

	return &CypherBookCodec{
		password:     password,
		cypherDict:   cypherDict,
		decypherDict: decypherDict,
	}, nil
}

func (codec CypherBookCodec) Encode(plainText []byte) []byte {

	byteArr := make([]byte, len(plainText))
	copy(byteArr, plainText)
	for i, v := range byteArr {
		byteArr[i] = codec.cypherDict[v]
	}
	return byteArr
}

func (codec CypherBookCodec) Decode(encrypted []byte) []byte {

	byteArr := make([]byte, len(encrypted))
	copy(byteArr, encrypted)
	for i, v := range byteArr {
		byteArr[i] = codec.decypherDict[v]
	}
	return byteArr
}
