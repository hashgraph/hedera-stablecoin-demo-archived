package operation

import (
	"encoding/hex"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/api/notification"
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
		statusMessage := fmt.Sprintf("address `%s` does not exist", senderAddressHex)
		// can't notify an unknown user
		// notification.SendNotification(payload.Address, true, statusMessage)
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: statusMessage,
			FromAddress:   &senderAddress,
		}, nil
	}

	senderUserName := senderUserNameI.(string)

	if _, exists = state.Balance.Load(senderUserName); !exists {
		statusMessage := fmt.Sprintf("user `%s` does not exist", senderUserName)
		notification.SendNotification(senderUserName, true, statusMessage)
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: statusMessage,
			FromAddress:   &senderAddress,
		}, nil
	}

	if frozenUserI, exists := state.Frozen.Load(senderUserName); exists {
		if frozenUserI.(bool) == true {
			statusMessage := fmt.Sprintf("user `%s` is frozen", senderUserName)
			notification.SendNotification(senderUserName, true, statusMessage)
			return domain.Operation{
				Operation:     domain.OpRedeem,
				Status:        domain.OpStatusFailed,
				StatusMessage: statusMessage,
				FromAddress:   &senderAddress,
			}, nil
		}
	}

	senderBalanceI, _ := state.Balance.Load(senderUserName)
	senderBalance := senderBalanceI.(uint64)

	if senderBalance < payload.Amount {
		statusMessage := fmt.Sprintf("user `%s` has an insufficient balance", senderUserName)
		notification.SendNotification(senderUserName, true, statusMessage)
		return domain.Operation{
			Operation:     domain.OpRedeem,
			Status:        domain.OpStatusFailed,
			StatusMessage: statusMessage,
			FromAddress:   &senderAddress,
		}, nil
	}

	state.UpdateBalance(senderUserName, func(balance uint64) uint64 {
		return balance - payload.Amount
	})

	statusMessage := fmt.Sprintf("redeem successful")
	notification.SendNotification(senderUserName, false, statusMessage)

	return domain.Operation{
		Operation:   domain.OpRedeem,
		Status:      domain.OpStatusComplete,
		FromAddress: &senderAddress,
		Amount:      int64(payload.Amount),
	}, nil
}
