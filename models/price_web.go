package models

import "github.com/shopspring/decimal"

type Price struct {
	Value decimal.Decimal `json:"value"`
	Currency string `json:"currency"`
}
