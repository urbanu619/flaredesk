package main

import (
	"fmt"
	"go_server/base/core"
)

func main() {
	sign, err := core.BuildSignMessage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(sign)
}
