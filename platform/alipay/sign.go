package alipay

import (
	"bytes"
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
	"strings"
)

func OpenSSLSign(data []byte, signType, privateKey string) (string, error) {
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
	pk, err := ParsePrivateKey(FormatPrivateKey(privateKey))
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

func OpenSSLVerify(data, sign []byte, signType, publicKey string) error {
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
	pk, err := ParsePublicKey(FormatPublicKey(publicKey))
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pk, hType, d, sign)
}

func ParsePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		err := errors.New("private key error:" + string(privateKey))
		return nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key, err
}

func ParsePublicKey(publicKey []byte) (key *rsa.PublicKey, err error) {
	var block *pem.Block
	block, _ = pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}

	var pubInterface interface{}
	pubInterface, err = x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("public key error")
	}

	return key, err
}

func FormatPublicKey(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN PUBLIC KEY-----", "-----END PUBLIC KEY-----")
}

func FormatPrivateKey(raw string) (result []byte) {
	return formatKey(raw, "-----BEGIN RSA PRIVATE KEY-----", "-----END RSA PRIVATE KEY-----")
}

func formatKey(raw, prefix, suffix string) (result []byte) {
	if raw == "" {
		return nil
	}
	raw = strings.Replace(raw, prefix, "", 1)
	raw = strings.Replace(raw, suffix, "", 1)
	raw = strings.Replace(raw, " ", "", -1)
	raw = strings.Replace(raw, "\n", "", -1)
	raw = strings.Replace(raw, "\r", "", -1)
	raw = strings.Replace(raw, "\t", "", -1)

	var ll = 64
	var sl = len(raw)
	var c = sl / ll
	if sl%ll > 0 {
		c = c + 1
	}

	var buf bytes.Buffer
	buf.WriteString(prefix + "\n")
	for i := 0; i < c; i++ {
		var b = i * ll
		var e = b + ll
		if e > sl {
			buf.WriteString(raw[b:])
		} else {
			buf.WriteString(raw[b:e])
		}
		buf.WriteString("\n")
	}
	buf.WriteString(suffix)
	return buf.Bytes()
}
