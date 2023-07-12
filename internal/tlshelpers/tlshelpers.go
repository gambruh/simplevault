package tlshelpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

const privatekeyfile = "privatekey.pem"
const certfile = "sert.pem"

// ParseRsaPublicKeyFromPem parses a PEM-encoded RSA public key
func ParseRsaPublicKeyFromPem(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the public key")
	}
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DER-encoded public key: %v", err)
	}
	return publicKey, nil
}

// CreateSelfSignedCertificate creates a certificate out of a server's private key file
func CreateSelfSignedCertificate(privatekeyPath string) (certFilepath string, err error) {
	certFilepath = certfile
	privateKeyFile, err := os.ReadFile(privatekeyPath)
	if err != nil {
		log.Fatalln("can't open file with private key:", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyFile)
	if privateKeyBlock == nil || privateKeyBlock.Type != "RSA PRIVATE KEY" {
		return "", fmt.Errorf("failed to decode PEM block containing private key in helpers.CreateSelfSignedCertificate")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key in helpers.CreateSelfSignedCertificate: %w", err)
	}

	// Create a self-signed certificate
	template := x509.Certificate{
		SerialNumber:          big.NewInt(time.Now().Unix()),
		Subject:               pkix.Name{CommonName: "practicum.yandex.ru"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", err
	}

	// Write the certificate to a PEM file
	certOut, err := os.Create(certFilepath)
	if err != nil {
		return "", err
	}
	defer certOut.Close()

	certBlock := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	}
	if err := pem.Encode(certOut, &certBlock); err != nil {
		panic(err)
	}
	return certFilepath, nil
}

// CreateConfigTLS creates TLS config for a server out of tls certificate and private key file
func CreateConfigTLS(certFile string, privatekeyFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, privatekeyFile)
	if err != nil {
		return nil, fmt.Errorf("error when loading tls certificate and key: %s", err)
	}

	tlsconfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			fmt.Printf("Received TLS handshake request from client %s:%d\n", info.Conn.RemoteAddr().String(), info.Conn.RemoteAddr().(*net.TCPAddr).Port)
			cert, err := tls.LoadX509KeyPair(certFile, privatekeyfile)
			if err != nil {
				return nil, fmt.Errorf("error when loading KeyPair in CreateConfigTLS: %s", err)
			}
			return &cert, nil
		},
		Certificates: []tls.Certificate{cert},
	}
	return tlsconfig, nil
}
