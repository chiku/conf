package conf

import (
	"encoding/json"
	"flag"
	"os"
)

const (
	JSON     = "JSON"
	Flag     = "Flags"
	Env      = "Environment"
	Defaults = "Defaults"
)

type MultiLoader struct {
	JSON      string
	Mandatory []string
	Optional  []string
	Defaults  map[string]string
	Args      []string
}

func (l MultiLoader) Load() (config map[string]string, origin map[string]string, err error) {
	config = make(map[string]string)
	origin = make(map[string]string)

	_ = json.Unmarshal([]byte(l.JSON), &config)
	for _, item := range l.Mandatory {
		origin[item] = JSON
	}
	for _, item := range l.Optional {
		origin[item] = JSON
	}

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flagVals := make(map[string]*string)
	for _, item := range l.Mandatory {
		flagVals[item] = flags.String(item, "", item)
	}
	for _, item := range l.Optional {
		flagVals[item] = flags.String(item, "", item)
	}
	_ = flags.Parse(l.Args)
	for _, item := range l.Mandatory {
		if config[item] == "" {
			config[item] = *flagVals[item]
			origin[item] = Flag
		}
	}
	for _, item := range l.Optional {
		if config[item] == "" {
			config[item] = *flagVals[item]
			origin[item] = Flag
		}
	}

	for _, item := range l.Mandatory {
		if config[item] == "" {
			config[item] = os.Getenv(item)
			origin[item] = Env
		}
	}
	for _, item := range l.Optional {
		if config[item] == "" {
			config[item] = os.Getenv(item)
			origin[item] = Env
		}
	}

	for _, item := range l.Mandatory {
		if config[item] == "" {
			config[item] = l.Defaults[item]
			origin[item] = Defaults
		}
	}
	for _, item := range l.Optional {
		if config[item] == "" {
			config[item] = l.Defaults[item]
			origin[item] = Defaults
		}
	}

	return config, origin, err
}
