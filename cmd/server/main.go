package main

import (
	"fmt"
	"os"

	// pb "getitqec.com/server/auth/pkg/api/v1"

	cmd "getitqec.com/server/file/pkg/cmd"
)

func main() {
	fmt.Printf("Starting file service...\n")
	if err := cmd.RunServer(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
