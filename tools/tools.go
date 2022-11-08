package tools

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params"
)

func EtherToWei(val *big.Int) *big.Int {
	return new(big.Int).Mul(val, big.NewInt(params.Ether))
}

func WeiToEther(val *big.Int) *big.Int {
	return new(big.Int).Div(val, big.NewInt(params.Ether))
}
