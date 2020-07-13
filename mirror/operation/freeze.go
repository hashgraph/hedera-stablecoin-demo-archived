package operation

import (
	"crypto/ed25519"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/api/notification"
	"github.io/hashgraph/stable-coin/mirror/state"
)

func Freeze(adminAddress []byte, account string, freeze bool) (domain.Operation, error) {
	operation := domain.OpFreeze
	operationString := "freeze"
	if freeze {
		log.Trace().
			Str("username", account).
			Msg("Freeze")
	} else {
		operation = domain.OpUnFreeze
		operationString = "unfreeze"
		log.Trace().
			Str("username", account).
			Msg("UnFreeze")
	}

	var userKeyI interface{}
	var exists bool
	if userKeyI, exists = state.User.Load(account); ! exists {
		statusMessage := fmt.Sprintf("user to %s `%s` does not exist", operationString, account)
		notification.SendNotification("Admin", true, statusMessage)
		return domain.Operation{
			Operation:     operation,
			Status:        domain.OpStatusFailed,
			StatusMessage: statusMessage,
		}, nil
	}
	userKey := []byte(userKeyI.(ed25519.PublicKey))

	state.UpdateFrozen(account, func(frozen bool) bool {
		return freeze
	})

	statusMessage := fmt.Sprintf("user %s for `%s` successful", operationString, account)
	notification.SendNotification("Admin", false, statusMessage)

	return domain.Operation{
		Operation:   operation,
		Status:      domain.OpStatusComplete,
		FromAddress: &adminAddress,
		ToAddress:   &userKey,
		Amount:      int64(0),
	}, nil
}
