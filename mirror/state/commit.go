package state

import (
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
	"time"
)

var commitInterval = "1s"

func init() {
	commitIntervalDur, err := time.ParseDuration(commitInterval)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			time.Sleep(commitIntervalDur)
			go commit()
		}
	}()
}

func commit() {
	start := time.Now()
	numOperations := 0
	numBalanaces := 0
	numUsers := 0

	if len(pendingNewUser) > 0 {
		// there are pending operations that should be committed
		users := pendingNewUser
		numUsers = len(users)

		// erase current maps
		pendingNewUser = nil

		// insert the new user records
		for _, userName := range users {
			err := data.InsertNewAddress(userName, User[userName])
			if err != nil {
				panic(err)
			}
		}
	}

	if len(pendingOperations) > 0 {
		// there are pending operations that should be committed
		operations := pendingOperations
		numOperations = len(operations)

		// erase current maps
		pendingOperations = nil
		pendingOperationsForUser = map[string][]domain.Operation{}

		// iterate and insert all the new operations
		for _, op := range operations {
			err := data.InsertOperation(op)
			if err != nil {
				panic(err)
			}
		}
	}

	if len(pendingBalances) > 0 {
		// there are pending operations that should be committed
		balances := pendingBalances
		numBalanaces = len(balances)

		// erase current maps
		pendingBalances = map[string]uint64{}

		// update the balance records
		for userName, newBalance := range balances {
			err := data.UpdateUserBalance(userName, newBalance)
			if err != nil {
				panic(err)
			}
		}
	}

	if numOperations > 0 || numBalanaces > 0 {
		log.Info().
			Dur("elapsed", time.Since(start)).
			Int("operations", numOperations).
			Int("users", numUsers).
			Int("balances", numBalanaces).
			Msg("Commit")
	}
}
