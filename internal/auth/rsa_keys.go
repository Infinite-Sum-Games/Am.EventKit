package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"os"
)

// Generates RSA keys only if they don't exist
func InitRSAKeys(privateKeyPath, publicKeyPath string) error {
	if fileExists(privateKeyPath) && fileExists(publicKeyPath) {
		log.Println("RSA keys already exist. Skipping generation.")
		return nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privBytes,
	})
	if err := ioutil.WriteFile(privateKeyPath, privPem, 0600); err != nil {
		return err
	}

	pubBytes := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubBytes,
	})
	if err := ioutil.WriteFile(publicKeyPath, pubPem, 0644); err != nil {
		return err
	}

	log.Println("RSA key pair generated successfully.")
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// Verifies that public key matches the private key
func VerifyRSAKeys(privateKeyPath, publicKeyPath string) error {
	privData, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return err
	}
	pubData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return err
	}

	privBlock, _ := pem.Decode(privData)
	pubBlock, _ := pem.Decode(pubData)

	if privBlock == nil || pubBlock == nil {
		return errors.New("invalid PEM blocks")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return err
	}
	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return err
	}

	if privKey.PublicKey.N.Cmp(pubKey.N) != 0 {
		return errors.New("public key does not match private key")
	}

	return nil
}

// Loads keys from files into memory
func LoadRSAKeysFromFiles(privateKeyPath, publicKeyPath string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, nil, err
	}

	privBlock, _ := pem.Decode(privBytes)
	if privBlock == nil || privBlock.Type != "RSA PRIVATE KEY" {
		return nil, nil, errors.New("failed to decode PEM block containing private key")
	}

	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	pubBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, nil, err
	}

	pubBlock, _ := pem.Decode(pubBytes)
	if pubBlock == nil || pubBlock.Type != "RSA PUBLIC KEY" {
		return nil, nil, errors.New("failed to decode PEM block containing public key")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(pubBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	return privKey, pubKey, nil
}
