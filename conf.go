package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
)

const (
	JSON     = "JSON"
	Flag     = "Flags"
	Env      = "Environment"
	Defaults = "Defaults"
)

type MultiLoader struct {
	JSONKey   string
	Mandatory []string
	Optional  []string
	Defaults  map[string]string
	Args      []string
}

func (l MultiLoader) Load() (config map[string]string, origin map[string]string, err error) {
	config = make(map[string]string)
	origin = make(map[string]string)

	flags := flag.NewFlagSet("", flag.ContinueOnError)
	flagVals := make(map[string]*string)
	for _, item := range l.Mandatory {
		flagVals[item] = flags.String(item, "", item)
	}
	for _, item := range l.Optional {
		flagVals[item] = flags.String(item, "", item)
	}

	var jsonFile *string
	if l.JSONKey != "" {
		jsonFile = flags.String(l.JSONKey, "", "JSON configuration file")
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

	var jsonConfig map[string]string
	if jsonFile != nil {
		jsonContent, _ := ioutil.ReadFile(*jsonFile)
		_ = json.Unmarshal(jsonContent, &jsonConfig)
	}

	for _, item := range l.Mandatory {
		if config[item] == "" {
			config[item] = jsonConfig[item]
			origin[item] = JSON
		}
	}
	for _, item := range l.Optional {
		if config[item] == "" {
			config[item] = jsonConfig[item]
			origin[item] = JSON
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
