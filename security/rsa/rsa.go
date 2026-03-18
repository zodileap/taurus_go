package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
)

var (
	generatePrivateKey  = func(reader io.Reader, bits int) (*rsa.PrivateKey, error) { return rsa.GenerateKey(reader, bits) }
	marshalPublicKeyDER = x509.MarshalPKIXPublicKey
)

// 生成公钥和私钥
//
// 返回：
//   - private：私钥
//   - public：公钥
//   - err：错误信息
func GenerateKey() (private string, public string, err error) {
	// Generate private key
	privateKey, err := generatePrivateKey(rand.Reader, 1024)
	if err != nil {
		return "", "", err
	}

	// Encode private key to PEM format
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Generate public key
	publicKey := &privateKey.PublicKey

	// Encode public key to PEM format
	publicKeyDER, err := marshalPublicKeyDER(publicKey)
	if err != nil {
		return "", "", err
	}

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyDER,
	})

	return string(privateKeyPEM), string(publicKeyPEM), nil
}
