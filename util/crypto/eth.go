package crypto

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func chainMessage() {

	client, err := ethclient.Dial("http://18.117.169.91:8545")
	if err != nil {
		logrus.Fatal(err)
	}

	bn, err := client.BlockNumber(context.Background())
	if err != nil {
		logrus.Fatal(err)
	}
	fmt.Print(bn)
}
