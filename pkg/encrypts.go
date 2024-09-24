package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// Cryptor 加/解密
type Cryptor interface {
	Encrypt(clearText string) string
	Decrypt(ciphertext string) string
}

type CryptoImpl struct {
	key []byte
}

func NewCrypto(key string) CryptoImpl {
	return CryptoImpl{key: StringToBytes(key)}
}

func (e CryptoImpl) Encrypt(clearText string) string {
	clearTextByte := StringToBytes(clearText)
	c, err := aes.NewCipher(e.key)
	if err != nil {
		panic(err)
	}
	cipherByte := make([]byte, aes.BlockSize+len(clearTextByte))
	iv := cipherByte[:aes.BlockSize]
	_, err = io.ReadFull(rand.Reader, iv)
	if err != nil {
		panic(err)
	}
	cfb := cipher.NewCFBEncrypter(c, iv)
	cfb.XORKeyStream(cipherByte[aes.BlockSize:], clearTextByte)
	return hex.EncodeToString(cipherByte)
}

func (e CryptoImpl) Decrypt(ciphertext string) string {
	cipherByte, err := hex.DecodeString(ciphertext)
	if err != nil {
		panic(err)
	}
	c, err := aes.NewCipher(e.key)
	if err != nil {
		panic(err)
	}
	iv := cipherByte[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(c, iv)
	plainByte := make([]byte, len(cipherByte)-aes.BlockSize)
	stream.XORKeyStream(plainByte, cipherByte[aes.BlockSize:])
	return BytesToString(plainByte)
}

func Md5(str string) string {
	hash := md5.New()
	hash.Write(StringToBytes(str))
	return hex.EncodeToString(hash.Sum(nil))
}
