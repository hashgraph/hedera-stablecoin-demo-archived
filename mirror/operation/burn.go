package operation

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Burn(senderAddress []byte, payload *pb.Burn) (domain.Operation, error) {
	senderAddressHex := hex.EncodeToString(senderAddress)

	log.Trace().
		Str("from", senderAddressHex).
		Uint64("quantity", payload.Amount).
		Msg("Redeem")

	var senderUserNameI interface{}
	var exists bool

	if senderUserNameI, exists = state.Address.Load(senderAddressHex); !exists {
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("address `%s` does not exist", senderAddressHex),
			FromAddress:   &senderAddress,
		}, nil
	}

	senderUserName := senderUserNameI.(string)

	if _, exists = state.Balance.Load(senderUserName); !exists {
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user `%s` does not exist", senderUserName),
			FromAddress:   &senderAddress,
		}, nil
	}

	if frozenUserI, exists := state.Frozen.Load(senderUserName); exists {
		if frozenUserI.(bool) == true {
			return domain.Operation{
				Operation:     domain.OpRedeem,
				Status:        domain.OpStatusFailed,
				StatusMessage: fmt.Sprintf("user `%s` is frozen", senderUserName),
				FromAddress:   &senderAddress,
			}, nil
		}
	}

	senderBalanceI, _ := state.Balance.Load(senderUserName)
	senderBalance := senderBalanceI.(uint64)

	if senderBalance < payload.Amount {
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user `%s` has an insufficient balance", senderUserName),
			FromAddress:   &senderAddress,
		}, nil
	}

	// TODO: Handle response to the UI

	state.UpdateBalance(senderUserName, func(balance uint64) uint64 {
		return balance - payload.Amount
	})

	return domain.Operation{
		Operation:   domain.OpRedeem,
		Status:      domain.OpStatusComplete,
		FromAddress: &senderAddress,
		Amount:      int64(payload.Amount),
	}, nil
}
