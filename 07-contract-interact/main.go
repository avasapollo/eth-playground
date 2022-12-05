package main

import (
	"context"
	"math/big"
	"os"

	"github.com/avasapollo/eth-playground/config"
	todo "github.com/avasapollo/eth-playground/gen"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()
	lgr := logrus.New().WithFields(logrus.Fields{
		"app": "07-contract-interact",
	})

	ctx := context.Background()
	b, err := os.ReadFile(c.ContractOwner)
	if err != nil {
		lgr.WithError(err).Fatal("can't get the contract owner")
	}

	key, err := keystore.DecryptKey(b, "password")
	if err != nil {
		lgr.WithError(err).Fatal("can't decrypt contract key")
	}

	client, err := ethclient.DialContext(ctx, c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatal("can't get client")
	}
	defer client.Close()

	contractAdd := common.HexToAddress(c.ContractAddress)
	todoClient, err := todo.NewTodo(contractAdd, client)

	chanID, err := client.NetworkID(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get network id")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get gas price")
	}

	// Add Task to the contract
	if err != nil {
		lgr.WithError(err).Fatal("can't get todo client")
	}
	transactorOpt, err := bind.NewKeyedTransactorWithChainID(key.PrivateKey, chanID)
	if err != nil {
		lgr.WithError(err).Fatal("can't get transaction")
	}
	transactorOpt.GasLimit = 3000000
	transactorOpt.GasPrice = gasPrice

	// Add operation
	// tr, err := todoClient.Add(transactorOpt, "first task")
	//if err != nil {
	//	lgr.WithError(err).Fatal("can't get transaction")
	//}
	//lgr.WithFields(logrus.Fields{
	//	"tr_hex": tr.Hash().Hex(),
	//	"op":     "Add",
	//}).Info("added task")

	// update operation
	tr, err := todoClient.Update(transactorOpt, big.NewInt(0), "updated task content")
	if err != nil {
		lgr.WithError(err).Fatal("can't get transaction")
	}
	lgr.WithFields(logrus.Fields{
		"tr_hex": tr.Hash().Hex(),
		"op":     "Update",
	}).Info("updated task")

	// LIST Task from the contract
	//callerOpt := &bind.CallOpts{
	//	From: crypto.PubkeyToAddress(key.PrivateKey.PublicKey),
	//}
	//
	//listRes, err := todoClient.List(callerOpt)
	//if err != nil {
	//	lgr.WithError(err).Fatal("can't get transaction")
	//}
	//lgr.WithFields(logrus.Fields{
	//	"list_res": listRes,
	//	"op":       "List",
	//}).Info("list tasks")
}
