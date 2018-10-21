package main

import (
	"fmt"
	"os"
	"simple-chatroom-in-GoLang-RPC/help"
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		fmt.Println(help.HowToRun())
		os.Exit(1)
	}
}
