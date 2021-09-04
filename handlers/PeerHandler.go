package handlers

import (
	"dap2pnet/rendezvous/models"
	"dap2pnet/rendezvous/rendezvous"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	PeerHandlerErrUnvalidPort  = errors.New("you must select a valid port")
	PeerHandlerErrUnauthorized = errors.New("you must subscribe to rendezvous in order to access peer lists")
	PeerHandlerErrMinLinks     = errors.New("not enough peers to link")
)

func OnSubscribe(ren *rendezvous.Rendezvous) gin.HandlerFunc {
	return func(c *gin.Context) {
		var subReq models.SubscribeRequest
		err := c.BindJSON(&subReq)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, PeerHandlerErrUnvalidPort)
			return
		}

		id := c.GetString("Identity")
		err = ren.AddTriplet(id, c.GetHeader("X-Forwarded-For"), fmt.Sprint(subReq.Port))
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.Status(http.StatusOK)
	}
}

func OnGetPeers(ren *rendezvous.Rendezvous) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("Identity")
		err := ren.IsPeerSubscribed(id)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		peerList, err := ren.MakePeerExchangeList(id)
		if peerList == nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		c.JSON(http.StatusOK, peerList)
	}
}
