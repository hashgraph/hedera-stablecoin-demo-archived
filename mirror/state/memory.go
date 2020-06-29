package state

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
)

// username -> balance
var Balance = map[string]uint64{}

// username -> address (public key)
var User = map[string]ed25519.PublicKey{}

// address (public key hex) -> username
var Address = map[string]string{}

// pending new users (usernames)
var pendingNewUser []string

// pending balance changes
var pendingBalances = map[string]uint64{}

// pending operations to be committed to the database
var pendingOperations []domain.Operation
var pendingOperationsForUser = map[string][]domain.Operation{}

func init() {
	addressRows, err := data.GetAllAddress()
	if err != nil {
		panic(err)
	}

	fmt.Printf("address? %v\n")

	for _, row := range addressRows {
		Balance[row.Username] = uint64(row.Balance)
		User[row.Username] = row.PublicKey
		Address[hex.EncodeToString(row.PublicKey)] = row.Username
	}
}

// AddUser adds a new user
func AddUser(username string, publicKey ed25519.PublicKey) {
	User[username] = publicKey
	Address[hex.EncodeToString(publicKey)] = username
	Balance[username] = 0

	// on the next commit, add the user
	pendingNewUser = append(pendingNewUser, username)
}

// UpdateBalance updates the balance for a user and ensures that
// it is eventually persisted
func UpdateBalance(userName string, update func(uint64) uint64) {
	Balance[userName] = update(Balance[userName])

	// on the next commit, update our balance
	pendingBalances[userName] = Balance[userName]
}

// AddOperation adds an operation to the pending store to be committed on the commit interval
func AddOperation(op domain.Operation) {
	pendingOperations = append(pendingOperations, op)

	if op.FromAddress != nil {
		if fromUserName, ok := Address[hex.EncodeToString(*op.FromAddress)]; ok {
			pendingOperationsForUser[fromUserName] = append(pendingOperationsForUser[fromUserName], op)
		}
	}

	if op.ToAddress != nil && op.FromAddress != op.ToAddress {
		if toUserName, ok := Address[hex.EncodeToString(*op.ToAddress)]; ok {
			pendingOperationsForUser[toUserName] = append(pendingOperationsForUser[toUserName], op)
		}
	}
}

func GetPendingOperationsForUser(userName string) []domain.Operation {
	return pendingOperationsForUser[userName]
}
