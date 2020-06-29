package operation

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Transfer(senderAddress []byte, payload *pb.Transfer) (domain.Operation, error) {
	// FIXME: proto says address but we are receiving a username
	log.Trace().
		Str("to", payload.ToAddress).
		Uint64("quantity", payload.Quantity).
		Msg("Transfer")

	var senderUserName string
	var exists bool
	senderAddressHex := hex.EncodeToString(senderAddress)

	if senderUserName, exists = state.Address[senderAddressHex]; !exists {
		return domain.Operation{
			Operation:     domain.OpTransfer,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("address `%s` does not exist", senderAddressHex),
			FromAddress:   &senderAddress,
		}, nil
	}

	if _, exists = state.Balance[payload.ToAddress]; !exists {
		return domain.Operation{
			Operation:     domain.OpTransfer,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user `%s` does not exist", payload.ToAddress),
			FromAddress:   &senderAddress,
		}, nil
	}

	toAddress := []byte(state.User[payload.ToAddress])
	senderBalance := state.Balance[senderUserName]

	if senderBalance < payload.Quantity {
		return domain.Operation{
			Operation:     domain.OpTransfer,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user `%s` has an insufficient balance", senderUserName),
			FromAddress:   &senderAddress,
			ToAddress:     &toAddress,
		}, nil
	}

	// TODO: Handle response to the UI

	state.UpdateBalance(payload.ToAddress, func(balance uint64) uint64 {
		return balance + payload.Quantity
	})

	state.UpdateBalance(senderUserName, func(balance uint64) uint64 {
		return balance - payload.Quantity
	})

	return domain.Operation{
		Operation:   domain.OpTransfer,
		Status:      domain.OpStatusComplete,
		FromAddress: &senderAddress,
		ToAddress:   &toAddress,
		Amount:      int64(payload.Quantity),
	}, nil
}
