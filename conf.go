// conf.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2016-2017. All rights reserved
// License::   MIT

package conf

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
)

type MultiLoader struct {
	Options map[string]Option
	JSONKey string
	Usage   string
}

func (l MultiLoader) Load() (config map[string]string, origin map[string]string, err error) {
	program, args := os.Args[0], os.Args[1:]
	flagsHandler := func(flags *flag.FlagSet) {
		flags.Usage = func() {
			fmt.Fprintf(os.Stderr, "%s: %s\n\nParameters:\n", program, l.Usage)
			flags.PrintDefaults()
			os.Exit(0)
		}
	}

	return l.load(args, flagsHandler)
}

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

type mappingFunc func(key string) (value string)

func (l MultiLoader) configure(config map[string]string, origin map[string]string, mapping mappingFunc, from string) {
	for name, _ := range l.Options {
		if config[name] == "" {
			config[name] = mapping(name)
			origin[name] = from
		}
	}
}

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
