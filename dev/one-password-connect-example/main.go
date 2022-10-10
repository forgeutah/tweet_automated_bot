package main

import (
	"fmt"

	"github.com/1Password/connect-sdk-go/connect"
)

func main() {
	// this requires OP_CONNECT_HOST OP_CONNECT_TOKEN to be send in the environment?
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		panic(err)
	}
	fmt.Println(client.GetVaults())
}
