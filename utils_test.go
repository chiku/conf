package conf

// utils_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2016-2017. All rights reserved
// License::   MIT

import (
	"os"
	"testing"
)

func TestParseJSON(t *testing.T) {
	jsonFile := createFile(t, `{ "foo": "abc", "bar": "xyz" }`)
	defer os.Remove(jsonFile)

	data, err := parseJSON(&jsonFile)

	requireNoError(t, err, "Expected no error parsing valid JSON")
	expectedData := map[string]string{
		"foo": "abc",
		"bar": "xyz",
	}
	assertEqual(t, data, expectedData, "Expected JSON parse data to create a string map")
}

func TestParseJSONWithoutFileName(t *testing.T) {
	name := ""
	emptyData, err := parseJSON(&name)

	requireNoError(t, err, "Expected no error parsing empty JSON file name")
	assertEqual(t, len(emptyData), 0, "Expected no JSON data")

	nilData, err := parseJSON(nil)

	requireNoError(t, err, "Expected no error parsing nil JSON file name")
	assertEqual(t, len(nilData), 0, "Expected no JSON data")
}

func TestParseJSONWithNonExistingFileName(t *testing.T) {
	name := "does-not-exist"
	data, err := parseJSON(&name)

	requireError(t, err, "Expected error parsing non-existing JSON file name")
	assertContains(t, err.Error(), "error reading JSON file: ", "Expected file does not exist error")
	assertEqual(t, len(data), 0, "Expected no JSON data")
}

func TestParseJSONWithMalformedJSON(t *testing.T) {
	jsonFile := createFile(t, `MALFORMED`)
	defer os.Remove(jsonFile)

	data, err := parseJSON(&jsonFile)

	requireError(t, err, "Expected error parsing malformed JSON file")
	assertContains(t, err.Error(), "json: syntax error at offset 1: ", "Expected JSON parse error")
	assertEqual(t, len(data), 0, "Expected no JSON data")
}

func TestParseJSONWithNonStringJSONValues(t *testing.T) {
	jsonFile := createFile(t, `{"foo": true}`)
	defer os.Remove(jsonFile)

	data, err := parseJSON(&jsonFile)

	requireError(t, err, "Expected error parsing malformed JSON file")
	assertContains(t, err.Error(), "json: type error at offset 12: ", "Expected JSON parse error")
	assertEqual(t, len(data), 0, "Expected no JSON data")
}
