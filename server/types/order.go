package types

import "time"

const (
	ActionTypeOrder = "order"
)

type DbOrder struct {
	Time time.Time
	User string
	Type string
	Order
	Block int64
	Hash  string `gorm:"primary_key"`
}

type Transaction struct {
	Time   int64       `json:"time"`
	User   string      `json:"user"`
	Action ActionOrder `json:"action"`
	Block  int64       `json:"block"`
	Hash   string      `json:"hash"`
	Error  *string     `json:"error"`
}

type Action struct {
	Type string `json:"type"`
}

type ActionOrder struct {
	Type     string  `json:"type"`
	Orders   []Order `json:"orders"`
	Grouping string  `json:"grouping"`
}

type Order struct {
	// asset id
	A int64 `json:"a"`
	// is buy
	B bool `json:"b"`
	// price
	P string `json:"p"`
	// size
	S string `json:"s"`
	// Is reduce-only
	R bool `json:"r"`
	// order type
	T T `json:"t"`
	// Client Order ID.
	C string `json:"c"`
}

type T struct {
	Limit   *Limit   `json:"limit"`
	Trigger *Trigger `json:"trigger"`
}

type Limit struct {
	Tif string `json:"tif"`
}

type Trigger struct {
	// Is market order
	IsMarket bool `json:"isMarket"`
	// Trigger price
	TriggerPx string `json:"triggerPx"`
	// Indicates whether it is take profit or stop loss.
	Tpsl string `json:"tpsl"`
}
