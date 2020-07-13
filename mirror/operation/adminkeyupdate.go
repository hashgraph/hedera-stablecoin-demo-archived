package operation

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/api/notification"
	"github.io/hashgraph/stable-coin/mirror/state"
)

func AdminKeyUpdate(adminAddress []byte, newAdminPublicKey string) (domain.Operation, error) {
	log.Trace().
		Str("username", "Admin").
		Str("keyUpdate", newAdminPublicKey).
		Msg("Admin Key Update")

	newHederaPublicKey, err := hedera.Ed25519PublicKeyFromString(newAdminPublicKey)
	newPublicKey := ed25519.PublicKey(newHederaPublicKey.Bytes())
	if err != nil {
		return domain.Operation{}, err
	}

	newPublicKeyBytes := []byte(newPublicKey)
	adminAddressHex := hex.EncodeToString(adminAddress)

	state.UpdateAdminKey(adminAddressHex, func(newKey ed25519.PublicKey) ed25519.PublicKey {
		return newPublicKey
	})

	statusMessage := fmt.Sprintf("admin key successfully replaced")
	notification.SendNotification("Admin", false, statusMessage)

	return domain.Operation{
		Operation:   domain.OpAdminKeyUpdate,
		Status:      domain.OpStatusComplete,
		FromAddress: &adminAddress,
		ToAddress:   &newPublicKeyBytes,
		Amount:      int64(0),
	}, nil
}
