package main

import (
	"fmt"

	"github.com/1Password/connect-sdk-go/connect"
)

func main() {
	client, err := connect.NewClientFromEnvironment()
	if err != nil {
		panic(err)
	}
	fmt.Println(client.GetVaults())
}
