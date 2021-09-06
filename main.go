package main

import (
	"dap2pnet/rendezvous/rendezvous"
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

	// st, err := storage.Open()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = st.CreateTriplet(models.Triplet{
	// 	ID:         "2231214124524124",
	// 	IP:         "192.168.1.39",
	// 	Port:       "6000",
	// 	Expiration: time.Now().Add(time.Second * 5).Unix(),
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = st.CreateTriplet(models.Triplet{
	// 	ID:         "2414124524124",
	// 	IP:         "192.168.1.39",
	// 	Port:       "6000",
	// 	Expiration: time.Now().Add(time.Second * 5).Unix(),
	// })

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// trs, err := st.GetTriplets()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(trs)
	// fmt.Println(st.GetTripletCount())

	ren := rendezvous.NewRendezvous()
	err := server.Run(ren)
	if err != nil {
		log.Fatal("cannot initialise rendezvous server: " + err.Error())
	}

	// GIN SERVER IS INTERNAL (LOCAL) ONLY SEEN BY THE GATEWAY (gin doesn't allow mutual TLS yet)

}
