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

func main() {
	flag.Parse()
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
