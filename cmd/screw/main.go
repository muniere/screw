package main

import (
	"github.com/muniere/screw/internal/app/screw"
	"github.com/muniere/screw/internal/pkg/sys"
)

func main() {
	cmd := screw.NewCommand()
	err := cmd.Execute()
	sys.CheckError(err)
}
