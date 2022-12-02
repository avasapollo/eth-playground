package main

import (
	"context"
	"math/big"

	"github.com/avasapollo/eth-playground/config"
	"github.com/avasapollo/eth-playground/tools/convert"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func main() {
	lgr := logrus.New().WithField("app", "05-recurring")
	c := config.Get()

	// generate keystore
	keyStorage := keystore.NewKeyStore(c.KeyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)

	// receiver balance
	receiverAcc, err := keyStorage.NewAccount(c.KeyStoreAccountPassword)
	if err != nil {
		lgr.WithError(err).Fatal("can't create the account")
	}

	lgr = lgr.WithField("receiver", receiverAcc.Address.Hex())
	lgr.Info("receiver created")
	ctx := context.Background()
	client, err := ethclient.DialContext(ctx, c.NetworkURL)
	if err != nil {
		lgr.WithError(err).Fatal("can't create the client")
	}

	bReceiver, err := client.BalanceAt(ctx, receiverAcc.Address, nil)
	if err != nil {
		lgr.WithError(err).Fatal("can't get sender balance")
	}

	bEth, _ := convert.FromWei(convert.Ether, bReceiver)
	lgr.Infof("receiver balance %v", bEth)

	// sender balance
	senderAdd := common.HexToAddress(c.MyAccount)
	lgr = lgr.WithField("sender", senderAdd.Hex())
	bSender, err := client.BalanceAt(ctx, senderAdd, nil)
	if err != nil {
		lgr.WithError(err).Fatal("can't get sender balance")
	}

	bEth, _ = convert.FromWei(convert.Ether, bSender)
	lgr.Infof("sender balance %v", bEth)

	// create transactions
	// get nonce
	nonce, err := client.PendingNonceAt(ctx, senderAdd)
	if err != nil {
		lgr.WithError(err).Fatal("can't get sender nonce")
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get suggested gas price")
	}
	lgr = lgr.WithField("gas_price_suggested", gasPrice)

	// 1 ether = 1000000000000000000
	lgr.WithField("nonce", nonce)
	amount := big.NewInt(100000000000000000)
	txData := &types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      21000,
		To:       &receiverAcc.Address,
		Value:    amount,
	}

	// fetch the chainID, we need when we signed the transaction
	chanID, err := client.NetworkID(ctx)
	if err != nil {
		lgr.WithError(err).Fatal("can't get chanID")
	}

	lgr = lgr.WithField("chan_id", chanID)

	// generate the private key from the raw private key that is inside the MY_ACCOUNT_PRIVATE_KEY
	prKey := convert.ToPrivateKey(c.MyAccountPrivateKey)

	tr, err := types.SignNewTx(prKey, types.NewEIP155Signer(chanID), txData)
	if err != nil {
		lgr.WithError(err).Fatal("can't sign the transaction")
	}

	err = client.SendTransaction(ctx, tr)
	if err != nil {
		lgr.WithError(err).Fatal("can't send the transaction")
	}
	trCost, _ := convert.FromWei(convert.Ether, tr.Cost())
	lgr = lgr.WithFields(logrus.Fields{
		"tr_hash":  tr.Hash().Hex(),
		"tr_nonce": tr.Nonce(),
		"tr_cost":  trCost,
	})

	lgr.Info("transaction sent!!!")
}
