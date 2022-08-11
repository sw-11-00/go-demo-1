package main

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func main() {
	TestAccAddressFromBech32("cosmos19lr8j0agk45hxmxfz2s3hn97hlgkfvvwt0ydfc")
}

func GenAddress() {
	pk := ed25519.GenPrivKey().PubKey()
	addr := pk.Address()
	fmt.Println(sdk.AccAddress(addr).String())
}

func TestAccAddressFromBech32(address string) {
	fmt.Println(sdk.AccAddressFromBech32(address))
}
