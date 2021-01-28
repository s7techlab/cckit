package encryption

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/pkg/errors"
)

// From bccsp/sw/aes
func pkcs7Padding(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs7UnPadding(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > aes.BlockSize || unpadding == 0 {
		return nil, errors.New("Invalid pkcs7 padding (unpadding > aes.BlockSize || unpadding == 0)")
	}

	pad := src[len(src)-unpadding:]
	for i := 0; i < unpadding; i++ {
		if pad[i] != byte(unpadding) {
			return nil, errors.New("Invalid pkcs7 padding (pad[i] != unpadding)")
		}
	}

	return src[:(length - unpadding)], nil
}

func aesCBCEncryptWithIV(IV []byte, key, s []byte) ([]byte, error) {
	if len(s)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid plaintext. It must be a multiple of the block size")
	}

	if len(IV) != aes.BlockSize {
		return nil, errors.New("Invalid IV. It must have length the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(s))
	copy(ciphertext[:aes.BlockSize], IV)

	mode := cipher.NewCBCEncrypter(block, IV)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], s)

	return ciphertext, nil
}

// AESCBCPKCS7Decrypt combines CBC decryption and PKCS7 unpadding
func AESCBCPKCS7Decrypt(key, src []byte) ([]byte, error) {
	// First decrypt
	pt, err := aesCBCDecrypt(key, src)
	if err == nil {
		return pkcs7UnPadding(pt)
	}
	return nil, err
}

func aesCBCDecrypt(key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(src) < aes.BlockSize {
		return nil, errors.New("Invalid ciphertext. It must be a multiple of the block size")
	}
	iv := src[:aes.BlockSize]
	src = src[aes.BlockSize:]

	if len(src)%aes.BlockSize != 0 {
		return nil, errors.New("Invalid ciphertext. It must be a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(src, src)

	return src, nil
}

func EncryptBytes(key, value []byte) ([]byte, error) {
	// IV temporary blank
	return aesCBCEncryptWithIV(make([]byte, 16), key, pkcs7Padding(value))
}

func DecryptBytes(key, value []byte) ([]byte, error) {
	return AESCBCPKCS7Decrypt(key, value)
}
