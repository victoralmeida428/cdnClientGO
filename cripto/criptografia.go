package cripto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
)

type Criptografia struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func New() *Criptografia {
	var cript Criptografia
	cript.loadPublicKey("./public_key.pem")
	cript.loadPrivateKey("./private_key.pem")
	return &cript
}

func (g *Criptografia) loadPrivateKey(path string) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic("failed to decode PEM block")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	g.privateKey = privateKey
}

func (g *Criptografia) loadPublicKey(path string) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic("failed to decode PEM block")
	}

	// Alterado para ParsePKIXPublicKey
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	publicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		panic("not an RSA public key")
	}
	g.publicKey = publicKey
}

func (c Criptografia) Encode(s string) (string, error) {
	msg := []byte(s)
	hash := sha256.New()

	ciphertext, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		c.publicKey,
		msg,
		nil,
	)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c Criptografia) Decode(s string) (string, error) {
	// Decodifica a string base64
	ciphertext, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	hash := sha256.New()

	plaintext, err := rsa.DecryptOAEP(
		hash,
		rand.Reader,
		c.privateKey,
		ciphertext,
		nil,
	)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
