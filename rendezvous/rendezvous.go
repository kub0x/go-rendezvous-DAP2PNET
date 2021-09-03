package rendezvous

import (
	"dap2pnet/rendezvous/models"
	"math/rand"
	"sync"
	"time"
)

type Rendezvous struct {
	Peers     PeerList
	MaxLinks  int
	MinLinks  int
	listMutex sync.Mutex // for controlling write and iterating on peer list
}

func NewRendezvous() *Rendezvous {
	return &Rendezvous{
		Peers: PeerList{
			List: make(map[string]*models.Triplet),
		},
		MaxLinks: 20,
		MinLinks: 5,
	}
}

func (ren *Rendezvous) AddTriplet(ID string, IP string, port string) {
	ren.listMutex.Lock()
	defer ren.listMutex.Unlock()

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
	// TODO danger here as locks Add and Exchange primitives
	ren.listMutex.Lock()
	defer ren.listMutex.Unlock()

	for _, triplet := range ren.Peers.List {
		if triplet.Expiration > time.Now().UnixNano() {
			delete(ren.Peers.List, triplet.ID)
		}
	}
}

func (ren *Rendezvous) doWholePeerList(ID string) *models.PeerInfo {
	restPeerInfo := &models.PeerInfo{}
	for k, v := range ren.Peers.List {
		if k == ID { // exclude requester node from the list
			continue
		}
		restPeerInfo.Triplets = append(restPeerInfo.Triplets, *v)
	}

	return restPeerInfo
}

func (ren *Rendezvous) doRandomPeerList(ID string) *models.PeerInfo {
	restPeerInfo := &models.PeerInfo{}
	keys := make([]string, 0, len(ren.Peers.List))
	for k := range ren.Peers.List {
		keys = append(keys, k)
	}
	rands := make(map[int]int, ren.MaxLinks)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < ren.MaxLinks; i++ {
		rnd := rand.Intn(len(ren.Peers.List))
		if keys[rnd] == ID { // exclude requester node from the list
			i--
			continue
		} else if rands[rnd] != rnd {
			rands[rnd] = rnd
			restPeerInfo.Triplets = append(restPeerInfo.Triplets, *ren.Peers.List[keys[rnd]])
		} else if rands[rnd] == rnd {
			i--
			continue
		}
	}

	return restPeerInfo
}

func (ren *Rendezvous) IsPeerSubscribed(id string) bool {
	ren.listMutex.Lock()
	defer ren.listMutex.Unlock()

	ret := false
	if ren.Peers.List[id] != nil {
		ret = true
	}
	return ret
}

func (ren *Rendezvous) MakePeerExchangeList(ID string) *models.PeerInfo {
	ren.listMutex.Lock()
	defer ren.listMutex.Unlock()

	if len(ren.Peers.List) <= ren.MinLinks {
		return nil
	}

	var restPeerInfo *models.PeerInfo
	if len(ren.Peers.List) < 2*ren.MaxLinks { // last probability of choice is 1/2 as it has n+1/2n ~ 1/2
		restPeerInfo = ren.doWholePeerList(ID)
	} else {
		restPeerInfo = ren.doRandomPeerList(ID)
	}

	return restPeerInfo
}
