package operation

import (
	"encoding/hex"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Announce(payload *pb.Join) error {
	// FIXME: address should be transmitted as raw bytes to remove the parsing work needed here
	publicKey, err := hedera.Ed25519PublicKeyFromString(payload.Address)
	if err != nil {
		return err
	}

	publicKeyHex := hex.EncodeToString(publicKey.Bytes())

	log.Trace().
		Str("username", payload.Username).
		Str("key", publicKeyHex).
		Msg("Announce")

	// TODO: Handle "user already exists"
	// TODO: Handle recording the operation
	// TODO: Handle response to the UI

	state.User[payload.Username] = publicKey.Bytes()
	state.Address[publicKeyHex] = payload.Username
	state.Balance[payload.Username] = 0

	return nil
}
