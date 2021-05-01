package cmd

import (
	"fmt"
)

func printError(err error) int {
	fmt.Println(err)

	return 1
}
