package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var PUBLIC_KEY_PATH = "private.pem"
var PRIVATE_KEY_PATH = "public.pem"

type Crypt struct {
	PublicKey    *rsa.PublicKey
	PrivateKey   *rsa.PrivateKey
	PublicKeyPem []byte
}

func KeysExists() bool {
	_, err := os.Stat(PRIVATE_KEY_PATH)
	privateOk := !errors.Is(err, os.ErrNotExist)
	_, err = os.Stat(PUBLIC_KEY_PATH)
	publicOk := !errors.Is(err, os.ErrNotExist)
	return privateOk && publicOk
}

func CreateKeyPair() {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	err = os.WriteFile(PRIVATE_KEY_PATH, privateKeyPEM, 0644)
	if err != nil {
		panic(err)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	err = os.WriteFile(PUBLIC_KEY_PATH, publicKeyPEM, 0644)
	if err != nil {
		panic(err)
	}
}

func (crypt *Crypt) loadKeys() {
	publicKeyPEM, err := os.ReadFile(PUBLIC_KEY_PATH)
	if err != nil {
		panic(err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyPEM)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	privateKeyPEM, err := os.ReadFile(PRIVATE_KEY_PATH)
	if err != nil {
		panic(err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyPEM)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		panic(err)
	}

	crypt.PublicKey = publicKey.(*rsa.PublicKey)
	crypt.PrivateKey = privateKey
	crypt.PublicKeyPem = publicKeyPEM
}

// func (c Crypt) encrypt(plaintext []byte) string {
// 	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, c.PublicKey, plaintext)
// 	if err != nil {
// 		panic(err)
// 	}
// 	// fmt.Printf("Encrypted: %x\n", ciphertext)
// 	encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
// 	return encodedCiphertext
// }

func (c Crypt) sign(message []byte) []byte {
	hashed := sha256.Sum256(message)
	// Sign the hashed message
	signature, err := rsa.SignPKCS1v15(rand.Reader, c.PrivateKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}

	// encodedCiphertext := base64.StdEncoding.EncodeToString(ciphertext)
	return signature
}

func (c Crypt) decrypt(ciphertext []byte) []byte {
	plaintext, err := rsa.DecryptPKCS1v15(rand.Reader, c.PrivateKey, ciphertext)
	if err != nil {
		panic(err)
	}
	// fmt.Printf("Decrypted: %s\n", plaintext)
	return plaintext
}
