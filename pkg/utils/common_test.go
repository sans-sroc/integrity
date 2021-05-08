package utils

import (
	"reflect"
	"testing"

	"github.com/sans-sroc/integrity/pkg/common"
)

func TestGetFiles(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/five.txt",
		},
		{
			Name: "one.txt",
			Path: "testdata/one.txt",
		},
		{
			Name: "two.txt",
			Path: "testdata/two.txt",
		},
	}

	files, err := GetFiles("testdata/")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, files) {
		t.Errorf("files do not match")
	}

}
