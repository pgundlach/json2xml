package json2xml

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSimple(t *testing.T) {
	f, err := os.Open(filepath.Join("testdata", "simple.json"))
	if err != nil {
		t.Error("Could not read test data simple.json")
	}
	str, err := ToXML(f)
	if err != nil {
		t.Error(err)
	}

	expected := `<data><map><array key="whatever"><entry>foo</entry><entry>3.45</entry><entry>bar</entry><entry>1</entry></array><map key="something"><entry key="another">object</entry><array key="and an"><entry>array</entry></array></map></map></data>`
	if str != expected {
		t.Errorf("ToXML() got: %q, want %q", str, expected)
	}
}
