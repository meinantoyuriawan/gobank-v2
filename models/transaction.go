package models

type Transfer struct {
	Source    int64  `json:"source"`
	Pin       string `json:"pin"`
	Recipient int64  `json:"recipient"`
	Amount    int64  `json:"amount"`
}
