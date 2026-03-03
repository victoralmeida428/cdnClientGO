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

func New(public, private string) *Criptografia {
	c := &Criptografia{}
	c.loadPublicKey(public)
	c.loadPrivateKey(private)
	return c
}

func (g *Criptografia) loadPrivateKey(path string) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic("falha ao decodificar PEM")
	}

	// Tenta PKCS1, se falhar tenta PKCS8 (o seu formato atual)
	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
		g.privateKey = key
	} else {
		key8, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			panic("formato de chave privada não suportado")
		}
		g.privateKey = key8.(*rsa.PrivateKey)
	}
}

func (g *Criptografia) loadPublicKey(path string) {
	keyBytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic("falha ao decodificar PEM")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	g.publicKey = pub.(*rsa.PublicKey)
}

func (c *Criptografia) Encode(s string) (string, error) {
	// Importante: SHA256 deve ser instanciado aqui
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, c.publicKey, []byte(s), nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (c *Criptografia) Decode(s string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	// Importante: SHA256 deve ser instanciado aqui também
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, c.privateKey, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}