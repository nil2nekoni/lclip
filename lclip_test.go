package lclip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/mitchellh/go-homedir"
)

var tempPath string

func TestMain(m *testing.M) {
	f, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		fmt.Fprintln(os.Stderr, "lclip_test:", err)
		return
	}
	f.Close()
	tempPath = f.Name()

	e := m.Run()
	defer os.Exit(e)

	os.Remove(tempPath)
}

func TestDefaultPath(t *testing.T) {
	h, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	expect := filepath.Join(h, ".lclip.json")
	actual, err := DefaultPath()
	if err != nil {
		t.Errorf("DefaultPath returns %q; want nil", err)
	}
	if actual != expect {
		t.Errorf("DefaultPath = %q; want %q",
			actual, expect)
	}
}

type GetTextTest struct {
	Src string
	Dst string
}

var indexTestsGetText = []GetTextTest{
	{Src: "foo", Dst: "bar"},
	{Src: "hoge", Dst: "piyo"},
}

func TestGetText(t *testing.T) {
	w := bytes.NewBuffer(make([]byte, 0))
	m := make(map[string]string)
	for _, test := range indexTestsGetText {
		m[test.Src] = test.Dst
	}
	if err := json.NewEncoder(w).Encode(m); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(tempPath, w.Bytes(), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	for _, test := range indexTestsGetText {
		expect := test.Dst
		actual := c.Get(test.Src)
		if actual != expect {
			t.Errorf("Get(%q) = %q; want %q",
				test.Src, actual, expect)
		}
	}
}

type SetTextTest struct {
	Label string
	Data  string
}

var indexTestsSetText = []SetTextTest{
	{Label: "a", Data: "aaa"},
	{Label: "abc", Data: "def"},
	{Label: "", Data: ""},
}

func TestSetText(t *testing.T) {
	if err := ioutil.WriteFile(tempPath, []byte(`{}`), 0644); err != nil {
		t.Fatal(err)
	}

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Errorf("NewClipboard returns %q; want nil", err)
	}
	for _, test := range indexTestsSetText {
		c.Set(test.Label, test.Data)
		expect := test.Data
		actual := c.Get(test.Label)
		if actual != expect {
			t.Errorf("after Set(%q, %q), Get(%q) = %q; want %q",
				test.Label, test.Data,
				test.Label, actual, expect)
		}
	}
}
