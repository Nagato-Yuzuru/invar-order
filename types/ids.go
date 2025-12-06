package types

import (
	"github.com/shopspring/decimal"
)

type UserID string
type QuoteID string

type ProductCode string

// InternalSKUID is an internal unique identifier for the project
type InternalSKUID string

type Money decimal.Decimal

type Unit string

type ProjectScope string
