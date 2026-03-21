package certs

import (
	"path/filepath"
)

type Certs struct {
	Cert   string
	Key    string
	CertCA string
}

func NewCerts(filename string) (*Certs, error) {
	cert, err := filepath.Abs(filepath.Join("certs", filename+".crt"))

	if err != nil {
		return nil, err
	}

	key, err := filepath.Abs(filepath.Join("certs", filename+".key"))

	if err != nil {
		return nil, err
	}

	certCA, err := filepath.Abs(filepath.Join("certs", "ca.crt"))

	if err != nil {
		return nil, err
	}

	return &Certs{Cert: cert, Key: key, CertCA: certCA}, nil
}
