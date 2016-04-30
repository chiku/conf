package conf

import (
	"fmt"
	"io/ioutil"
	"path"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func requireNoError(t *testing.T, err error, msg string) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %s\n\n", fileBase, line, err.Error())
		t.FailNow()
	}
}

func requireError(t *testing.T, err error, msg string) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		t.FailNow()
	}
}

func assertEqual(t *testing.T, actual, expected, msg interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %#v != %#v\n\n", fileBase, line, actual, expected)
		t.Fail()
	}
}

func assertContains(t *testing.T, total, part, msg string) {
	if !strings.Contains(total, part) {
		_, file, line, _ := runtime.Caller(1)
		fileBase := path.Base(file)

		fmt.Printf("\t%v:%v: %s\n", fileBase, line, msg)
		fmt.Printf("\t%v:%v: %#v doesn't contain %#v\n\n", fileBase, line, total, part)
		t.Fail()
	}
}

func createFile(t *testing.T, content string) string {
	tmpfile, err := ioutil.TempFile("", "example")
	requireNoError(t, err, "Expected no error creating temporary file")
	_, err = tmpfile.Write([]byte(content))
	requireNoError(t, err, "Expected no error writing to temporary file")
	err = tmpfile.Close()
	requireNoError(t, err, "Expected no error closing temporary file")

	return tmpfile.Name()
}
