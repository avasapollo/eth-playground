package main

import (
	"os"
	"path"

	"github.com/avasapollo/eth-playground/config"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()
	lgr := logrus.New()

	// generate keystore
	key := keystore.NewKeyStore(c.KeyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	// generate an account
	acc, err := key.NewAccount(c.KeyStoreAccountPassword)
	if err != nil {
		lgr.WithError(err).Fatal("can't create the account")
	}
	lgr.WithField("account_url", acc.URL).Info("account key stored")

	// read the account file to get private, public key and public address of the account
	_, file := path.Split(acc.URL.Path)
	accountFile := c.KeyStorePath + "/" + file
	b, err := os.ReadFile(accountFile)
	if err != nil {
		lgr.WithError(err).Fatal("can't read account file from keystore")
	}
	pvKey, err := keystore.DecryptKey(b, c.KeyStoreAccountPassword)
	if err != nil {
		lgr.WithError(err).Fatal("can't decrypt the json account key")
	}

	// get private key
	privateKey := hexutil.Encode(crypto.FromECDSA(pvKey.PrivateKey))
	// get public key from private key
	publicKey := hexutil.Encode(crypto.FromECDSAPub(&pvKey.PrivateKey.PublicKey))
	// get public address from public key
	pubAddres := crypto.PubkeyToAddress(pvKey.PrivateKey.PublicKey).Hex()
	lgr.WithFields(logrus.Fields{
		"account_private_key":    privateKey,
		"account_public_key":     publicKey,
		"account_public_address": pubAddres,
	}).Info("account details")
}
