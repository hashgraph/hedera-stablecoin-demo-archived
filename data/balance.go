package data

func GetUserBalanceByUsername(username string) (int64, bool, error) {
	var balance int64 = 0
	var frozen bool = false
	err := db.QueryRow(`
SELECT balance, frozen
FROM address
WHERE username = $1
	`, username).Scan(&balance, &frozen)

	return balance, frozen, err
}

func UpdateUserBalances(balances map[string]uint64) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for userName, newBalance := range balances {
		_, err := tx.Exec("UPDATE address SET balance = $1 WHERE username = $2", newBalance, userName)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
