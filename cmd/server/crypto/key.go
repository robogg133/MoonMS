package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GenerateKeyPair(bitsAmount int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitsAmount)
	if err != nil {
		return err
	}

	privateKeyMarshalized := x509.MarshalPKCS1PrivateKey(privateKey)

	header := make(map[string]string)
	header["BIT AMOUNT"] = fmt.Sprint(bitsAmount)

	privateKeyBlock := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: header,
		Bytes:   privateKeyMarshalized,
	}
	file, err := os.OpenFile("server-private-key.pem", os.O_CREATE|os.O_WRONLY, 0400)
	if err != nil {
		file.Close()
		return err
	}
	if err := pem.Encode(file, privateKeyBlock); err != nil {
		return err
	}
	file.Close()

	publicKey := privateKey.Public()

	publicKeyMarshalized, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	return os.WriteFile("server-public.key", publicKeyMarshalized, 0644)
}
