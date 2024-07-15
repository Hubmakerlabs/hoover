package schema

import "math/big"

type Receipt struct {
	AccId    string
	TokenTag string
	Amount   *big.Int
}
