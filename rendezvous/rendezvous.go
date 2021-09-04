package rendezvous

import (
	"dap2pnet/rendezvous/models"
	"dap2pnet/rendezvous/storage"
	"errors"
	"log"
	"math/rand"
	"time"
)

var (
	RendezvousErrPeerUnsuscribed = errors.New("you are not subscribed to the rendezvous")
	RendezvousErrPeerExpired     = errors.New("your session has expired")
	RendezvousErrPeerMinlinks    = errors.New("not enough peers in the list")
)

type Rendezvous struct {
	MaxLinks int
	MinLinks int
	st       *storage.Storage
}

func NewRendezvous() *Rendezvous {
	st, err := storage.Open()
	if err != nil {
		log.Fatal(err)
	}

	return &Rendezvous{
		MaxLinks: 20,
		MinLinks: 5,
		st:       st,
	}
}

func (ren *Rendezvous) AddTriplet(ID string, IP string, port string) error {
	return ren.st.CreateTriplet(models.Triplet{
		ID:         ID,
		IP:         IP,
		Port:       port,
		Expiration: time.Now().Add(time.Minute * 2).Unix(),
	})
}

func (ren *Rendezvous) doPeerList(ID string) (*models.PeerInfo, error) {
	restPeerInfo := &models.PeerInfo{}
	trs, err := ren.st.GetTriplets()
	if err != nil {
		return nil, err
	}

	limit := ren.MaxLinks
	if len(trs) < 2*ren.MaxLinks {
		limit = len(trs)
	}

	perm := rand.Perm(limit) // pseudo random permutation ftw
	from := 0
	to := perm[from]
	for i := 0; i < limit; i++ {
		v := perm[from]
		from = to
		if trs[v].ID == ID {
			continue // exclude requester node from the list
		}
		restPeerInfo.Triplets = append(restPeerInfo.Triplets, trs[v])
	}

	return restPeerInfo, nil
}

func (ren *Rendezvous) IsPeerSubscribed(id string) error {
	triplet, err := ren.st.GetTriplet(id)
	if err != nil {
		err = RendezvousErrPeerUnsuscribed
	}

	now := time.Now().Unix()
	if now > triplet.Expiration { // peer has ttl lease on etcd but who knows...
		err = RendezvousErrPeerExpired
	}

	return err
}

func (ren *Rendezvous) MakePeerExchangeList(ID string) (*models.PeerInfo, error) {
	count, err := ren.st.GetTripletCount()
	if err != nil {
		return nil, err
	}

	if int(count) <= ren.MinLinks {
		return nil, RendezvousErrPeerMinlinks
	}

	restPeerInfo, err := ren.doPeerList(ID)
	if err != nil {
		return nil, err
	}

	if len(restPeerInfo.Triplets) <= ren.MinLinks {
		err = RendezvousErrPeerMinlinks
	}

	return restPeerInfo, err
}
