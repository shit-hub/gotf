package main

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const bufSize = 2048

func createDir(filename string) {
	index := strings.LastIndex(filename, "/")
	if index != 0 {
		err := os.MkdirAll(filename[:index], 0o755)
		if err != nil {
			log.Println("Failed to create dir: ", err)
		}
	}
}

func receiveFile(conn net.Conn, p *mpb.Progress, path string) {
	defer conn.Close()

	// Get file name length
	var fnLen uint32
	fnLenBuf := make([]byte, 4)
	_, err := conn.Read(fnLenBuf)
	if err != nil {
		log.Println("Failed to read file name length: ", err)
		return
	}
	binary.Read(bytes.NewBuffer(fnLenBuf), binary.BigEndian, &fnLen)

	// Get file name
	filename := make([]byte, fnLen)
	_, err = conn.Read(filename)
	if err != nil {
		log.Println("Failed to read file name: ", err)
		return
	}

	// Get file size
	var fileSize int64
	fileSizeBuf := make([]byte, 8)
	_, err = conn.Read(fileSizeBuf)
	if err != nil {
		log.Println("Failed to read file size: ", err)
		return
	}
	binary.Read(bytes.NewBuffer(fileSizeBuf), binary.BigEndian, &fileSize)

	// Receive file
	fnStr := path + "/" + string(filename)
	createDir(fnStr)
	f, err := os.OpenFile(fnStr, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0o666)
	if err != nil {
		log.Printf("Failed to open file[%s]: %v", fnStr, err)
		return
	}
	defer f.Close()

	var size int64
	buf := make([]byte, bufSize)
	bar := p.AddBar(int64(fileSize),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(string(filename), decor.WC{W: len(filename) + 1, C: decor.DidentRight}),
			// decor.DSyncWidth bit enables column width synchronization
			decor.Percentage(decor.WCSyncSpace),
		),
		mpb.AppendDecorators(
			// replace ETA decorator with "done" message, OnComplete event
			decor.OnComplete(
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}), "done",
			),
		),
	)
	for {
		n, err := conn.Read(buf)

		if n > 0 {
			if err != nil {
				log.Printf("Failed to get file[%s] stat: %v", fnStr, err)
				break
			}
			f.WriteAt(buf[0:n], size)
			size += int64(n)
			bar.IncrInt64(int64(n))
		}

		if err != nil {
			if err != io.EOF {
				log.Printf("Failed to read from connection: %v", err)
			}
			break
		}
	}
}
