package main

import (
	"context"

	"github.com/avasapollo/eth-playground/config"
	"github.com/avasapollo/eth-playground/tools/convert"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()

	lgr := logrus.New().WithFields(logrus.Fields{
		"app": "02-balance",
	})

	client, err := ethclient.DialContext(context.Background(), c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatalf("can't get the client")
	}
	defer client.Close()

	// get the last block of the chain
	account := common.HexToAddress(c.MyAccount)
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		lgr.WithError(err).Fatal("can't get the balance of the last block")
	}
	// 1 ether = 10^18 wei

	valueEth, _ := convert.FromWei(convert.Ether, balance)
	lgr.WithField("balance", valueEth).Info("balance")
}
