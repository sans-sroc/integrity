package utils

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/sans-sroc/integrity/pkg/common"
)

func TestDefaultHashFileSha256(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/default/five.txt",
			Hash: "ee0874170b7f6f32b8c2ac9573c428d35b575270a66b757c2c0185d2bd09718d",
		},
		{
			Name: "folder-one/folder-two/two-two.txt",
			Path: "testdata/default/folder-one/folder-two/two-two.txt",
			Hash: "7f37056e2015eeeb5e45a948ba1979fe4d733b38174050776090db2255ed8907",
		},
		{
			Name: "folder-one/one-one.txt",
			Path: "testdata/default/folder-one/one-one.txt",
			Hash: "2410ee537f887eb9f3063c8bb6da0b73dce38dcf00729bb3664c9e24d6c85304",
		},
		{
			Name: "one.txt",
			Path: "testdata/default/one.txt",
			Hash: "a7937b64b8caa58f03721bb6bacf5c78cb235febe0e70b1b84cd99541461a08e",
		},
		{
			Name: "two.txt",
			Path: "testdata/default/two.txt",
			Hash: "8b5b9db0c13db24256c829aa364aa90c6d2eba318b9232a4ab9313b954d3555f",
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

	for i, r := range files {
		e := expected[i]

		if r.Name == e.Name && e.Hash != r.Hash {
			t.Errorf("%s hash does not match. Expected: %s, Actual: %s", r.Name, e.Hash, r.Hash)
		}
	}
}

func TestWithGetFirstHashFileSha256(t *testing.T) {
	expected := []*common.File{
		{
			Name: "five.txt",
			Path: "testdata/with-get-first/five.txt",
			Hash: "ee0874170b7f6f32b8c2ac9573c428d35b575270a66b757c2c0185d2bd09718d",
		},
		{
			Name: "get_first/folder-two/two-two.txt",
			Path: "testdata/with-get-first/get_first/folder-two/two-two.txt",
			Hash: "7f37056e2015eeeb5e45a948ba1979fe4d733b38174050776090db2255ed8907",
		},
		{
			Name: "get_first/one-one.txt",
			Path: "testdata/with-get-first/get_first/one-one.txt",
			Hash: "2410ee537f887eb9f3063c8bb6da0b73dce38dcf00729bb3664c9e24d6c85304",
		},
		{
			Name: "one.txt",
			Path: "testdata/with-get-first/one.txt",
			Hash: "a7937b64b8caa58f03721bb6bacf5c78cb235febe0e70b1b84cd99541461a08e",
		},
		{
			Name: "two.txt",
			Path: "testdata/with-get-first/two.txt",
			Hash: "8b5b9db0c13db24256c829aa364aa90c6d2eba318b9232a4ab9313b954d3555f",
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

	for i, r := range files {
		e := expected[i]

		if r.Name == e.Name && e.Hash != r.Hash {
			t.Errorf("%s hash does not match. Expected: %s, Actual: %s", r.Name, e.Hash, r.Hash)
		}
	}
}
