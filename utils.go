package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func isPresentInside(list []string, key string) bool {
	for _, item := range list {
		if item == key {
			return true
		}
	}

	return false
}

func partitionByUniqueness(list []string) (uniq, dupl []string) {
	presentInUniq := make(map[string]bool)
	presentInDupl := make(map[string]bool)

	for _, item := range list {

		isPresentInUniq := presentInUniq[item]
		isPresentInDupl := presentInDupl[item]

		if !isPresentInUniq && !isPresentInDupl {
			uniq = append(uniq, item)
			presentInUniq[item] = true
		}

		if isPresentInUniq && !isPresentInDupl {
			dupl = append(dupl, item)
			presentInDupl[item] = true
		}
	}

	var truelyUniq []string

	for _, item := range uniq {
		if !presentInDupl[item] {
			truelyUniq = append(truelyUniq, item)
		}
	}

	uniq = truelyUniq

	return uniq, dupl
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
		details := ""
		if serr, ok := err.(*json.SyntaxError); ok {
			details = fmt.Sprintf(" (syntax error at offset: %d)", serr.Offset)
		}

		if serr, ok := err.(*json.UnmarshalTypeError); ok {
			details = fmt.Sprintf(" (type error at offset: %d)", serr.Offset)
		}

		return nil, fmt.Errorf("error parsing JSON file: %s%s", err, details)
	}

	return config, nil
}
