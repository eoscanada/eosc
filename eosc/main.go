package main

import "github.com/eoscanada/eosc/eosc/cmd"

var version = "dev"

func main() {
	cmd.Version = version
	cmd.Execute()
}
