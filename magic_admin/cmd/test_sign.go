package main

import (
	"fmt"
	"go_server/utils"
)

func main() {
	sign, err := utils.BuildSignMessage()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(sign)
}
