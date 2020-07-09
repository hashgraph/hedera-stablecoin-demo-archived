package operation

import (
	"crypto/ed25519"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
)

func Freeze(adminAddress []byte, account string, freeze bool) (domain.Operation, error) {
	operation := domain.OpFreeze
	if freeze {
		log.Trace().
			Str("username", account).
			Msg("Freeze")
	} else {
		operation = domain.OpUnFreeze
		log.Trace().
			Str("username", account).
			Msg("UnFreeze")
	}

	var userKeyI interface{}
	var exists bool
	if userKeyI, exists = state.User.Load(account); ! exists {
		return domain.Operation{
			Operation:     operation,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user to freeze `%s` does not exist", account),
		}, nil
	}
	userKey := []byte(userKeyI.(ed25519.PublicKey))

	// TODO: Handle response to the UI

	state.UpdateFrozen(account, func(frozen bool) bool {
		return freeze
	})

	return domain.Operation{
		Operation:   operation,
		Status:      domain.OpStatusComplete,
		FromAddress: &adminAddress,
		ToAddress:   &userKey,
		Amount:      int64(0),
	}, nil
}
