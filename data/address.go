package data

import "github.io/hashgraph/stable-coin/domain"

func GetAllAddress() ([]domain.Address, error) {
	var r []domain.Address
	err := db.Select(&r, "SELECT username, balance, public_key FROM address")

	return r, err
}

func InsertNewAddress(username string, publicKey []byte) error {
	_, err := db.Exec(`INSERT INTO address (username, public_key, balance) VALUES ($1, $2, 0)`,
		username, publicKey)

	return err
}
