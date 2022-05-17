package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"log"
	"math/big"
	"time"
)

func handleConnection(sess quic.Connection, password string, reverse bool) {
	// defer sess.Close()
	dur, _ := time.ParseDuration("10h")
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	stream, err := sess.AcceptStream(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	// defer stream.Close()
	clientPassword, err := readString(stream)
	if err != nil {
		log.Println(err)
		return
	}
	if clientPassword != password {
		log.Println("Incorrect password specified:", clientPassword)
		return
	}
	localPath, err := readString(stream)
	if err != nil {
		log.Println(err)
		return
	}

	transferFile(stream, localPath, reverse)
}

// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-wget"},
	}
}

func runServer(target, password string, reverse bool) {
	var conf quic.Config
	dur, _ := time.ParseDuration("10h")
	conf.HandshakeIdleTimeout = dur
	conf.MaxIdleTimeout = dur
	listener, err := quic.ListenAddr(target, generateTLSConfig(), &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		ctx, cancel := context.WithTimeout(context.Background(), dur)
		sess, err := listener.Accept(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		cancel()
		go handleConnection(sess, password, reverse)
	}
}
