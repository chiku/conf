package conf

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

type MultiLoader struct {
	JSONKey     string
	Mandatory   []string
	Optional    []string
	Description map[string]string
	Defaults    map[string]string
}

func (l MultiLoader) Load() (config map[string]string, origin map[string]string, err error) {
	args := os.Args[1:]
	return l.load(args)
}

func (l MultiLoader) load(args []string) (config map[string]string, origin map[string]string, err error) {
	config = make(map[string]string)
	origin = make(map[string]string)

	if err = l.verifyPresence(); err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	if err = l.verifyUniqueness(); err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	if err = l.verifyDescriptions(); err != nil {
		return nil, nil, fmt.Errorf("conf.Load: %s", err)
	}

	flagVals, err := l.parseFlags(args)
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

func (l MultiLoader) verifyPresence() error {
	var missing []string

	if isPresentInside(l.Mandatory, "") {
		missing = append(missing, "mandatory")
	}
	if isPresentInside(l.Optional, "") {
		missing = append(missing, "optional")
	}

	if len(missing) > 0 {
		return fmt.Errorf("empty keys exist: %s", strings.Join(missing, ", "))
	}

	return nil
}

func uniqueness(items, existingDuplMsgs []string, key string) (uniq, duplMsgs []string) {
	uniq, dupl := partitionByUniqueness(items)

	if len(dupl) > 0 {
		return uniq, append(existingDuplMsgs, fmt.Sprintf("%s(%s)", key, strings.Join(dupl, ", ")))
	}

	return uniq, existingDuplMsgs
}

func (l MultiLoader) verifyUniqueness() error {
	var dulpMsgs []string

	uniqMandatory, dulpMsgs := uniqueness(l.Mandatory, dulpMsgs, "mandatory")
	uniqOptional, dulpMsgs := uniqueness(l.Optional, dulpMsgs, "optional")
	_, dulpMsgs = uniqueness(append(uniqMandatory, uniqOptional...), dulpMsgs, "mandatory+optional")
	_, dulpMsgs = uniqueness(append(uniqMandatory, l.JSONKey), dulpMsgs, "mandatory+jsonkey")
	_, dulpMsgs = uniqueness(append(uniqOptional, l.JSONKey), dulpMsgs, "optional+jsonkey")

	if len(dulpMsgs) > 0 {
		return fmt.Errorf("configuration keys are duplicated: %s", strings.Join(dulpMsgs, ", "))
	}

	return nil
}

func (l MultiLoader) verifyDescriptions() error {
	descs := keysIn(l.Description)
	knowns := append(l.Mandatory, l.Optional...)

	extras := extraItems(descs, knowns)
	if len(extras) == 0 {
		return nil
	}

	return fmt.Errorf("description keys are unknown: %s", strings.Join(extras, ", "))
}

func (l MultiLoader) description(flags *flag.FlagSet, item string) *string {
	if desc, ok := l.Description[item]; ok {
		return flags.String(item, "", desc)
	}
	return flags.String(item, "", item)
}

func (l MultiLoader) parseFlags(args []string) (flagVals map[string]*string, err error) {
	flagVals = make(map[string]*string)
	flags := flag.NewFlagSet("", flag.ContinueOnError)

	for _, item := range l.Mandatory {
		flagVals[item] = l.description(flags, item)
	}
	for _, item := range l.Optional {
		flagVals[item] = l.description(flags, item)
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
