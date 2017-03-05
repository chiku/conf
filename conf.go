// Package conf is for extracting application configuration.
// It uses configuration from command-line arguments, JSON file,
// environment variable or default value.
//
// example.go
//
//  import (
//      "fmt"
//
//      "github.com/chiku/conf"
//  )
//
//  func main() {
//      options := map[string]conf.Option{
//          "foo": conf.Option{
//              Desc:      "a description for foo",
//              Default:   "default foo",
//              Mandatory: true,
//          }
//          "bar": conf.Option{Mandatory: true},
//          "baz": conf.Option{Desc: "a description for baz"},
//          "qux": conf.Option{},
//      }
//
//      loader := conf.MultiLoader{
//          JSONKey: "shr",
//          Options: options,
//          Usage:   "Example application",
//      }
//
//      config, origin, err := loader.Load()
//
//      if err != nil {
//          fmt.Printf("error: %s\n", err)
//          return
//      }
//
//      fmt.Printf("configuration: %#v\n", config)
//      fmt.Printf("origin: %#v\n", origin)
//  }
//
// Usage
//
//     go build -o example example.go
//     ./example -foo fooval -bar barval -shr file.json
//
package conf

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

// A Loader represents a configuration loader.
type Loader interface {
	// Load extracts configuration from different sources. It returns the
	// configuration and their origin, and an error if present.
	Load() (config map[string]string, origin map[string]string, err error)
}

// An Option represents a configuration for github.com/chiku/conf
type Option struct {
	// Default is the value used if not provided in command-line flag,
	// JSON file or environment variable.
	Default string

	// Desc is the command line argument description.
	Desc string

	// Mandatory is true if the configuration must be specified.
	Mandatory bool
}

// MultiLoader is a configuration loader with different sources.
// It extracts values from command-line arguments, JSON configuration file,
// environment variable and a fallback default value.
type MultiLoader struct {
	// Options is a map of Option for a given configuration key. The
	// configuration and origin returned by Load() use the same keys.
	Options map[string]Option

	// JSONKey, if not empty, is the configuration key name expected
	// for the JSON configuration file.
	JSONKey string

	// Usage is a description for the application. Usage shows up when
	// the application is run with "-help".
	Usage string
}

// Load extracts configuration from different sources. It returns the
// configuration and their origin, and an error if present.
// The configurations are loaded in following order.
//   1. Command-line arguments
//   2. JSON file mentioned in JSONKey
//   3. Environment variable
//   4. Default values.
//
// The origin is returned as a string and can be one of "Flags", "JSON",
// "Environment" or "Defaults"
// based on what was matched when looking up for the configuration.
// The configuration is always returned as a map[string]string.
// Load() returns an error in the following cases.
//   1. Command-line argument parse fails.
//   2. JSON parse fails.
//   3. Mandatory configuration was not provided.
func (l MultiLoader) Load() (config map[string]string, origin map[string]string, err error) {
	program, args := os.Args[0], os.Args[1:]
	flagsHandler := func(flags *flag.FlagSet) {
		flags.Usage = func() {
			fmt.Printf("%s: %s\n\nParameters:\n", program, l.Usage)
			flags.PrintDefaults()
			os.Exit(0)
		}
	}

	return l.load(args, flagsHandler)
}

// load extracts configuration from different sources. It returns the
// configuration and their origin, and an error if present.
func (l MultiLoader) load(args []string, flagsHandler func(flags *flag.FlagSet)) (config map[string]string, origin map[string]string, err error) {
	config = make(map[string]string)
	origin = make(map[string]string)

	flagVals, err := l.parseFlags(args, flagsHandler)
	if err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	jsonFile := flagVals[l.JSONKey]
	jsonConfig, err := parseJSON(jsonFile)
	if err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	l.configure(config, origin, func(key string) string { return *flagVals[key] }, "Flags")
	l.configure(config, origin, func(key string) string { return jsonConfig[key] }, "JSON")
	l.configure(config, origin, func(key string) string { return os.Getenv(key) }, "Environment")
	l.configure(config, origin, func(key string) string { return l.Options[key].Default }, "Defaults")

	if err = l.verifyMandatoryPresent(config); err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	return config, origin, nil
}

// parseFlags parses application-level command-line flags. The flags
// are based on the configuration value and JSON-key flag. It returns
// the parsed values as a map of string to pointer of strings and an
// error if parse fails.
func (l MultiLoader) parseFlags(args []string, flagsHandler func(*flag.FlagSet)) (flagVals map[string]*string, err error) {
	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flagsHandler(flags)

	flagVals = make(map[string]*string)
	for name, option := range l.Options {
		if desc := option.Desc; desc != "" {
			flagVals[name] = flags.String(name, "", desc)
		} else {
			flagVals[name] = flags.String(name, "", name)
		}
	}

	if l.JSONKey != "" {
		flagVals[l.JSONKey] = flags.String(l.JSONKey, "", "JSON configuration file")
	}

	err = flags.Parse(args)
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %s", err)
	}

	return flagVals, nil
}

// A mappingFunc on running returns a value against a key. MappingFuncs
// are processed by configure.
type mappingFunc func(key string) (value string)

// Configure adds value and origin against a key if not already present.
func (l MultiLoader) configure(config map[string]string, origin map[string]string, mapping mappingFunc, from string) {
	for name := range l.Options {
		if config[name] == "" {
			config[name] = mapping(name)
			origin[name] = from
		}
	}
}

// VerifyMandatoryPresent returns an error if one or more mandatory
// parameters are missing. The error message reports all the missing
// configuration keys.
func (l MultiLoader) verifyMandatoryPresent(config map[string]string) error {
	var missing []string
	for name, option := range l.Options {
		if config[name] == "" && option.Mandatory {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		sort.Strings(missing)
		return fmt.Errorf("missing mandatory configurations: %s", strings.Join(missing, ", "))
	}

	return nil
}
