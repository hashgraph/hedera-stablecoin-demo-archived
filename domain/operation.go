package domain

type Op string

const (
	OpAnnounce Op = "announce"
	OpMint     Op = "mint"
	OpTransfer Op = "transfer"
	OpRedeem   Op = "redeem"
	OpFreeze   Op = "freeze"
	OpUnFreeze Op = "unfreeze"
	OpClawback   Op = "clawback"
)

type OpStatus string

const (
	OpStatusComplete OpStatus = "complete"
	OpStatusFailed   OpStatus = "failed"
)

type Operation struct {
	Consensus int64 `db:"consensus" json:"consensus"`

	Operation Op `db:"operation" json:"operation"`

	Signature []byte `db:"signature" json:"-"`

	ToAddress   *[]byte `db:"to_address" json:"-"`
	FromAddress *[]byte `db:"from_address" json:"-"`

	ToUsername   *string `db:"to_username" json:"toUsername"`
	FromUsername *string `db:"from_username" json:"fromUsername"`

	Amount int64 `db:"amount" json:"amount"`

	Status        OpStatus `db:"status" json:"status"`
	StatusMessage string   `db:"status_message" json:"statusMessage"`
}
