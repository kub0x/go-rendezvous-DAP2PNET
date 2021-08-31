package rendezvous

import "dap2pnet/rendezvous/models"

type PeerList struct {
	List map[string]*models.Triplet
}

func (pl *PeerList) Add(triplet *models.Triplet) {
	if pl.List[triplet.ID] == nil {
		pl.List[triplet.ID] = triplet
	}
}

// func (pl *PeerList) ToJSON(t) models.Triplet {
// 	return pl.List[[]]
// }
