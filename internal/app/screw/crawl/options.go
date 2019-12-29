package crawl

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
	spider.DownloadOptions
	Verbose bool
}

func assemble(cmd *cobra.Command) {
	cmd.Flags().String("focus", "href", "Focus to crawl [href/href-text/href-image/image/script]")
	cmd.Flags().String("grep", "", "Grep contents by URI with regex")
	cmd.Flags().StringP("prefix", "p", "", "Directory to download files")
	cmd.Flags().Int("concurrency", 1, "Concurrency of crawling")
	cmd.Flags().Bool("overwrite", false, "Overwrite existing files")
	cmd.Flags().BoolP("dry-run", "n", false, "Do not execute actually")
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
	if focus == spider.None {
		focus = spider.Href
	}

	grep, err := flags.GetString("grep")
	if err != nil {
		return nil, err
	}

	re, err := regexp.CompilePOSIX(grep)
	if err != nil {
		return nil, err
	}

	prefix, err := flags.GetString("prefix")
	if err != nil {
		return nil, err
	}

	concurrency, err := flags.GetInt("concurrency")
	if err != nil {
		return nil, err
	}

	overwrite, err := flags.GetBool("overwrite")
	if err != nil {
		return nil, err
	}

	dryRun, err := flags.GetBool("dry-run")
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
		DownloadOptions: spider.DownloadOptions{
			Prefix:      prefix,
			Concurrency: concurrency,
			Blocking:    false,
			Overwrite:   overwrite,
			DryRun:      dryRun,
		},
		Verbose: verbose,
	}
	return opts, nil
}
