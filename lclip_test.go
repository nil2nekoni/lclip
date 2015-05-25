package lclip

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
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
	expect := filepath.Join(h, ".lclip.db")
	actual, err := DefaultPath()
	if err != nil {
		t.Errorf("DefaultPath returns %q; want nil", err)
	}
	if actual != expect {
		t.Errorf("DefaultPath = %q; want %q",
			actual, expect)
	}
}

func TestCreateStorageFileIfNotExists(t *testing.T) {
	os.Remove(tempPath)

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if _, err := os.Stat(tempPath); err != nil {
		t.Error("not create storage file; want create")
	}
}

type AccessTest struct {
	Label string
	Data  []byte
}

var indexTestsAccess = []AccessTest{
	{Label: "foo", Data: []byte("bar")},
	{Label: "hoge", Data: []byte("piyo")},
	{Label: "日本語", Data: []byte("日本語")},
}

func TestSetText(t *testing.T) {
	os.Remove(tempPath)

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for _, test := range indexTestsAccess {
		c.Set(test.Label, test.Data)
		expect := test.Data
		actual := c.Get(test.Label)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("after Set(%q, %q), Get(%q) = %q; want %q",
				test.Label, test.Data,
				test.Label, actual, expect)
		}
	}
}

var indexTestsLabels = [][]string{
	{"foo", "bar", "baz"},
	{"hoge", "piyo", "fuga"},
}

func TestListLabels(t *testing.T) {
	for _, labels := range indexTestsLabels {
		os.Remove(tempPath)

		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Fatal(err)
		}
		for _, label := range labels {
			c.Set(label, []byte(``))
		}

		expect := append(make([]string, 0, len(labels)), labels...)
		actual := c.Labels()
		sort.Strings(expect)
		sort.Strings(actual)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("got %q; want %q", actual, expect)
		}
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func TestSave(t *testing.T) {
	os.Remove(tempPath)

	{
		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Fatal(err)
		}
		for _, test := range indexTestsAccess {
			c.Set(test.Label, test.Data)
		}
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}
	{
		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Fatal(err)
		}
		defer c.Close()
		for _, test := range indexTestsAccess {
			expect := test.Data
			actual := c.Get(test.Label)
			if !reflect.DeepEqual(actual, expect) {
				t.Errorf("Get(%q) = %q; want %q",
					test.Label, actual, expect)
			}
		}
	}
}
