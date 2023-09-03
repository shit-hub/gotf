package main

import (
	"flag"
	"log"
	"net"
	"sync"

	"github.com/vbauerster/mpb/v8"
)

var host = flag.String("l", "0.0.0.0:3000", "Listen Host")
var path = flag.String("f", ".", "Download Path")
var aes = flag.String("aes", "", "enable AES encrypt and set mode: CBC/ECB/CFB")
var aesKey = flag.String("aes-key", "ABCDEFGHIJKLMNOP", "the key of AES encrypt")

func main() {
	flag.Parse()
	log.Println("Listen Host: ", *host)
	log.Println("Download Path: ", *path)

	ln, err := net.Listen("tcp", *host)
	if err != nil {
		log.Fatalf("Failed to listen host[%s]: %v", *host, err)
	}

	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Accept error: ", err)
			continue
		}

		wg.Add(1)
		go receiveFile(conn, p, *path)
	}
}
