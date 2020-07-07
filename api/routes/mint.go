package routes

import (
	"github.com/labstack/echo/v4"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
	"strconv"
)

func SendMint(c echo.Context) error {
	var req struct {
		// FIXME: UI sends the username where it calls it the address
		// NOTE: I (@mehcode) prefer the username here, but we should change the field name
		Username  string `json:"address"`
		Quantity string `json:"quantity"`
	}

	err := c.Bind(&req)
	if err == nil {
		quantity, err := strconv.Atoi(req.Quantity)
		if err == nil {

			v := &pb.MintTo{
				Address:  req.Username,
				Quantity: uint64(quantity),
			}
			err = sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_MintTo{MintTo: v}})
		}
	}

	if (err == nil) {
		return c.JSON(http.StatusAccepted, transactionResponse{
			Status:  true,
			Message: "MintTo request sent",
		})
	} else {
		return c.JSON(http.StatusInternalServerError, transactionResponse{
			Status:  false,
			Message: err.Error(),
		})
	}
}
