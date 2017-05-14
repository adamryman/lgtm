package main

import (
	"fmt"
	"os"

	"github.com/StudentRND/lgtm"
)

import _ "github.com/joho/godotenv/autoload"

func main() {
	if err := lgtm.Start(lgtm.DefaultConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
