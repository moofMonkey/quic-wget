package main

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"os"
	"time"
)

func runClient(target, password, downloadPath, localPath string) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-wget"},
	}
	dur, _ := time.ParseDuration("10h")
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	var conf quic.Config
	conf.HandshakeIdleTimeout = dur
	conf.MaxIdleTimeout = dur
	session, err := quic.DialAddrContext(ctx, target, tlsConf, &conf)
	if err != nil {
		log.Fatalln(err)
	}

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer stream.Close()

	log.Println("Session start")

	if err := writeString(stream, password); err != nil {
		log.Fatalln("Failed to write password", err)
	}
	if err := writeString(stream, downloadPath); err != nil {
		log.Fatalln("Failed to write downloadPath", err)
	}
	size, err := readUint64(stream)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Downloading:", size)

	f, err := os.Create(localPath)
	if err != nil {
		log.Fatalln("Failed to create file", localPath, err)
	}
	defer f.Close()

	if _, err = io.CopyN(f, stream, int64(size)); err != nil {
		log.Fatalln("Failed to transfer file", err)
	}
}
