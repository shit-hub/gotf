package main

import (
	"bytes"
	"encoding/binary"
	"gotf/encrypt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

const bufSize = 2048

func getAllFile(path string) ([]string, error) {
	var s []string
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			s = append(s, fi.Name())
		} else {
			ss, err := getAllFile(path + "/" + fi.Name())
			if err != nil {
				return s, err
			}
			for _, n := range ss {
				s = append(s, fi.Name()+"/"+n)
			}
		}
	}
	return s, nil
}

func sendFile(ch chan int, wg *sync.WaitGroup, p *mpb.Progress, host, path, filename string) {
	defer func() {
		wg.Done()
		<-ch
	}()

	// Connect to the server
	conn, err := net.Dial("tcp", host)
	if err != nil {
		log.Println("Failed to connect to server:", err)
		return
	}
	defer conn.Close()

	// Openfile
	fullName := path + "/" + filename
	f, err := os.OpenFile(fullName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Printf("Failed to open file[%s]: %v", fullName, err)
		return
	}
	defer f.Close()

	// Send filename length
	fnLenBuf := bytes.NewBuffer([]byte{})
	binary.Write(fnLenBuf, binary.BigEndian, uint32(len(filename)))
	_, err = conn.Write(fnLenBuf.Bytes())
	if err != nil {
		log.Printf("Failed to send filename leanth: %v", err)
		return
	}

	// Send filename
	_, err = conn.Write([]byte(filename))
	if err != nil {
		log.Printf("Failed to send filename: %v", err)
		return
	}

	// Send file size
	s, err := f.Stat()
	if err != nil {
		log.Printf("Failed to get file stat: %v", err)
		return
	}
	fileSize := s.Size()
	fileSizeBuf := bytes.NewBuffer([]byte{})
	binary.Write(fileSizeBuf, binary.BigEndian, fileSize)
	_, err = conn.Write(fileSizeBuf.Bytes())
	if err != nil {
		log.Printf("Failed to send file size: %v", err)
	}

	// Send file
	buf := make([]byte, bufSize)
	bar := p.AddBar(int64(fileSize),
		mpb.PrependDecorators(
			// display our name with one space on the right
			decor.Name(filename, decor.WC{W: len(filename) + 1, C: decor.DidentRight}),
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
		n, readErr := f.Read(buf)

		if n > 0 {
			// Encrypt buffer
			body := encrypt.AesEncrypt(buf[0:n], []byte(*aesKey), *aes)
			// Send buffer length
			bufSizeBuf := bytes.NewBuffer([]byte{})
			err = binary.Write(bufSizeBuf, binary.BigEndian, uint32(len(body)))
			if err != nil {
				log.Printf("Failed to write buffer size to buffer: ", err)
				break
			}
			_, err = conn.Write(bufSizeBuf.Bytes())
			if err != nil {
				log.Printf("Failed to send buffer size: %v", err)
				break
			}

			//Send buffer
			_, err := conn.Write(body)
			if err != nil {
				log.Printf("Failed to send buffer[%s]: %v", filename, err)
				break
			}

			// Update bar
			bar.IncrInt64(int64(n))
		}

		// Send Complete
		if readErr != nil {
			if readErr != io.EOF {
				log.Printf("Failed to read file[%s]: %v", filename, readErr)
			}
			break
		}
	}
}
