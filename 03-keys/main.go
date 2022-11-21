package main

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

func main() {
	lgr := logrus.New()
	// generate private key
	pvKey, err := crypto.GenerateKey()
	if err != nil {
		lgr.WithError(err).Fatal("can't generate private key")
	}

	// string of private key
	pvData := crypto.FromECDSA(pvKey)
	// encodes private key data as a hex string with 0x prefix.
	pvKeyString := hexutil.Encode(pvData)

	// generate public key
	pubData := crypto.FromECDSAPub(&pvKey.PublicKey)
	// encodes public key data as a hex string with 0x prefix.
	pubKeyString := hexutil.Encode(pubData)

	// generate hex public address
	pubAddress := crypto.PubkeyToAddress(pvKey.PublicKey).Hex()

	lgr.WithFields(logrus.Fields{
		"private_key":    pvKeyString,
		"public_key":     pubKeyString,
		"public_address": pubAddress,
	}).Info("print keys")

}
