package utils

import (
	"reflect"
	"testing"

	"github.com/sans-sroc/integrity/pkg/common"
)

func TestHashFileSha256(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/five.txt",
			Hash: "ee0874170b7f6f32b8c2ac9573c428d35b575270a66b757c2c0185d2bd09718d",
		},
		{
			Name: "one.txt",
			Path: "testdata/one.txt",
			Hash: "a7937b64b8caa58f03721bb6bacf5c78cb235febe0e70b1b84cd99541461a08e",
		},
		{
			Name: "two.txt",
			Path: "testdata/two.txt",
			Hash: "8b5b9db0c13db24256c829aa364aa90c6d2eba318b9232a4ab9313b954d3555f",
		},
	}

	files, err := GetFiles("testdata/")
	if err != nil {
		t.Error(err)
		return
	}

	for _, file := range files {
		hash, err := HashFileSha256(file.Path)
		if err != nil {
			t.Error(err)
			return
		}

		file.Hash = hash
	}

	if !reflect.DeepEqual(expected, files) {
		t.Error("files do not match")
	}
}
