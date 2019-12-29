package screw

import (
	"github.com/spf13/cobra"

	"github.com/muniere/screw/internal/app/screw/crawl"
	"github.com/muniere/screw/internal/app/screw/index"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "screw",
	}

	cmd.AddCommand(index.NewCommand())
	cmd.AddCommand(crawl.NewCommand())

	return cmd
}
