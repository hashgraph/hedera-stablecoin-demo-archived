package data

import (
	"context"
	"crypto/ed25519"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.io/hashgraph/stable-coin/domain"
)

func GetAllAddress() ([]domain.Address, error) {
	var r []domain.Address
	err := db.Select(&r, "SELECT username, balance, public_key FROM address")

	return r, err
}

func InsertNewAddresses(newUsers []string, userToAddress map[string]ed25519.PublicKey) error {
	var rows = make([][]interface{}, 0, len(newUsers))

	for _, newUser := range newUsers {
		rows = append(rows, []interface{}{
			newUser,
			userToAddress[newUser],
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
