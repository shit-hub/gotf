package main

import (
	"flag"
	"log"
	"runtime"
	"sync"

	"github.com/vbauerster/mpb/v8"
)

var host = flag.String("l", "0.0.0.0:3000", "Server Host")
var path = flag.String("f", ".", "Upload Path")
var aes = flag.String("aes", "", "enable AES encrypt and set mode: CBC/ECB/CFB")
var aesKey = flag.String("aes-key", "ABCDEFGHIJKLMNOP", "the key of AES encrypt")

func main() {
	flag.Parse()
	if len(*aesKey) != 16 && len(*aesKey) != 24 && len(*aesKey) != 32 {
		log.Fatalln("Unvalid AES Key, the key length need to be 16, 24 or 32")
	}
	log.Println("Server Host: ", *host)
	log.Println("Upload Path: ", *path)

	files, err := getAllFile(*path)
	if err != nil {
		log.Fatalln(err)
	}

	concurrent := make(chan int, runtime.NumCPU())
	defer close(concurrent)

	var wg sync.WaitGroup
	p := mpb.New(mpb.WithWaitGroup(&wg))
	for _, file := range files {
		wg.Add(1)
		concurrent <- 1
		go sendFile(concurrent, &wg, p, *host, *path, file)
	}

	wg.Wait()

	log.Println("Completed !!")
}
