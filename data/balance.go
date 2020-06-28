package data

import (
	"context"
)

func GetUserBalanceByAddress(publicKey []byte) (int64, bool, error) {
	var balance int64 = 0
	var frozen bool = false
	err := db.QueryRow(context.TODO(),`
SELECT balance, frozen
FROM address
WHERE public_key = $1
	`, publicKey).Scan(&balance, &frozen)

	return balance, frozen, err
}
