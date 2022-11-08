package main

import (
	"context"
	"fmt"
	"math"
	"math/big"

	"github.com/avasapollo/eth-playground/config"
	"github.com/avasapollo/eth-playground/tools"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()

	lgr := logrus.New()

	client, err := ethclient.DialContext(context.Background(), c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatalf("can't get the client")
	}
	defer client.Close()

	// get the last block of the chain
	account := common.HexToAddress(c.MyAccount)
	balanceAt, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		lgr.WithError(err).Error("can't get the balance of the last block")
	}
	print(tools.WeiToEther(balanceAt))
	fbalance := new(big.Float)
	fbalance.SetString(balanceAt.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Println(ethValue)
	b := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	lgr.WithFields(logrus.Fields{
		"balance": b,
	}).Info("balance at last block")
}
