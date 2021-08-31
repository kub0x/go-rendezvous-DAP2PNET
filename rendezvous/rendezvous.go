package rendezvous

import (
	"dap2pnet/rendezvous/models"
	"time"
)

type Rendezvous struct {
	Peers    PeerList
	MaxLinks int
}

func NewRendezvous() *Rendezvous {
	return &Rendezvous{
		Peers:    PeerList{},
		MaxLinks: 5,
	}
}

func (ren *Rendezvous) AddTriplet(ID string, IP string, port string) {
	ren.Peers.Add(
		&models.Triplet{
			ID:         ID,
			IP:         IP,
			Port:       port,
			Expiration: time.Now().Add(time.Minute * 2).UnixNano(),
		},
	)
}

func (ren *Rendezvous) ClearPeerList() { // delete all triplets that exceeded expiration time
	for _, triplet := range ren.Peers.List {
		if triplet.Expiration > time.Now().UnixNano() {
			delete(ren.Peers.List, triplet.ID)
		}
	}
}

func (ren *Rendezvous) MakePeerExchangeList(ID string) *models.PeerInfo {
	restPeerInfo := &models.PeerInfo{}
	ctr := 0
	for _, triplet := range ren.Peers.List {
		if triplet.ID == ID { // exclude requester node from the list
			continue
		} else {
			restPeerInfo.Triplets = append(restPeerInfo.Triplets, *triplet)
		}
		if ctr > ren.MaxLinks {
			break
		} else {
			ctr++
		}
	}
	return restPeerInfo
}
