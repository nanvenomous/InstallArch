/*
Copyright © 2025 nanvenomous mrgarelli@gmail.com
*/
package main

import (
	_ "embed"

	"github.com/nanvenomous/InstallArch/cmd"
)

//go:embed version
var version string

func main() {
	cmd.Execute(version)
}
