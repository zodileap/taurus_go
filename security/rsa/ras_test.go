package rsa

import (
	"testing"

	tlog "github.com/yohobala/taurus_go/tlog"
)

func TestGenerateKey(t *testing.T) {

	privateKey, publicKey := GenerateKey()

	tlog.Print("",
		tlog.Int("privateKey长度", len(privateKey)),
		tlog.String("privateKey", privateKey),
		tlog.Int("publicKey长度", len(publicKey)),
		tlog.String("publicKey", publicKey),
	)
}
