package main

import (
	"dap2pnet/rendezvous/TLSProxy"
	"log"
)

func main() {

	//RUN TLS SERVER AS AN EXPOSED GATEWAY

	tlsProxy := TLSProxy.NewTLSProxy()

	err := tlsProxy.Listen()
	if err != nil {
		log.Fatal("cannot initialise TLS Proxy: " + err.Error())
	}

	// GIN SERVER IS INTERNAL (LOCAL) ONLY SEEN BY THE GATEWAY (gin doesn't allow mutual TLS yet)

}
