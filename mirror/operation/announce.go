package operation

import (
	"encoding/hex"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Announce(payload *pb.Join) (domain.Operation, error) {
	// FIXME: address should be transmitted as raw bytes to remove the parsing work needed here
	publicKey, err := hedera.Ed25519PublicKeyFromString(payload.Address)
	if err != nil {
		return domain.Operation{}, err
	}

	publicKeyHex := hex.EncodeToString(publicKey.Bytes())

	log.Trace().
		Str("username", payload.Username).
		Str("key", publicKeyHex).
		Msg("Announce")

	// TODO: Handle response to the UI

	if _, exists := state.User.Load(payload.Username); exists {
		// duplicate user name
		return domain.Operation{
			Operation:     domain.OpAnnounce,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("user name `%s` already exists", payload.Username),
		}, nil
	}

	if _, exists := state.Address.Load(publicKeyHex); exists {
		// duplicate public key
		return domain.Operation{
			Operation:     domain.OpAnnounce,
			Status:        domain.OpStatusFailed,
			StatusMessage: fmt.Sprintf("address `%s` already exists", publicKeyHex),
		}, nil
	}

	state.AddUser(payload.Username, publicKey.Bytes())

	return domain.Operation{
		Operation: domain.OpAnnounce,
		Status:    domain.OpStatusComplete,
	}, nil
}
