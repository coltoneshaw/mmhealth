/*
Copyright © 2024 Colton Shaw
*/
package main

import (
	"embed"

	"github.com/coltoneshaw/mmhealth/mmhealth/cmd"
)

//go:embed checks.yaml config.yaml
var f embed.FS

func main() {
	cmd.Execute()
}
