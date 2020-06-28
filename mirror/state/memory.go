package state

import "crypto/ed25519"

// username -> balance
var Balances = map[string]uint64{}

// username -> address (public key)
var Users = map[string]ed25519.PublicKey{}
