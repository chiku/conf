package conf

import (
	"io/ioutil"
	"testing"
)

// createFile creates a temporary file with the given contents.
// It returns the file-name of the created file.
// The caller is expected to delete the created file.
// The test aborts on failure.
func createFile(t *testing.T, content string) string {
	t.Helper()

	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatalf("Unexpected error creating temporary file: %s", err)
	}

	_, err = tmpfile.WriteString(content)
	if err != nil {
		t.Fatalf("Unexpected error writing to temporary file: %s", err)
	}

	err = tmpfile.Close()
	if err != nil {
		t.Fatalf("Unexpected error closing temporary file: %s", err)
	}

	return tmpfile.Name()
}
