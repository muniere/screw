package index

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/muniere/screw/internal/pkg/spider"
)

type options struct {
	spider.IndexOptions
	Verbose bool
}

func assemble(cmd *cobra.Command) {
	cmd.Flags().String("focus", "href", "Focus to crawl [href/href-text/href-image/image/script]")
	cmd.Flags().String("grep", "", "Grep contents by URI with regex")
	cmd.Flags().BoolP("verbose", "v", false, "Show verbose log messages")
}

func decode(flags *pflag.FlagSet) (*options, error) {
	focus_, err := flags.GetString("focus")
	if err != nil {
		return nil, err
	}

	focus := spider.NewFocus(focus_)
	if focus == spider.None {
		return nil, errors.New(fmt.Sprintf("unsupported type of focus: %s", focus_))
	}

	grep, err := flags.GetString("grep")
	if err != nil {
		return nil, err
	}

	re, err := regexp.CompilePOSIX(grep)
	if err != nil {
		return nil, err
	}

	verbose, err := flags.GetBool("verbose")
	if err != nil {
		return nil, err
	}

	opts := &options{
		IndexOptions: spider.IndexOptions{
			Focus: focus,
			Grep:  re,
		},
		Verbose: verbose,
	}
	return opts, nil
}
