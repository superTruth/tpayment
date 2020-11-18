package algorithmutils

import (
	"crypto/aes"
	"crypto/cipher"
)

func AESGCMDecrypt(data, key, iv []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	d, err := cipher.NewGCMWithNonceSize(aesCipher, len(iv))
	if err != nil {
		return nil, err
	}
	return d.Open(nil, iv, data, nil)
}

func AESGCMEncrypt(data, key, iv []byte) ([]byte, error) {
	aesCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	d, err := cipher.NewGCMWithNonceSize(aesCipher, len(iv))
	if err != nil {
		return nil, err
	}
	return d.Seal(nil, iv, data, nil), nil
}
