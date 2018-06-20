package main

import (
	// Load all contracts here, so we can always read and decode
	// transactions with those contracts.
	_ "github.com/eoscanada/eos-go/msig"
	_ "github.com/eoscanada/eos-go/system"
	_ "github.com/eoscanada/eos-go/token"

	"github.com/eoscanada/eosc/eosc/cmd"
)

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	cmd.Execute()
}
