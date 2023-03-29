package main

import (
	"log"
	"net"
	"net/http"
	"server/chatgpt"
	"server/mrpc"
	"server/registry"
	"sync"
	"time"
)

func startRegistry(wg *sync.WaitGroup) {
	l, _ := net.Listen("tcp", ":9999")
	registry.HandleHTTP()
	wg.Done()
	_ = http.Serve(l, nil)
}

func startServer(registryAddr string, wg *sync.WaitGroup) {
	var c chatgpt.OpenAIClient
	l, _ := net.Listen("tcp", ":0")
	server := mrpc.NewServer()
	_ = server.Register(&c)
	registry.Heartbeat(registryAddr, "tcp@"+l.Addr().String(), 0)
	wg.Done()
	server.Accept(l)
}

func main() {
	log.SetFlags(0)
	registryAddr := "http://localhost:9999/_mrpc_/registry"
	var wg sync.WaitGroup
	wg.Add(1)
	go startRegistry(&wg)
	wg.Wait()

	time.Sleep(time.Second)
	wg.Add(1)
	go startServer(registryAddr, &wg)
	// go startServer(registryAddr, &wg)
	wg.Wait()

	c := make(chan bool)
	<-c
}
