package operation

import (
	"crypto/ed25519"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Clawback(payload *pb.Clawback) (domain.Operation, error) {
	// FIXME: proto says address but we are receiving a username
	log.Trace().
		Str("from", payload.Account).
		Msg("Clawback")

	var clawBackUserAddressI interface{}
	var exists bool

	if clawBackUserAddressI, exists = state.User.Load(payload.Account); !exists {
		return domain.Operation{
			Operation:     domain.OpClawback,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("username `%s` does not exist", payload.Account),
		}, nil
	}
	clawBackUserAddress := []byte(clawBackUserAddressI.(ed25519.PublicKey))

	senderBalanceI, _ := state.Balance.Load(payload.Account)
	senderBalance := senderBalanceI.(uint64)

	// TODO: Handle response to the UI

	state.UpdateBalance(payload.Account, func(balance uint64) uint64 {
		return 0
	})

	return domain.Operation{
		Operation:   domain.OpClawback,
		Status:      domain.OpStatusComplete,
		ToAddress:   &clawBackUserAddress,
		Amount:      int64(-senderBalance),
	}, nil
}
