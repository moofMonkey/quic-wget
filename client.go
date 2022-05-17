package main

import (
	"context"
	"crypto/tls"
	"github.com/lucas-clemente/quic-go"
	"io"
	"log"
	"time"
)

func runClient(target, password, remotePath, localPath string, reverse, tcp bool) {
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-wget"},
	}
	dur, _ := time.ParseDuration("10h")
	ctx, cancel := context.WithTimeout(context.Background(), dur)
	defer cancel()
	var conn io.ReadWriteCloser
	var err error
	if !tcp {
		var conf quic.Config
		conf.HandshakeIdleTimeout = dur
		conf.MaxIdleTimeout = dur
		session, err := quic.DialAddrContext(ctx, target, tlsConf, &conf)
		if err != nil {
			log.Fatalln(err)
		}
		if conn, err = session.OpenStreamSync(ctx); err != nil {
			log.Fatalln(err)
		}
	} else {
		if conn, err = tls.Dial("tcp", target, tlsConf); err != nil {
			log.Fatalln(err)
		}
	}
	defer conn.Close()

	log.Println("Session start")

	if err := writeString(conn, password); err != nil {
		log.Fatalln("Failed to write password", err)
	}
	if err := writeString(conn, remotePath); err != nil {
		log.Fatalln("Failed to write remotePath", err)
	}
	transferFile(conn, localPath, !reverse)
}
