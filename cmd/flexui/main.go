// Command flexui is the FlexUI 2.0 compiler and command line interface.
package main

import (
	"os"

	"github.com/Ahaduzzamankhan/flexui/internal/cli"
)

func main() {
	os.Exit(cli.Main(os.Args[1:], os.Stdout, os.Stderr))
}
