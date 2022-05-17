package main

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"log"
	"time"
)

func runClient(target, password, remotePath, localPath string, reverse bool) {
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
	if err := writeString(stream, remotePath); err != nil {
		log.Fatalln("Failed to write remotePath", err)
	}
	transferFile(stream, localPath, !reverse)
}
