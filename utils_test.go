package conf

import (
	"os"
	"testing"
)

func TestIsPresentInside(t *testing.T) {
	list := []string{"foo", "bar"}

	assertEqual(t, isPresentInside(list, "foo"), true, "Expected item to exist inside slice")
	assertEqual(t, isPresentInside(list, "foo1"), false, "Expected different item to not exist inside slice")
	assertEqual(t, isPresentInside(list, "FOO"), false, "Expected item with wrong case to not exist inside slice")
}

func TestIsPresentInsideForEmptySlice(t *testing.T) {
	assertEqual(t, isPresentInside([]string{}, ""), false, "Expected empty item to not exist inside empty slice")
}

func TestPartitionByUniqueness(t *testing.T) {
	list := []string{"uniq1", "dupl2", "dupl1", "dupl2", "dupl2", "dupl2", "dupl1", "uniq2"}

	uniq, dupl := partitionByUniqueness(list)

	assertEqual(t, uniq, []string{"uniq1", "uniq2"}, "Expected unique to maintain original order")
	assertEqual(t, dupl, []string{"dupl2", "dupl1"}, "Expected duplicate items to maintain original order")
}

func TestPartitionByUniquenessForEmptySlice(t *testing.T) {
	uniq, dupl := partitionByUniqueness([]string{})

	assertEqual(t, len(uniq), 0, "Expected no unique items")
	assertEqual(t, len(dupl), 0, "Expected no duplicate items")
}

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
	assertContains(t, err.Error(), "error parsing JSON file: ", "Expected JSON parse error")
	assertContains(t, err.Error(), "(syntax error at offset: 1)", "Expected JSON syntax error")
	assertEqual(t, len(data), 0, "Expected no JSON data")
}

func TestParseJSONWithNonStringJSONValues(t *testing.T) {
	jsonFile := createFile(t, `{"foo": true}`)
	defer os.Remove(jsonFile)

	data, err := parseJSON(&jsonFile)

	requireError(t, err, "Expected error parsing malformed JSON file")
	assertContains(t, err.Error(), "error parsing JSON file: ", "Expected JSON parse error")
	assertContains(t, err.Error(), "(type error at offset: 12)", "Expected JSON type error")
	assertEqual(t, len(data), 0, "Expected no JSON data")
}
