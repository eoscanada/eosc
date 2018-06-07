package main

import "github.com/eoscanada/eosc/eosc/cmd"

var version = "dev"

func init() {
	cmd.Version = version
}

func main() {
	cmd.Execute()
}
