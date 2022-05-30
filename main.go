package main

import (
	"fmt"
	"os"

	"github.com/Clink-n-Clank/Eitri/cmd/eitri"
)

func main() {
	if execErr := eitri.ExecCommandTool(); execErr != nil {
		fmt.Println(execErr.Error())
		os.Exit(1)
	}
}
