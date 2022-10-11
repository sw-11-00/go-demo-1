package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/machinebox/graphql"
)

type PoolCreateResp struct {
	PoolCreateEntities []PoolCreateEntity `json:"poolCreateEntities"`
}

type PoolCreateEntity struct {
	Id string `json:"id"`
}

var client *graphql.Client

func setupClient() {
	client = graphql.NewClient("http://localhost:8000/subgraphs/name/perpetual/spidex")
}

func getAllProducts() {
	query := graphql.NewRequest(`
		{
  			poolCreateEntities {
    			id
  			}
		}
	`)

	ctx := context.Background()

	var poolCreateResp PoolCreateResp

	if err := client.Run(ctx, query, &poolCreateResp); err != nil {
		log.Fatal(err)
	}

	jsonResponse, _ := json.MarshalIndent(poolCreateResp, "", "\t")

	fmt.Printf("%+v\n", string(jsonResponse))
}

func main() {
	setupClient()
	getAllProducts()
}
