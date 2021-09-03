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
		ren.AddTriplet(id, c.GetHeader("X-Forwarded-For"), fmt.Sprint(subReq.Port))

		c.Status(http.StatusOK)
	}
}

func OnGetPeers(ren *rendezvous.Rendezvous) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetString("Identity")
		if !ren.IsPeerSubscribed(id) {
			c.AbortWithError(http.StatusUnauthorized, PeerHandlerErrUnauthorized)
			return
		}

		peerList := ren.MakePeerExchangeList(id)
		if peerList == nil {
			c.AbortWithError(http.StatusUnauthorized, PeerHandlerErrMinLinks)
			return
		}

		c.JSON(http.StatusOK, peerList)
	}
}
