package alipay

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"hash"
)

func Sign(data []byte, signType, privateKey string) (string, error) {
	var h hash.Hash
	var hType crypto.Hash
	if signType == SignTypeRSA {
		h = sha1.New()
		hType = crypto.SHA1
	} else {
		h = sha256.New()
		hType = crypto.SHA256
	}
	h.Write(data)
	d := h.Sum(nil)
	pk, err := ParsePrivateKey(privateKey)
	if err != nil {
		return "", err
	}
	bs, err := rsa.SignPKCS1v15(rand.Reader, pk, hType, d)

	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(bs)
	return signature, nil
}

func ParsePrivateKey(privateKey string) (pk *rsa.PrivateKey, err error) {
	block, _ := pem.Decode([]byte(privateKey))
	if block == nil {
		err = errors.New("私钥格式错误1:" + privateKey)
		return nil, err
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err == nil {
			pk = rsaPrivateKey
		}
	default:
		err = errors.New("私钥格式错误2:" + privateKey)
	}
	return
}
