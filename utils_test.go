package conf

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestParseJSON(t *testing.T) {
	jsonFile := createFile(t, `{ "foo": "abc", "bar": "xyz" }`)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	data, err := parseJSON(&jsonFile)
	if err != nil {
		t.Fatalf("Unexpected error parsing valid JSON file: %s", err)
	}

	expectedData := map[string]string{
		"foo": "abc",
		"bar": "xyz",
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Error("Invalid parsed data")
		t.Errorf("Actual:   %#v", data)
		t.Errorf("Expected: %#v", expectedData)
	}
}

func TestParseJSONWithoutFileName(t *testing.T) {
	name := ""
	data, err := parseJSON(&name)
	if err != nil {
		t.Fatalf("Unexpected error parsing empty JSON file name: %s", err)
	}
	if len(data) != 0 {
		t.Errorf("Unexpected data for empty JSON file name: %#v", data)
	}

	data, err = parseJSON(nil)
	if err != nil {
		t.Errorf("Unexpected data for nil JSON file name: %s", err)
	}
	if len(data) != 0 {
		t.Errorf("Unexpected data for nil JSON file name: %#v", data)
	}
}

func TestParseJSONWithNonExistingFileName(t *testing.T) {
	name := "does-not-exist"
	data, err := parseJSON(&name)

	if expectedMsg := "error reading JSON file: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error when parsing missing JSON file")
		t.Errorf("\tActual:        %q", err)
		t.Errorf("\tExpected part: %q", expectedMsg)
	}

	if len(data) != 0 {
		t.Errorf("Unexpected data for missing JSON file: %#v", data)
	}
}

func TestParseJSONWithMalformedJSON(t *testing.T) {
	jsonFile := createFile(t, "MALFORMED")
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	data, err := parseJSON(&jsonFile)

	if expectedMsg := "json: syntax error at offset 1: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error when parsing a file with malformed JSON")
		t.Errorf("\tActual:        %q", err)
		t.Errorf("\tExpected part: %q", expectedMsg)
	}

	if len(data) != 0 {
		t.Errorf("Unexpected data for malformed JSON file: %#v", data)
	}
}

func TestParseJSONWithNonStringJSONValues(t *testing.T) {
	jsonFile := createFile(t, `{"foo": true}`)
	defer func() {
		if err := os.Remove(jsonFile); err != nil {
			t.Fatalf("Unexpected error deleting temporary file: %s", err)
		}
	}()

	data, err := parseJSON(&jsonFile)

	if expectedMsg := "json: type error at offset 12: "; !strings.Contains(err.Error(), expectedMsg) {
		t.Error("Invalid error when parsing a file with JSON having non-string values")
		t.Errorf("\tActual:        %q", err)
		t.Errorf("\tExpected part: %q", expectedMsg)
	}

	if len(data) != 0 {
		t.Errorf("Unexpected data for JSON file having non-string values: %#v", data)
	}
}
