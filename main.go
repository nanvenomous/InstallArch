/*
Copyright © 2025 nanvenomous mrgarelli@gmail.com
*/
package main

import (
	"embed"

	"github.com/nanvenomous/InstallArch/cmd"
)

//go:embed version
var version string

//go:embed rsrc
var rsrc embed.FS

func main() {
	cmd.Execute(version, rsrc)
}
