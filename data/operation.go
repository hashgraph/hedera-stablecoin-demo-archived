package data

import (
	"context"
	"sync"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.io/hashgraph/stable-coin/domain"
)

// returns operation ID -> true as there is no sync.Set type
func GetAllOperationIDs() (*sync.Map, error) {
	operationIDs := new(sync.Map)

	rows, err := db.Query(`
SELECT encode(operation.signature, 'hex')
FROM operation
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}

		operationIDs.Store(id, true)
	}

	return operationIDs, nil
}

func GetOperationsForUsername(username string) ([]domain.Operation, error) {
	operations := []domain.Operation{}
	err := db.Select(&operations, `
SELECT operation.*, from_address.username as from_username, to_address.username as to_username
FROM operation
LEFT JOIN address from_address ON from_address.public_key = operation.from_address
LEFT JOIN address to_address ON to_address.public_key = operation.to_address
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

func InsertOperations(operations []domain.Operation) error {
	var rows = make([][]interface{}, 0, len(operations))

	for _, op := range operations {
		rows = append(rows, []interface{}{
			op.Consensus,
			string(op.Operation),
			op.Signature,
			op.FromAddress,
			op.ToAddress,
			op.Amount,
			string(op.Status),
			op.StatusMessage,
		})
	}

	conn, err := stdlib.AcquireConn(db.DB)

	defer func() {
		err := stdlib.ReleaseConn(db.DB, conn)
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		return err
	}

	_, err = conn.CopyFrom(context.Background(), pgx.Identifier{"operation"}, []string{
		"consensus",
		"operation",
		"signature",
		"from_address",
		"to_address",
		"amount",
		"status",
		"status_message",
	}, pgx.CopyFromRows(rows))

	return err
}
