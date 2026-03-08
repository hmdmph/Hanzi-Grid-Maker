//go:build !gui

package main

import "fmt"

func launchGUI() {
	fmt.Println("Native GUI is not available in this build.")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  1. Use  -web  flag to launch browser-based UI")
	fmt.Println("  2. Rebuild with GUI support:")
	fmt.Println("       CGO_ENABLED=1 go build -tags gui -o print-square-gui .")
	fmt.Println("       make build-gui")
}
