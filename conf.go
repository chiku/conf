package conf

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
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

	var jsonConfig map[string]string
	if jsonFile != nil {
		jsonContent, _ := ioutil.ReadFile(*jsonFile)
		_ = json.Unmarshal(jsonContent, &jsonConfig)
	}

	l.configure(config, origin, func(key string) string { return *flagVals[key] }, "Flags")
	l.configure(config, origin, func(key string) string { return jsonConfig[key] }, "JSON")
	l.configure(config, origin, func(key string) string { return os.Getenv(key) }, "Environment")
	l.configure(config, origin, func(key string) string { return l.Defaults[key] }, "Defaults")

	return config, origin, err
}

type mappingFunc func(key string) (value string)

func (l MultiLoader) configure(config map[string]string, origin map[string]string, mapping mappingFunc, from string) {
	for _, item := range l.Mandatory {
		if config[item] == "" {
			config[item] = mapping(item)
			origin[item] = from
		}
	}
	for _, item := range l.Optional {
		if config[item] == "" {
			config[item] = mapping(item)
			origin[item] = from
		}
	}
}
