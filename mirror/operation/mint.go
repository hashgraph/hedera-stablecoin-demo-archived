package operation

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Mint(payload *pb.MintTo) (domain.Operation, error) {
	log.Trace().
		Str("username", payload.Address).
		Uint64("quantity", payload.Quantity).
		Msg("Mint")

	if _, exists := state.Balance.Load(payload.Address); !exists {
		return domain.Operation{
			Operation:     domain.OpMint,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user `%s` does not exist", payload.Address),
		}, nil
	}

	// TODO: Handle response to the UI

	// FIXME: UI sends the username where it calls it the address
	// NOTE: I (@mehcode) prefer the username here, but we should change the field name
	state.UpdateBalance(payload.Address, func(balance uint64) uint64 {
		return balance + payload.Quantity
	})

	userPublicKey, _ := state.User.Load(payload.Address)
	userPublicKeyBytes := userPublicKey.([]byte)

	return domain.Operation{
		Operation: domain.OpMint,
		Status:    domain.OpStatusComplete,
		ToAddress: &userPublicKeyBytes,
		Amount:    int64(payload.Quantity),
	}, nil
}
