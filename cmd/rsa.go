package cmd

import (
	"fmt"
	"os"

	"aidanwoods.dev/go-paseto"
)

func CheckRSAKeyPairExists() error {
	privateKeyPath := "app.rsa"
	publicKeyPath := "app.pub.rsa"
	if _, err := os.Stat(privateKeyPath); err != nil {
		return err
	}
	if _, err := os.Stat(publicKeyPath); err != nil {
		return err
	}
	return nil
}

func GenerateRSAKeyPair() error {
	privateKey := paseto.NewV4AsymmetricSecretKey()
	publicKey := privateKey.Public()

	// Writing both keys to individual files
	if err := os.WriteFile("app.rsa", privateKey.ExportBytes(), 0644); err != nil {
		return fmt.Errorf("error saving private key: %w", err)
	}
	if err := os.WriteFile("app.pub.rsa", publicKey.ExportBytes(), 0644); err != nil {
		return fmt.Errorf("error saving public key: %w", err)
	}
	return nil
}
