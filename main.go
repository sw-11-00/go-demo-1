package main

import (
	"context"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"

	"go-demo-1/util/math"
)

type PoolOperate struct {
	charge   []*math.Decimal
	withdraw []*math.Decimal
}

type PoolDepositResp struct {
	PoolCreateEntities []PoolDepositEntity `json:"poolDepositEntities"`
}

type PoolDepositEntity struct {
	UserAddr  string `json:"userAddr"`
	Share     string `json:"share"`
	Amount    string `json:"amount"`
	PoolAddr  string `json:"poolAddr"`
	Timestamp string `json:"timestamp"`
}

type PoolWithdrawResp struct {
	PoolWithdrawEntities []PoolWithdrawEntity `json:"poolWithdrawEntities"`
}

type PoolWithdrawEntity struct {
	UserAddr  string `json:"userAddr"`
	Share     string `json:"share"`
	Amount    string `json:"amount"`
	PoolAddr  string `json:"poolAddr"`
	Timestamp string `json:"timestamp"`
}

var client *graphql.Client

func main() {
	call, err := PoolOperateCall("0x65dcfd504c4b5078005257827744539c282cb029")
	if err != nil {
		return
	}
	fmt.Println(call)
}

func PoolOperateCall(poolAddress string) (*PoolOperate, error) {
	client = graphql.NewClient("http://localhost:8000/subgraphs/name/perpetual/spidex")
	poolDeposits, err := poolDeposit(poolAddress)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	fmt.Println(poolDeposits)

	poolWithdraws, err := poolWithdraw(poolAddress)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	fmt.Println(poolWithdraws)

	return nil, nil
}

func poolDeposit(poolAddress string) (*math.Decimal, error) {
	query := graphql.NewRequest(`
		{
  			poolDepositEntities(where: {poolAddr: "` + poolAddress + `"}) {
    			timestamp
    			share
    			userAddr
    			amount
				poolAddr
  			}
		}
	`)

	var poolDepositResp PoolDepositResp

	if err := client.Run(context.Background(), query, &poolDepositResp); err != nil {
		log.Fatal(err)
	}

	return nil, nil
}

func poolWithdraw(poolAddress string) (*math.Decimal, error) {
	query := graphql.NewRequest(`
		{
			poolWithdrawEntities(where: {poolAddr: "` + poolAddress + `"}){
    			share
    			timestamp
    			userAddr
    			amount
				poolAddr
  			}
		}
	`)

	var poolWithdrawResp PoolWithdrawResp

	if err := client.Run(context.Background(), query, &poolWithdrawResp); err != nil {
		log.Fatal(err)
	}

	return nil, nil
}
