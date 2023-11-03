package main

import (
	"github.com/abh1sheke/utrooper/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		return
	}
}
