package integrity

import "testing"

var valid = []string{
	"123.21.1",
	"123A.21.1",
	"123.21.1A",
	"123A.21.1A",
}

var invalid = []string{
	"SEC123.21.1",
	"123.21.10",
	"123AA.21.1",
	"123.21.1AA",
	"12.12.1",
	"123.2.2",
}

func TestValidNames(t *testing.T) {
	for _, v := range valid {
		integrity, err := New("/tmp", false)
		if err != nil {
			t.Error(err)
		}

		if err := integrity.SetName(v); err != nil {
			t.Error(err)
		}
	}
}

func TestInvalidNames(t *testing.T) {
	for _, v := range invalid {
		integrity, err := New("/tmp", false)
		if err != nil {
			t.Error(err)
		}

		if err := integrity.SetName(v); err == nil {
			t.Errorf("Name matched and it should not have: %s", v)
		}
	}
}
