package conf

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

// parseJSON parses a JSON file with the given name into a map of key-value
// strings. It fail if the keys are not strings or the file has more than
// one level of nesting.
func parseJSON(file *string) (map[string]string, error) {
	var config map[string]string

	if file == nil || *file == "" {
		return nil, nil
	}

	content, err := ioutil.ReadFile(*file)
	if err != nil {
		return nil, fmt.Errorf("error reading JSON file: %w", err)
	}

	err = json.Unmarshal(content, &config)
	if err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return nil, fmt.Errorf("json: syntax error at offset %d: %w", syntaxErr.Offset, err)
		}

		var typeErr *json.UnmarshalTypeError
		if errors.As(err, &typeErr) {
			return nil, fmt.Errorf("json: type error at offset %d: %w", typeErr.Offset, err)
		}

		return nil, fmt.Errorf("json: %w", err)
	}

	return config, nil
}
