package event

import "time"

type BalanceUpdated struct {
	Name    string
	Payload interface{}
}

func NewBalanceUpdated() *BalanceUpdated {
	return &BalanceUpdated{
		Name:    "BalanceUpdated",
		Payload: nil,
	}
}

func (b *BalanceUpdated) GetName() string {
	return b.Name
}

func (b *BalanceUpdated) GetDatetime() time.Time {
	return time.Now()
}

func (b *BalanceUpdated) GetPayload() interface{} {
	return b.Payload
}

func (b *BalanceUpdated) SetPayload(p interface{}) {
	b.Payload = p
}
