package routes

import (
	"github.com/gofiber/fiber"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
	"strconv"
)

func SendMint(c *fiber.Ctx) {
	var req struct {
		// FIXME: UI sends the username where it calls it the address
		// NOTE: I (@mehcode) prefer the username here, but we should change the field name
		Username  string `json:"address"`
		Quantity string `json:"quantity"`
	}

	err := c.BodyParser(&req)
	if err != nil {
		return
	}

	quantity, err := strconv.Atoi(req.Quantity)
	if err != nil {
		log.Warn().Msgf("%v", err)
		c.SendStatus(400)
		return
	}

	v := &pb.MintTo{
		Address:  req.Username,
		Quantity: uint64(quantity),
	}

	sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_MintTo{MintTo: v}})

	c.Status(http.StatusAccepted)
	err = c.JSON(transactionResponse{
		Status:  true,
		Message: "MintTo request sent",
	})

	if err != nil {
		log.Error().Msgf("%v", err)
		c.SendStatus(500)
		return
	}
}
