package tools

import (
	"errors"
	"math"
	"math/big"
)

const Ether = "ether"

func FromWei(unit string, val *big.Int) (*big.Float, error) {
	if unit != Ether {
		return nil, errors.New("unit is not allowed")
	}
	// 1 ether = 10^18 wei
	fBalance := new(big.Float)
	fBalance.SetString(val.String())
	return new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18))), nil
}

