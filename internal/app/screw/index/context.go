package index

import (
	"net/url"

	"github.com/spf13/pflag"
)

type context struct {
	Uris    []*url.URL
	Options *options
}

func parse(args []string, flags *pflag.FlagSet) (*context, error) {
	uris, err := normalize(args)
	if err != nil {
		return nil, err
	}

	options, err := decode(flags)
	if err != nil {
		return nil, err
	}

	ctx := &context{
		Uris:    uris,
		Options: options,
	}
	return ctx, nil
}
