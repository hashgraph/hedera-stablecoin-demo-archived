package data

import (
	"github.io/hashgraph/stable-coin/domain"
	"time"
)

func GetOperationsForUsername(username string) ([]domain.Operation, error) {
	operations := []domain.Operation{}
	err := db.Select(&operations, `
SELECT operation.*, from_address.username as from_username, to_address.username as to_username
FROM operation
INNER JOIN address from_address ON from_address.public_key = operation.from_address
INNER JOIN address to_address ON to_address.public_key = operation.to_address
WHERE from_address.username = $1 
   OR to_address.username = $1
ORDER BY operation.consensus DESC
	`, username)

	return operations, err
}

func GetLatestOperationConsensus() (time.Time, error) {
	var consensusNanos int64
	err := db.Get(&consensusNanos, "SELECT consensus FROM operation ORDER BY consensus DESC LIMIT 1")

	return time.Unix(0, consensusNanos), err
}

func InsertOperation(op domain.Operation) error {
	_, err := db.Exec(`
INSERT INTO operation (
   consensus,
   "operation",
   signature,
   from_address,
   to_address,
   amount,
   status,
   status_message
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
	`, op.Consensus, op.Operation, op.Signature, op.FromAddress, op.ToAddress, op.Amount, op.Status, op.StatusMessage)

	return err
}
