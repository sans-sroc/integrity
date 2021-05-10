package utils

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sans-sroc/integrity/pkg/common"
)

func TestDefaultGetFiles(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/default/five.txt",
		},
		{
			Name: "folder-one/folder-two/two-two.txt",
			Path: "testdata/default/folder-one/folder-two/two-two.txt",
		},
		{
			Name: "folder-one/one-one.txt",
			Path: "testdata/default/folder-one/one-one.txt",
		},
		{
			Name: "one.txt",
			Path: "testdata/default/one.txt",
		},
		{
			Name: "two.txt",
			Path: "testdata/default/two.txt",
		},
	}

	for _, e := range expected {
		e.Name = filepath.ToSlash(e.Name)
		e.Path = filepath.ToSlash(e.Path)
	}

	directory := filepath.ToSlash("testdata/default/")

	files, err := GetFiles(directory)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, files) {
		t.Errorf("files do not match")
	}

}

func TestWithGetFirstGetFiles(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/with-get-first/five.txt",
		},
		{
			Name: "get_first/folder-two/two-two.txt",
			Path: "testdata/with-get-first/get_first/folder-two/two-two.txt",
		},
		{
			Name: "get_first/one-one.txt",
			Path: "testdata/with-get-first/get_first/one-one.txt",
		},
		{
			Name: "one.txt",
			Path: "testdata/with-get-first/one.txt",
		},
		{
			Name: "two.txt",
			Path: "testdata/with-get-first/two.txt",
		},
	}

	for _, e := range expected {
		e.Name = filepath.ToSlash(e.Name)
		e.Path = filepath.ToSlash(e.Path)
	}

	directory := filepath.ToSlash("testdata/with-get-first/")

	files, err := GetFiles(directory)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, files) {
		t.Errorf("files do not match")
	}

	for i, a := range files {
		if i > len(expected) {
			t.Errorf("File not in expected: Name: %s, Path: %s", a.Name, a.Path)
			return
		}

		e := expected[i]

		if a.Path != e.Path {
			t.Errorf("%s path does not match. Expected: %s, Actual: %s", a.Name, e.Path, a.Path)
		}
	}
}
