package conf

// utils.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2016-2017. All rights reserved
// License::   MIT

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// parseJSON parses a JSON file with the given name into a map of key-value strings.
// It fail if the keys are not strings or the file has more than one level of nesting.
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
		if serr, ok := err.(*json.SyntaxError); ok {
			return nil, fmt.Errorf("json: syntax error at offset %d: %s", serr.Offset, err)
		}

		if terr, ok := err.(*json.UnmarshalTypeError); ok {
			return nil, fmt.Errorf("json: type error at offset %d: %s", terr.Offset, err)
		}

		return nil, fmt.Errorf("json: %s", err)
	}

	return config, nil
}
