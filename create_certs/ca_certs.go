package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"time"
)

func GenerateCACertTemplate()(*x509.Certificate){
	return &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}
}

func GenerateKey(bits int)(*rsa.PrivateKey, error){
	caPrivKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return caPrivKey, nil
}

func CreateSelfSignedCACert(ca *x509.Certificate, caPrivateKey *rsa.PrivateKey)([]byte, error){
	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivateKey.PublicKey, caPrivateKey)
	if err != nil {
		return nil, err
	}
	return caBytes, nil
}

func EncodeCertPEMEncoded(certBytes []byte, compPrivateKey *rsa.PrivateKey)([]byte, []byte){
	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(compPrivateKey),
	})

	return caPEM.Bytes(), caPrivKeyPEM.Bytes()

}

func GenerateComponentCertsTemplate()*x509.Certificate{
	return &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}
}

func SignCertWithCA(compCert *x509.Certificate, caCert *x509.Certificate, compCertPrivateKey *rsa.PrivateKey)([]byte, error){
	certBytes, err := x509.CreateCertificate(rand.Reader, compCert, caCert, &compCertPrivateKey.PublicKey, compCertPrivateKey)
	if err != nil {
		return nil, err
	}

	return certBytes, nil
}

func main() {
	caCert := GenerateCACertTemplate()
	caKey, err := GenerateKey(2048)
	if err != nil{
		log.Fatalln(err)
	}

	rootCACert, err := CreateSelfSignedCACert(caCert, caKey)
	if err != nil{
		log.Fatalln(err)
	}

	pemKey, privateKey  := EncodeCertPEMEncoded(rootCACert, caKey)
	fmt.Println("PEM key",string(pemKey))
	fmt.Println("Private Key", string(privateKey))
}


