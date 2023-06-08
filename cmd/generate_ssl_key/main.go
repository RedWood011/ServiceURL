package main

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"log"
	"math/big"
	"net"
	"os"
	"time"
)

func generate() {
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Evgeniy Yaroshenko"},
			Country:      []string{"RU"},
			Locality:     []string{"Moscow"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 1 год
		NotAfter:     time.Now().AddDate(1, 0, 0),
		SubjectKeyId: []byte{0, 8, 9, 4, 7, 6, 6, 3, 4, 1, 4},
		// устанавливаем использование ключа для цифровой подписи, а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит

	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		log.Fatal(err)
	}

	// создаём x.509-сертификат
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// кодируем сертификат и ключ в формате PEM, который используется для хранения и обмена криптографическими ключами
	var certPEM bytes.Buffer
	err = pem.Encode(&certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	if err != nil {
		log.Fatal(err)
	}

	fileserver, err := os.Create("server_crt.crt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = fileserver.Write(certPEM.Bytes())

	if err != nil {
		log.Fatal(err)
	}

	var privateKeyPEM bytes.Buffer
	err = pem.Encode(&privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if err != nil {
		log.Fatal(err)
	}

	filekey, err := os.Create("server_key.key")
	_, err = filekey.Write(privateKeyPEM.Bytes())

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	generate()
}
