package convert

import (
	"crypto/ecdsa"
	"crypto/elliptic"
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

func ToPrivateKey(raw string) *ecdsa.PrivateKey {
	k := new(big.Int)
	k.SetString(raw, 16)

	prKey := new(ecdsa.PrivateKey)
	curve := elliptic.P256()
	prKey.PublicKey.Curve = curve
	prKey.D = k
	prKey.PublicKey.X, prKey.PublicKey.Y = curve.ScalarBaseMult(k.Bytes())
	return prKey
}
