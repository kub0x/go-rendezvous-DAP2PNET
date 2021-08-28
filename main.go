package main

import (
	"dap2pnet/rendezvous/server"
	"log"
)

func main() {

	//RUN TLS SERVER AS AN EXPOSED GATEWAY

	/*tlsProxy := TLSProxy.NewTLSProxy()
	go func() {
		err := tlsProxy.Listen()
		if err != nil {
			log.Fatal("cannot initialise TLS Proxy: " + err.Error())
		}
	}()*/

	err := server.Run()
	if err != nil {
		log.Fatal("cannot initialise rendezvous server: " + err.Error())
	}

	// GIN SERVER IS INTERNAL (LOCAL) ONLY SEEN BY THE GATEWAY (gin doesn't allow mutual TLS yet)

}
