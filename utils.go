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
