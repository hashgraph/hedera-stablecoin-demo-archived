package operation

import (
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Mint(payload *pb.MintTo) error {
	log.Trace().
		Str("username", payload.Address).
		Uint64("quantity", payload.Quantity).
		Msg("Mint")

	// TODO: Handle user not existing
	// TODO: Handle recording the operation
	// TODO: Handle response to the UI
	// TODO: Handle [Balance] not having the user

	// FIXME: UI sends the username where it calls it the address
	// NOTE: I (@mehcode) prefer the username here, but we should change the field name
	state.Balance[payload.Address] += payload.Quantity

	return nil
}
