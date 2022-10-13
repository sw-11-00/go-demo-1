package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/machinebox/graphql"
	"github.com/sirupsen/logrus"

	_ "go-demo-1/util/math"
)

var timeDeltaTmp = strconv.FormatInt(time.Now().Unix()-60*60, 10)

type PoolOperate struct {
	charges   []string
	withdraws []string
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
	poolDeposits, err := poolDeposit(poolAddress, timeDeltaTmp)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var charges []string
	for _, val := range poolDeposits.PoolCreateEntities {
		charges = append(charges, val.Amount)
	}

	poolWithdraws, err := poolWithdraw(poolAddress, timeDeltaTmp)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	var withdraws []string
	for _, val := range poolWithdraws.PoolWithdrawEntities {
		withdraws = append(withdraws, val.Amount)
	}
	return &PoolOperate{
		charges,
		withdraws,
	}, nil
}

func poolDeposit(poolAddress string, timeDelta string) (*PoolDepositResp, error) {
	query := graphql.NewRequest(`
		{
  			poolDepositEntities(where: {poolAddr: "` + poolAddress + `", timestamp_gt: "` + timeDelta + `"}) {
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

	return &poolDepositResp, nil
}

func poolWithdraw(poolAddress string, timeDelta string) (*PoolWithdrawResp, error) {
	query := graphql.NewRequest(`
		{
			poolWithdrawEntities(where: {poolAddr: "` + poolAddress + `", timestamp_gt: "` + timeDelta + `"}){
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

	return &poolWithdrawResp, nil
}
