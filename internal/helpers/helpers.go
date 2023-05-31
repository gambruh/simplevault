package helpers

import (
	"crypto/tls"
	"fmt"
	"net"

	"github.com/gambruh/gophkeeper/internal/config"
)

const privatekeyfile = "privatekey.pem"

func CreateConfigTLS() (*tls.Config, error) {

	cert, err := tls.LoadX509KeyPair(config.Cfg.Certificate, privatekeyfile)
	if err != nil {
		return nil, fmt.Errorf("error when loading tls certificate and key: %s", err)
	}

	tlsconfig := &tls.Config{
		GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
			fmt.Printf("Received TLS handshake request from client %s:%d\n", info.Conn.RemoteAddr().String(), info.Conn.RemoteAddr().(*net.TCPAddr).Port)
			cert, err := tls.LoadX509KeyPair(config.Cfg.Certificate, privatekeyfile)
			if err != nil {
				return nil, fmt.Errorf("error when loading KeyPair in CreateConfigTLS: %s", err)
			}
			return &cert, nil
		},
		Certificates: []tls.Certificate{cert},
	}
	return tlsconfig, nil
}
