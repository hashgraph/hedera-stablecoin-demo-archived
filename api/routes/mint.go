package routes

import (
	"github.com/gin-gonic/gin"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
	"strconv"
)

func SendMint(c *gin.Context) {
	var req struct {
		// FIXME: UI sends the username where it calls it the address
		// NOTE: I (@mehcode) prefer the username here, but we should change the field name
		Username  string `json:"address"`
		Quantity string `json:"quantity"`
	}

	err := c.BindJSON(&req)
	if err != nil {
		return
	}

	quantity, err := strconv.Atoi(req.Quantity)
	if err != nil {
		c.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	v := &pb.MintTo{
		Address:  req.Username,
		Quantity: uint64(quantity),
	}

	go sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_MintTo{MintTo: v}})

	c.JSON(http.StatusAccepted, transactionResponse{
		Status:  true,
		Message: "MintTo request sent",
	})
}
