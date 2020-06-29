package data

func GetUserBalanceByAddress(publicKey []byte) (int64, bool, error) {
	var balance int64 = 0
	var frozen bool = false
	err := db.QueryRow(`
SELECT balance, frozen
FROM address
WHERE public_key = $1
	`, publicKey).Scan(&balance, &frozen)

	return balance, frozen, err
}

func UpdateUserBalance(userName string, newBalance uint64) error {
	_, err := db.Exec("UPDATE address SET balance = $1 WHERE username = $2", newBalance, userName)
	return err
}
