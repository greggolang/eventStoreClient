package main

import (
	"flag"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	dummy_server "github.com/greggolang/eventStoreClient/dummy/server"
)

func main() {
	log.SetOutput(os.Stdout)
	rand.Seed(time.Now().UnixNano())
	port := flag.Int("port", 8080, "port to run eventeur server on")
	flag.Parse()

	server := dummy_server.NewDummyServer(dummy_server.Config{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := server.Serve(*port)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// TODO: gracefull shutdown
	wg.Wait()
}
