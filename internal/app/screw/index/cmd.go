package index

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"

	"github.com/muniere/screw/internal/pkg/spider"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "index",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(cmd, args)
		},
	}

	assemble(cmd)

	return cmd
}

func run(cmd *cobra.Command, args []string) error {
	ctx, err := parse(args, cmd.Flags())
	if err != nil {
		return err
	}

	if err := prepare(ctx.Options); err != nil {
		return err
	}

	if err := perform(ctx.Uris, ctx.Options); err != nil {
		return err
	}

	return nil
}

func prepare(options *options) error {
	if options.Verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: false,
		FullTimestamp:    true,
		TimestampFormat:  "15:04:05.000",
	})

	return nil
}

func perform(uris []*url.URL, options *options) error {
	if len(uris) == 0 {
		return errors.New("no uris to index")
	}

	x := uris[len(uris)-1]

	for _, u := range uris {
		uris, err := spider.Index(u, options.IndexOptions)
		if err != nil {
			return err
		}

		for _, v := range uris {
			fmt.Println(v)
		}
		if len(uris) > 1 && u != x {
			fmt.Println()
		}
	}

	return nil
}
