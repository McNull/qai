package app

import (
	"github.com/mcnull/qai/shared/envflags"
	"github.com/mcnull/qai/shared/provider"
	"strings"
)

func parseFlags(args []string, exitOnError bool, values *provider.FlagValues) (*provider.FlagValues, error) {

	fs := provider.CreateFlagSet(APP_NAME, values, exitOnError)

	options := envflags.NewParseOptions(fs)

	// Parse the command line arguments and merge them with environment variables
	remainingArgs, err := envflags.Parse(args, options)

	if err != nil {
		return nil, err
	}

	values.Prompt = strings.Join(remainingArgs, " ")

	return values, nil
}
