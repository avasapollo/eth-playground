package convert

import (
	"errors"
	"math"
	"math/big"
)

type Unit string

const Ether Unit = "ether"

func FromWei(unit Unit, val *big.Int) (*big.Float, error) {
	if unit != Ether {
		return nil, errors.New("unit is not allowed")
	}
	// 1 ether = 10^18 wei
	fBalance := new(big.Float)
	fBalance.SetString(val.String())
	return new(big.Float).Quo(fBalance, big.NewFloat(math.Pow10(18))), nil
}
