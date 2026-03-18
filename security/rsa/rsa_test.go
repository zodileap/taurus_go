package rsa

import (
	"encoding/pem"
	"errors"
	"strings"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	privateKey, publicKey, err := GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey 返回了意外错误: %v", err)
	}

	privateBlock, _ := pem.Decode([]byte(privateKey))
	if privateBlock == nil || privateBlock.Type != "RSA PRIVATE KEY" {
		t.Fatalf("私钥 PEM 不合法: %q", privateKey)
	}

	publicBlock, _ := pem.Decode([]byte(publicKey))
	if publicBlock == nil || publicBlock.Type != "RSA PUBLIC KEY" {
		t.Fatalf("公钥 PEM 不合法: %q", publicKey)
	}
}

func TestGenerateKeyReturnsError(t *testing.T) {
	original := marshalPublicKeyDER
	t.Cleanup(func() {
		marshalPublicKeyDER = original
	})

	expectedErr := errors.New("marshal public key failed")
	marshalPublicKeyDER = func(any) ([]byte, error) {
		return nil, expectedErr
	}

	privateKey, publicKey, err := GenerateKey()
	if !errors.Is(err, expectedErr) {
		t.Fatalf("期望错误 %v，实际 %v", expectedErr, err)
	}
	if strings.TrimSpace(privateKey) != "" || strings.TrimSpace(publicKey) != "" {
		t.Fatalf("失败时不应返回密钥内容，实际 private=%q public=%q", privateKey, publicKey)
	}
}
