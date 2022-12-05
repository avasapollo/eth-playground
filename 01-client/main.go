package main

import (
	"context"

	"github.com/avasapollo/eth-playground/config"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	c := config.Get()

	lgr := logrus.New().WithFields(logrus.Fields{
		"app": "01-client",
	})

	client, err := ethclient.DialContext(context.Background(), c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatalf("can't get the client")
	}
	defer client.Close()

	// get the last block of the chain
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		lgr.WithError(err).Error("can't get the last block")
	}
	lgr.WithFields(logrus.Fields{
		"block_number":  block.Number(),
		"block_receipt": block.Header().ReceiptHash.String(),
	}).Info("last block")
}
