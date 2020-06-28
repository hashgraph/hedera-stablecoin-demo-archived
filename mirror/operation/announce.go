package operation

import (
	"github.com/hashgraph/hedera-sdk-go"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
)

func Announce(payload *pb.Join) error {
	// FIXME: address should be transmitted as raw bytes to remove the parsing work needed here
	publicKey, err := hedera.Ed25519PublicKeyFromString(payload.Address)
	if err != nil {
		return err
	}

	// TODO: Handle "user already exists"
	// TODO: Handle recording the operation
	// TODO: Handle response to the UI

	state.Users[payload.Username] = publicKey.Bytes()

	return nil
}
