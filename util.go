package main

import (
	"encoding/binary"
	"io"
	"log"
	"os"
)

func readUint8(r io.Reader) (uint8, error) {
	var buf [1]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func readUint64(r io.Reader) (uint64, error) {
	var buf [8]byte
	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

func readString(r io.Reader) (string, error) {
	size, err := readUint8(r)
	if err != nil {
		return "", err
	}
	if size == 0 {
		return "", nil
	}
	buf := make([]byte, size)
	if _, err = io.ReadFull(r, buf); err != nil {
		return "", err
	}
	return string(buf), nil
}

func writeUint64(w io.Writer, x uint64) error {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	_, err := w.Write(buf[:])
	return err
}

func writeUint8(w io.Writer, x uint8) error {
	buf := [1]byte{x}
	_, err := w.Write(buf[:])
	return err
}

func writeString(w io.Writer, s string) error {
	if len(s) > 255 {
		return io.ErrShortWrite
	}
	b := []byte(s)
	if err := writeUint8(w, uint8(len(b))); err != nil {
		return err
	}
	if _, err := w.Write(b); err != nil {
		return err
	}
	return nil
}

func transferFile(conn io.ReadWriter, localPath string, download bool) {
	if download {
		size, err := readUint64(conn)
		if err != nil {
			log.Fatalln(err)
		}

		log.Println("Downloading:", size)

		f, err := os.Create(localPath)
		if err != nil {
			log.Fatalln("Failed to create file", localPath, err)
		}
		defer f.Close()

		if _, err = io.CopyN(f, conn, int64(size)); err != nil {
			log.Fatalln("Failed to transfer file", err)
		}
	} else {
		f, err := os.Open(localPath)
		if err != nil {
			log.Fatalln("Failed to open file", localPath, err)
			return
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			log.Println("Failed to get file stat", err)
			return
		}

		if err = writeUint64(conn, uint64(stat.Size())); err != nil {
			log.Println("Failed to write file size", err)
			return
		}

		if _, err = io.CopyN(conn, f, stat.Size()); err != nil {
			log.Fatalln("Failed to transfer file", err)
		}
	}
}
