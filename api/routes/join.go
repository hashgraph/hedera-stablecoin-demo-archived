package routes

import (
	"github.com/gofiber/fiber"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/pb"
	"net/http"
)

func SendAnnounce(c *fiber.Ctx) {
	var req struct {
		PublicKey string `json:"publicKey"`
		Username  string `json:"username"`
	}

	err := c.BodyParser(&req)
	if err != nil {
		log.Warn().Msgf("%v", err)
		c.SendStatus(400)
		return
	}

	v := &pb.Join{
		Address:  req.PublicKey,
		Username: req.Username,
	}

	sendTransaction(v, &pb.Primitive{Primitive: &pb.Primitive_Join{Join: v}})

	c.Status(http.StatusAccepted)
	err = c.JSON(transactionResponse{
		Status:  true,
		Message: "Join request sent",
	})

	if err != nil {
		log.Error().Msgf("%v", err)
		c.SendStatus(500)
		return
	}
}
