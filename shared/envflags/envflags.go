package envflags

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
)

// LookupEnvFunc is a function type that looks up an environment variable by its key.
type LookupEnvFunc func(key string) (string, bool)

// FlagToEnvKeyFunc is a function type that converts a flag to its corresponding environment variable key.
type FlagToEnvKeyFunc func(f *flag.Flag) (string, error)

// FlagEnvMap is a map containing flag names and their corresponding environment variable values.
type FlagEnvMap map[string]string

// ParseOptions contains options for parsing command line arguments and environment variables.
type ParseOptions struct {
	FlagSet      *flag.FlagSet
	LookupEnv    LookupEnvFunc
	FlagToEnvKey FlagToEnvKeyFunc
}

// NewParseOptions creates a new ParseOptions instance with the provided flag set and environment variable lookup function.
// If the flag set is nil, it defaults to flag.CommandLine.
func NewParseOptions(fs *flag.FlagSet) *ParseOptions {
	if fs == nil {
		fs = flag.CommandLine
	}

	return &ParseOptions{
		FlagSet:      fs,
		LookupEnv:    os.LookupEnv,
		FlagToEnvKey: defaultToEnvKey,
	}
}

// Parse parses the command line arguments and merges them with environment variables.
// It returns the remaining arguments after parsing.
// If the options parameter is nil, it uses default options with the flag.CommandLine flag set.
func Parse(args []string, options *ParseOptions) ([]string, error) {

	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}

	if options == nil {
		options = NewParseOptions(nil)
	}

	feMap, err := flagsToEnvMap(options.FlagSet, options.LookupEnv, options.FlagToEnvKey)

	if err != nil {
		err = fmt.Errorf("error creating flag to env map: %w", err)
		return nil, err
	}

	args, err = merge(args, feMap)
	if err != nil {
		err = fmt.Errorf("error merging args and env: %w", err)
		return nil, err
	}

	err = options.FlagSet.Parse(args)

	if err != nil {
		err = fmt.Errorf("error parsing flags: %w", err)
		return nil, err
	}

	return options.FlagSet.Args(), nil
}

// Returns a new array of arguments merged with the provided args.
// Only key starting with - or -- are considered to be merged.
// Keys already present in the args will be skipped.
func merge(args []string, argsMap FlagEnvMap) ([]string, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}

	if argsMap == nil {
		return nil, fmt.Errorf("argsMap is nil")
	}

	nn, err := argsToFlags(args)

	if err != nil {
		return nil, fmt.Errorf("error normalizing args: %w", err)
	}

	for k, v := range argsMap {
		// Check if the key is already present in the args
		found := slices.Contains(nn, k)
		if !found {
			// If not, append the key and value to the args
			args = append(args, fmt.Sprintf("--%s=%s", k, v))
		}
	}

	return args, nil
}

// Returns a new array containing all the arguments starting with - or --.
func argsToFlags(args []string) ([]string, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}

	var normalized []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Trim the - or -- prefix
			arg = strings.TrimLeft(arg, "-")
			// Split the argument by '='
			parts := strings.SplitN(arg, "=", 2)
			// Add the argument to the normalized list
			normalized = append(normalized, parts[0])
		}
	}

	return normalized, nil
}

// Creates a map of arguments that are defined in the flag set and have a value set in the environment.
func flagsToEnvMap(fs *flag.FlagSet, getter LookupEnvFunc, toEnvKey FlagToEnvKeyFunc) (FlagEnvMap, error) {
	if getter == nil {
		getter = os.LookupEnv
	}

	if toEnvKey == nil {
		toEnvKey = defaultToEnvKey
	}

	args := make(FlagEnvMap)
	var loop_err error = nil

	fs.VisitAll(func(f *flag.Flag) {

		fk, err := toEnvKey(f)

		if err != nil {
			loop_err = fmt.Errorf("error converting flag to env key: %w", err)
			return
		}

		v, ok := getter(fk)

		if !ok {
			// Not set in the environment
			return
		}

		if f.DefValue == v {
			// Default value, skip
			return
		}

		args[f.Name] = v

	})

	if loop_err != nil {
		return nil, loop_err
	}

	return args, nil
}

func defaultToEnvKey(f *flag.Flag) (string, error) {
	if f == nil {
		return "", fmt.Errorf("flag is nil")
	}

	if f.Name == "" {
		return "", fmt.Errorf("flag name is empty")
	}

	// Convert the flag name to upper case and replace "-" with "_"
	envKey := strings.ToUpper(f.Name)
	envKey = strings.ReplaceAll(envKey, "-", "_")

	return envKey, nil
}
