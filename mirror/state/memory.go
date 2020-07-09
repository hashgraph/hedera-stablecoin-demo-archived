package state

import (
	"crypto/ed25519"
	"encoding/hex"
	"sync"

	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
)

// username -> balance
// var Balance = map[string]uint64{}
var Balance = sync.Map{}

// username -> address (public key)
// var User = map[string]ed25519.PublicKey{}
var User = sync.Map{}

// address (public key hex) -> username
// var Address = map[string]string{}
var Address = sync.Map{}

// frozen (username) -> boolean
var Frozen = sync.Map{}

// pending new users (usernames)
var pendingNewUser []string
var pendingNewUserLock sync.Mutex

// pending freezes (user, status)
var pendingFreezes = map[string]bool{}

// pending balance changes
var pendingBalances = map[string]uint64{}
var pendingBalancesLock sync.Mutex

// pending operations to be committed to the database
var pendingOperations []domain.Operation
var pendingOperationsLock sync.Mutex
var pendingOperationsForUser = map[string][]domain.Operation{}

func init() {
	addressRows, err := data.GetAllAddress()
	if err != nil {
		panic(err)
	}

	for _, row := range addressRows {
		Balance.Store(row.Username, uint64(row.Balance))
		User.Store(row.Username, ed25519.PublicKey(row.PublicKey))
		Address.Store(hex.EncodeToString(row.PublicKey), row.Username)
		Frozen.Store(row.Username, row.Frozen)
	}
}

// AddUser adds a new user
func AddUser(username string, publicKey ed25519.PublicKey) {
	Balance.Store(username, uint64(0))
	User.Store(username, publicKey)
	Address.Store(hex.EncodeToString(publicKey), username)
	Frozen.Store(username, false)

	// on the next commit, add the user
	pendingNewUserLock.Lock()
	pendingNewUser = append(pendingNewUser, username)
	pendingNewUserLock.Unlock()
}

// UpdateBalance updates the balance for a user and ensures that
// it is eventually persisted
func UpdateBalance(userName string, update func(uint64) uint64) {
	v, _ := Balance.Load(userName)
	Balance.Store(userName, update(v.(uint64)))
	v, _ = Balance.Load(userName)

	// on the next commit, update our balance
	pendingBalancesLock.Lock()
	pendingBalances[userName] = v.(uint64)
	pendingBalancesLock.Unlock()
}

// UpdateFrozen updates the frozen status for a user and ensures that
// it is eventually persisted
func UpdateFrozen(userName string, update func(bool) bool) {
	v, _ := Frozen.Load(userName)
	Frozen.Store(userName, update(v.(bool)))
	v, _ = Frozen.Load(userName)

	// on the next commit, update the frozen status
	pendingFreezes[userName] = v.(bool)
}

// AddOperation adds an operation to the pending store to be committed on the commit interval
func AddOperation(op domain.Operation) {
	pendingOperationsLock.Lock()
	pendingOperations = append(pendingOperations, op)
	pendingOperationsLock.Unlock()

	if op.FromAddress != nil {
		if fromUserName, ok := Address.Load(hex.EncodeToString(*op.FromAddress)); ok {
			fromUserNameS := fromUserName.(string)
			pendingOperationsForUser[fromUserNameS] = append(pendingOperationsForUser[fromUserNameS], op)
		}
	}

	if op.ToAddress != nil && op.FromAddress != op.ToAddress {
		if toUserName, ok := Address.Load(hex.EncodeToString(*op.ToAddress)); ok {
			toUserNameS := toUserName.(string)
			pendingOperationsForUser[toUserNameS] = append(pendingOperationsForUser[toUserNameS], op)
		}
	}
}

func GetPendingOperationsForUser(userName string) []domain.Operation {
	return pendingOperationsForUser[userName]
}
