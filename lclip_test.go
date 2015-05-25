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

type GetTextTest struct {
	Src string
	Dst []byte
}

var indexTestsGetText = []GetTextTest{
	{Src: "foo", Dst: []byte("bar")},
	{Src: "hoge", Dst: []byte("piyo")},
}

func TestGetText(t *testing.T) {
	os.Remove(tempPath)

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Fatal(err)
	}
	for _, test := range indexTestsGetText {
		c.Set(test.Src, test.Dst)
	}
	defer c.Close()
	for _, test := range indexTestsGetText {
		expect := test.Dst
		actual := c.Get(test.Src)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Get(%q) = %q; want %q",
				test.Src, actual, expect)
		}
	}
}

type SetTextTest struct {
	Label string
	Data  []byte
}

var indexTestsSetText = []SetTextTest{
	{Label: "a", Data: []byte("aaa")},
	{Label: "abc", Data: []byte("def")},
	{Label: "", Data: []byte("")},
}

func TestSetText(t *testing.T) {
	os.Remove(tempPath)

	c, err := NewClipboard(tempPath)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()
	for _, test := range indexTestsSetText {
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

func TestSaveText(t *testing.T) {
	os.Remove(tempPath)

	k, v := "key", []byte("value")
	{
		c, err := NewClipboard(tempPath)
		if err != nil {
			t.Fatal(err)
		}
		c.Set(k, v)
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
		expect := v
		actual := c.Get(k)
		if !reflect.DeepEqual(actual, expect) {
			t.Errorf("Get(%q) = %q; want %q",
				k, actual, expect)
		}
	}
}
