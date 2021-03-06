package filter

import (
	"reflect"
	"testing"
)

var filterFilesTests = []struct {
	includes []string
	excludes []string
	files    []string
	expected []string
	err      bool
}{
	{
		nil,
		[]string{"*"},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{},
		false,
	},
	{
		[]string{"*"},
		nil,
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		false,
	},
	{
		[]string{"*"},
		[]string{"*.go"},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{"main.cpp", "main.h", "bar.py"},
		false,
	},
	// Invalid patterns won't match anything. This would trigger a warning at
	// runtime.
	{
		[]string{"*"},
		[]string{"[["},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		true,
	},

	{
		[]string{"main.*"},
		[]string{"*.cpp"},
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{"main.go", "main.h"},
		false,
	},
	{
		nil, nil,
		[]string{"main.cpp", "main.go", "main.h", "foo.go", "bar.py"},
		[]string{},
		false,
	},

	{
		[]string{"**/*"},
		nil,
		[]string{"foo", "/test/foo", "/test/foo.go"},
		[]string{"foo", "/test/foo", "/test/foo.go"},
		false,
	},
}

func TestFilterFiles(t *testing.T) {
	for i, tt := range filterFilesTests {
		result, err := Files(tt.files, tt.includes, tt.excludes)
		if !tt.err && err != nil {
			t.Errorf("Test %d: error %s", i, err)
		}
		if !reflect.DeepEqual(result, tt.expected) {
			t.Errorf(
				"Test %d (inc: %v, ex: %v), expected \"%v\" got \"%v\"",
				i, tt.includes, tt.excludes, tt.expected, result,
			)
		}
	}
}

var basePathTests = []struct {
	pattern  string
	expected string
}{
	{"foo", "."},
	{"test/foo", "test"},
	{"test/foo*", "test"},
	{"test/*.**", "test"},
	{"**/*", "."},
	{"foo*/bar", "."},
	{"foo/**/bar", "foo"},
	{"/voing/**", "/voing"},
}

func TestBasePath(t *testing.T) {
	for i, tt := range basePathTests {
		ret := BasePath(tt.pattern)
		if ret != tt.expected {
			t.Errorf("%d: %q - Expected %q, got %q", i, tt.pattern, tt.expected, ret)
		}
	}
}

var getBasePathTests = []struct {
	patterns []string
	expected []string
}{
	{[]string{"foo"}, []string{"."}},
	{[]string{"foo", "bar"}, []string{"."}},
	{[]string{"foo", "bar", "/voing/**"}, []string{".", "/voing"}},
	{[]string{"foo/**", "**"}, []string{"."}},
	{[]string{"foo/**", "**", "/bar/**"}, []string{".", "/bar"}},
}

func TestGetBasePaths(t *testing.T) {
	for i, tt := range getBasePathTests {
		bp := []string{}
		bp = GetBasePaths(bp, tt.patterns)
		if !reflect.DeepEqual(bp, tt.expected) {
			t.Errorf("%d: %#v - Expected %#v, got %#v", i, tt.patterns, tt.expected, bp)
		}
	}
}

var findTests = []struct {
	include  []string
	exclude  []string
	expected []string
}{
	{
		[]string{"**"},
		[]string{},
		[]string{"a/a.test1", "a/b.test2", "b/a.test1", "b/b.test2", "x", "x.test1"},
	},
	{
		[]string{"**/*.test1"},
		[]string{},
		[]string{"a/a.test1", "b/a.test1", "x.test1"},
	},
	{
		[]string{"**"},
		[]string{"*.test1"},
		[]string{"a/a.test1", "a/b.test2", "b/a.test1", "b/b.test2", "x"},
	},
	{
		[]string{"**"},
		[]string{"a"},
		[]string{"b/a.test1", "b/b.test2", "x", "x.test1"},
	},
	{
		[]string{"**"},
		[]string{"a/"},
		[]string{"b/a.test1", "b/b.test2", "x", "x.test1"},
	},
	{
		[]string{"**"},
		[]string{"**/*.test1", "**/*.test2"},
		[]string{"x"},
	},
}

func TestFind(t *testing.T) {
	for i, tt := range findTests {
		ret, err := Find("./test/find", tt.include, tt.exclude)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(ret, tt.expected) {
			t.Errorf(
				"%d: %#v, %#v - Expected\n%#v\ngot:\n%#v",
				i, tt.include, tt.exclude, tt.expected, ret,
			)
		}
	}
}
