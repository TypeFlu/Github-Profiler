package main

import (
	"fmt"
	"os"

	"github-profiler/cmd"
)

const (
	version = "1.0.0"
	author  = "github@Tyeflu"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
