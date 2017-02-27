package conf

// helpers_test.go
//
// Author::    Chirantan Mitra
// Copyright:: Copyright (c) 2016-2017. All rights reserved
// License::   MIT

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

// requireNoError verifies that err is nil. It prints the given message if err is not nil.
// The test is aborted on failure.
func requireNoError(t *testing.T, err error, msg string) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %s\n\n", fileBase, line, err.Error())
		t.FailNow()
	}
}

// requireError verifies that err not nil. It prints the given message if err is nil.
// The test is aborted on failure.
func requireError(t *testing.T, err error, msg string) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		t.FailNow()
	}
}

// assertError verifies the actual equal expected. It prints the given message if the two aren't equal.
// The equality is checked using reflect.DeepEquals. The test continues on failure.
func assertEqual(t *testing.T, actual, expected, msg interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %#v != %#v\n\n", fileBase, line, actual, expected)
		t.Fail()
	}
}

// assertContains verifies part is a sub-string of total. It prints the given message if it isn't.
// The test continues on failure.
func assertContains(t *testing.T, total, part, msg string) {
	if !strings.Contains(total, part) {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %#v doesn't contain %#v\n\n", fileBase, line, total, part)
		t.Fail()
	}
}

// createFile creates a temporary file with the given contents. It returns the file-name of the created file.
// The caller is expected to delete the created file.
// The test aborts on failure.
func createFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "example")
	requireNoError(t, err, "Expected no error creating temporary file")
	_, err = tmpfile.Write([]byte(content))
	requireNoError(t, err, "Expected no error writing to temporary file")
	err = tmpfile.Close()
	requireNoError(t, err, "Expected no error closing temporary file")

	return tmpfile.Name()
}
