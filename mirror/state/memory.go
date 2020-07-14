package state

import (
	"crypto/ed25519"
	"encoding/hex"
	"sync"

	"github.com/hashgraph/hedera-sdk-go"

	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
)

// store of operation "ID"s so we don't process duplicate operations
// in the future, we want to tweak the message body so its linked to the Transaction ID
// and we can let Hedera handle duplicates
var operationIDs *sync.Map

// username -> balance
var Balance = sync.Map{}

// username -> address (public key)
var User = sync.Map{}

// address (public key hex) -> username
var Address = sync.Map{}

// frozen (username) -> boolean
var Frozen = sync.Map{}

// pending new users (usernames)
var pendingNewUser []string
var pendingNewUserLock sync.Mutex

// pending freezes (user, status)
var pendingFreezes = map[string]bool{}
var pendingFreezesLock sync.Mutex

// pending key updates (old key)
var pendingKeyUpdates = map[string]ed25519.PublicKey{}
var pendingKeyUpdatesLock sync.Mutex

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

	operationIDs, err = data.GetAllOperationIDs()
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

// OperationExists checks for an existing operation ID
func OperationExists(signatureHex string) bool {
	_, exists := operationIDs.Load(signatureHex)
	return exists
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

// UpdateAdminKey updates the public key for the admin user and ensures that
// it is eventually persisted
func UpdateAdminKey(oldAdminAddress string, update func(ed25519.PublicKey) ed25519.PublicKey) {
	v, _ := User.Load("Admin")
	// replace the key for the Admin user
	publicKey, _ := hedera.Ed25519PublicKeyFromBytes(update(v.(ed25519.PublicKey)))
	User.Store("Admin", update(v.(ed25519.PublicKey)))
	// delete the old key entry
	Address.Delete(oldAdminAddress)
	// create an entry for the new one
	Address.Store(hex.EncodeToString(update(v.(ed25519.PublicKey))), "Admin")

	// on the next commit, update the frozen status
	pendingKeyUpdatesLock.Lock()
	pendingKeyUpdates[oldAdminAddress] = publicKey.Bytes()
	pendingKeyUpdatesLock.Unlock()
}

// UpdateFrozen updates the frozen status for a user and ensures that
// it is eventually persisted
func UpdateFrozen(userName string, update func(bool) bool) {
	v, _ := Frozen.Load(userName)
	Frozen.Store(userName, update(v.(bool)))
	v, _ = Frozen.Load(userName)

	// on the next commit, update the frozen status
	pendingFreezesLock.Lock()
	pendingFreezes[userName] = v.(bool)
	pendingFreezesLock.Unlock()
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

	// store the signature so we don't accept this operation again
	operationIDs.Store(hex.EncodeToString(op.Signature), true)
}

func GetPendingOperationsForUser(userName string) []domain.Operation {
	return pendingOperationsForUser[userName]
}
