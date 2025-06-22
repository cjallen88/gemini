package server

import (
	"crypto/sha256"
	"crypto/tls"
	"fmt"
	"net"
)

type Certificate struct {
	Name        string
	Fingerprint [32]byte
}

func ClientCert(conn net.Conn) (*Certificate, error) {
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return nil, fmt.Errorf("Connection is not a TLS connection")
	}

	if err := tlsConn.Handshake(); err != nil {
		return nil, fmt.Errorf("TLS handshake failed: %w", err)
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		return nil, nil
	}

	cert := state.PeerCertificates[0]
	fingerprint := sha256.Sum256(cert.Raw)
	return &Certificate{
		Name:        cert.Subject.CommonName,
		Fingerprint: fingerprint,
	}, nil
}
