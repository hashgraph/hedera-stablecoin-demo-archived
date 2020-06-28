package state

import "crypto/ed25519"

// username -> balance
var Balance = map[string]uint64{}

// username -> address (public key)
var User = map[string]ed25519.PublicKey{}

// address (public key hex) -> username
var Address = map[string]string{}
