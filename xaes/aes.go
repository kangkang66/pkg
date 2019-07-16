package xaes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

type Aes struct {
	key 	[]byte
}

func NewAes(key string) (*Aes) {
	fmt.Println("NewAes key:",key)
	sum := sha256.Sum256([]byte(key))
	return &Aes{
		key:sum[:],
	}
}

func (a *Aes) Encrypt(encodeBytes []byte) (val string, err error) {
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return
	}
	blockSize := block.BlockSize()
	encodeBytes = a.pkCS5Padding(encodeBytes, blockSize)

	iv := make([]byte, blockSize)
	_,err = rand.Read(iv)
	if err != nil {
		return
	}

	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	iv = append(iv,crypted...)
	val = base64.StdEncoding.EncodeToString(iv)
	return
}


func (a *Aes) pkCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (a *Aes) Decrypt(decodeStr string) (origData []byte,err error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(decodeStr)
	if err != nil {
		return
	}
	block, err := aes.NewCipher(a.key)
	if err != nil {
		return nil, err
	}
	iv := decodeBytes[:block.BlockSize()]
	decodeBytes = decodeBytes[block.BlockSize():]

	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData = make([]byte, len(decodeBytes))
	blockMode.CryptBlocks(origData, decodeBytes)
	origData = a.pkCS5UnPadding(origData)
	return
}

func (a *Aes) pkCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}



