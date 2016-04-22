package conf

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
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

	flagVals, err := l.parseFlags()
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
	l.configure(config, origin, func(key string) string { return l.Defaults[key] }, "Defaults")

	if err = l.verifyMandatoryPresent(config); err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	return config, origin, nil
}

func (l MultiLoader) parseFlags() (flagVals map[string]*string, err error) {
	flagVals = make(map[string]*string)
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	for _, item := range l.Mandatory {
		flagVals[item] = flags.String(item, "", item)
	}
	for _, item := range l.Optional {
		flagVals[item] = flags.String(item, "", item)
	}
	if l.JSONKey != "" {
		flagVals[l.JSONKey] = flags.String(l.JSONKey, "", "JSON configuration file")
	}

	err = flags.Parse(l.Args)
	if err != nil {
		return nil, fmt.Errorf("error parsing flags: %s", err)
	}

	return flagVals, nil
}

func parseJSON(file *string) (map[string]string, error) {
	var config map[string]string

	if file == nil || *file == "" {
		return nil, nil
	}

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %s", err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON file: %s", err)
	}

	return config, nil
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

func (l MultiLoader) verifyMandatoryPresent(config map[string]string) error {
	var missing []string
	for _, item := range l.Mandatory {
		if config[item] == "" {
			missing = append(missing, item)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing mandatory configurations: %s", strings.Join(missing, ", "))
	}

	return nil
}
