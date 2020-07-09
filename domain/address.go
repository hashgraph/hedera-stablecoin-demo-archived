package domain

type Address struct {
	Balance   int64  `db:"balance"`
	PublicKey []byte `db:"public_key"`
	Username  string `db:"username"`
	Frozen bool `db:"frozen"`
}
