package main

import (
	"context"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"errors"
	"encoding/binary"
	"fmt"
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

func main() {
	var target string
	var password string
	var downloadPath string
	var localPath string
	flag.StringVar(&target, "target", "127.0.0.1:58993", "Target IP:Port")
	flag.StringVar(&password, "password", "", "Password for specified target")
	flag.StringVar(&downloadPath, "downloadPath", "", "Path to be downloaded")
	flag.StringVar(&localPath, "localPath", "", "Path where file will be stored")
	flag.Parse()
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-wget"},
	}
	dur, _ := time.ParseDuration("10h")
	ctx, _ := context.WithTimeout(context.Background(), dur)
	var conf quic.Config
	conf.HandshakeTimeout = dur
	conf.MaxIdleTimeout = dur
	session, err := quic.DialAddrContext(ctx, target, tlsConf, &conf)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	stream, err := session.OpenStreamSync(ctx)
	if err != nil {
		panic(err)
	}
	defer stream.Close()
	
	fmt.Printf("Session start\n")
	
	WriteString(stream, password)
	WriteString(stream, downloadPath)
	size, err := ReadUint32(stream)
	if err != nil {
		panic(err)
	}
	
	log.Println("Downloading: ", size)
	
	bytes, err := ReadBytes(stream, size)
	ioutil.WriteFile(localPath, bytes, 0644)
	if err != nil {
		panic(err)
	}
}
