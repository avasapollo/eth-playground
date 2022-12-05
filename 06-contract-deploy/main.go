package main

import (
	"context"
	"math/big"
	"os"

	"github.com/avasapollo/eth-playground/config"
	todo "github.com/avasapollo/eth-playground/gen"
	"github.com/avasapollo/eth-playground/tools/convert"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()
	lgr := logrus.New().WithFields(logrus.Fields{
		"app": "06-contract-deploy",
	})

	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatal("can't create the client")
	}

	b, err := os.ReadFile(c.ContractOwner)
	if err != nil {
		lgr.WithError(err).Fatal("can't get wallet json file")
	}

	key, err := keystore.DecryptKey(b, c.KeyStoreAccountPassword)
	if err != nil {
		lgr.WithError(err).Fatal("can't decrypt the key")
	}

	add := crypto.PubkeyToAddress(key.PrivateKey.PublicKey)
	lgr.Infof("this is the wallet address: %s", crypto.PubkeyToAddress(key.PrivateKey.PublicKey).Hex())

	nonce, err := client.PendingNonceAt(ctx, add)
	if err != nil {
		lgr.WithError(err).Fatal("can't get sender nonce")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get suggested gas price")
	}

	// fetch the chainID, we need when we signed the transaction
	chanID, err := client.NetworkID(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get chanID")
	}

	auth, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chanID)
	if err != nil {
		lgr.WithError(err).Fatal("can't create auth request to deploy contract")
	}
	auth.GasPrice = gasPrice
	auth.GasLimit = 3000000
	auth.Nonce = big.NewInt(int64(nonce))

	conAddr, tr, _, err := todo.DeployTodo(auth, client)
	if err != nil {
		lgr.WithError(err).Fatal("can't create auth request to deploy contract")
	}
	// contract created address 0x8Ba1ded04D8170EaF6A8a2109E63B4bD85124617

	trCost, _ := convert.FromWei(convert.Ether, tr.Cost())
	lgr.WithFields(logrus.Fields{
		"tr_price":         trCost,
		"gas_limit":        3000000,
		"tr_hash":          tr.Hash().Hex(),
		"contract_address": conAddr.Hex(),
		"chan_id":          chanID,
	}).Info("contract deployed")
}
