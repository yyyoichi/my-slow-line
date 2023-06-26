package services_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"testing"

	"github.com/SherClockHolmes/webpush-go"
)

func TestVapid(t *testing.T) {
	privateKey, publicKey, err := webpush.GenerateVAPIDKeys()
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("private-key: %s\n", privateKey)
	fmt.Printf("public-key: %s\n", publicKey)
}

func TestP256key(t *testing.T) {
	// P-256 gen
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
		return
	}

	// DER encode
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		t.Error(err)
		return
	}

	// Base64URL encode
	publicKeyBase64 := base64.URLEncoding.EncodeToString(publicKeyBytes)

	privateKeyBit, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		t.Error(err)
	}
	block := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBit,
	}
	fmt.Println("private-key:")
	fmt.Println(string(pem.EncodeToMemory(block)))

	fmt.Println("public-key:")
	fmt.Println(publicKeyBase64)
}
