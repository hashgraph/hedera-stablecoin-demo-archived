package data

import (
	"context"
	"crypto/ed25519"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.io/hashgraph/stable-coin/domain"
	"sync"
)

func GetAllAddress() ([]domain.Address, error) {
	var r []domain.Address
	err := db.Select(&r, "SELECT username, balance, public_key FROM address")

	return r, err
}

func InsertNewAddresses(newUsers []string, userToAddress *sync.Map) error {
	var rows = make([][]interface{}, 0, len(newUsers))

	for _, newUser := range newUsers {
		address, _ := userToAddress.Load(newUser)

		rows = append(rows, []interface{}{
			newUser,
			address.(ed25519.PublicKey),
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

	_, err = conn.CopyFrom(context.Background(), pgx.Identifier{"address"}, []string{
		"username",
		"public_key",
	}, pgx.CopyFromRows(rows))

	return nil
}
func UpdateUserFrozenStatus(users map[string]bool) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	for userName, frozen := range users {
		_, err := tx.Exec("UPDATE address SET frozen = $1 WHERE username = $2", frozen, userName)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
