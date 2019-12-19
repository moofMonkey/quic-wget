package main

import (
	"encoding/binary"
	"errors"
	"io"
	"fmt"
	"io/ioutil"
	"log"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"crypto/rand"
	"crypto/rsa"
	"context"
	"math/big"
	"runtime"
	"time"
	"flag"

	quic "github.com/lucas-clemente/quic-go"
)

func Read(r io.Reader, dst []byte) error {
	nBytesRead, err := r.Read(dst)
	if err != nil {
		return err
	}
	for nBytesRead != len(dst) {
		nBytesRead2, err := r.Read(dst[nBytesRead:])
		if err != nil {
			return err
		}
		nBytesRead += nBytesRead2
	}
	if nBytesRead != len(dst) {
		return errors.New("unable to read bytes")
	}
	return nil
}

func ReadUint32(r io.Reader) (uint32, error) {
	var buf [4]byte
	err := Read(r, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}

// NextByte reads the next byte from the buffer
func ReadByte(r io.Reader) (byte, error) {
	var b [1]byte
	err := Read(r, b[:])
	return b[0], err
}

func ReadBytes(r io.Reader, n uint32) ([]byte, error) {
	buf := make([]byte, n)
	if err := Read(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadString reads a null terminated string
func ReadString(r io.Reader) (string, error) {
	buf := make([]byte, 0)
	for {
		b, err := ReadByte(r)
		if err != nil {
			return string(buf), err
		}
		if b == 0 {
			break
		}
		buf = append(buf, b)
	}

	return string(buf), nil
}

func Write(w io.Writer, src []byte) error {
	nBytesWritten, err := w.Write(src)
	if err != nil {
		return err
	}
	if nBytesWritten != len(src) {
		return errors.New("unable to write bytes")
	}
	return nil
}

func WriteUint32(w io.Writer, x uint32) error {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	return Write(w, buf[:])
}

// WriteByte writes a single byte
func WriteByte(w io.Writer, b byte) error {
	var ar [1]byte
	ar[0] = b
	err := Write(w, ar[:])
	if err != nil {
		return err
	}
	return nil
}

// WriteString writes a null-terminated string
func WriteString(w io.Writer, str string) error {
	if err := Write(w, []byte(str)); err != nil {
		return err
	}
	if err := WriteByte(w, 0); err != nil {
		return err
	}
	return nil
}

var password string
func handleConnection(sess quic.Session) {
	// defer sess.Close()
	dur, _ := time.ParseDuration("10h")
	ctx, _ := context.WithTimeout(context.Background(), dur)
	stream, err := sess.AcceptStream(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	// defer stream.Close()
	clientPassword, err := ReadString(stream)
	if err != nil {
		log.Println(err)
		return
	}
	if clientPassword != password {
		log.Println("Incorrect password specified:", clientPassword)
		return
	}
	path, err := ReadString(stream)
	if err != nil {
		log.Println(err)
		return
	}
	
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err.Error())
		return
	}

	err = WriteUint32(stream, uint32(len(buf)))
	if err != nil {
		log.Println(err.Error())
		return
	}
	err = Write(stream, buf)
	if err != nil {
		log.Println(err.Error())
		return
	}
}

func main() {
	var target string
	flag.StringVar(&password, "password", "", "Password for this server")
	flag.StringVar(&target, "target", "0.0.0.0:58993", "Target for this server to listen")
	flag.Parse()

	fmt.Printf("GOMAXPROCS is %d\n", runtime.GOMAXPROCS(0))

	var conf quic.Config
	dur, _ := time.ParseDuration("10h")
	conf.HandshakeTimeout = dur
	conf.MaxIdleTimeout = dur
	listener, err := quic.ListenAddr(target, generateTLSConfig(), &conf)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	for {
		ctx, _ := context.WithTimeout(context.Background(), dur)
		sess, err := listener.Accept(ctx)
		if err != nil {
			log.Fatal(err)
			return
		}
		go handleConnection(sess)
	}
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
