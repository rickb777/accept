package accept

import (
	"fmt"
	"reflect"
	"testing"
)

func TestCodingString(t *testing.T) {

	examples := []struct {
		c   Coding
		exp string
	}{
		{Coding{"gzip", 1, nil}, "gzip"},
		{Coding{"gzip", 0.1234, nil}, "gzip;q=0.123"},
		{Coding{"text/plain", 0, map[string]string{"format": "flowed"}}, "text/plain;format=flowed;q=0"},
		{Coding{"text/plain", 1, nil}, "text/plain"},
		{Coding{"text/html", 0.4, map[string]string{"level": "2"}}, "text/html;level=2;q=0.4"},
	}

	for i, ex := range examples {
		s := ex.c.String()
		if s != ex.exp {
			t.Errorf("Test %d: got '%s', expected '%s'", i, s, ex.exp)
		}
	}
}

func TestCodingsString(t *testing.T) {
	c := Codings{
		Coding{"text/html", 0.4, map[string]string{"level": "2"}},
		Coding{"image/png", 0.1234, nil},
		Coding{"*/*", 0, nil},
	}

	s := c.String()

	exp := "text/html;level=2;q=0.4, image/png;q=0.123, */*;q=0"
	if s != exp {
		t.Errorf("Got '%s', expected '%s'", s, exp)
	}
}

//"gzip;q=1.0, identity; q=0.5, *;q=0": {"gzip": 1.0, "identity": 0.5, "*": 0.0},
func TestParseHappy(t *testing.T) {

	examples := []struct {
		c   Codings
		hdr string
	}{
		{Codings{Coding{"gzip", 1, nil}},
			"gzip"},

		{Codings{Coding{"*", 1, nil}},
			"*"},

		{Codings{},
			""},

		{Codings{Coding{"compress", 1, nil}, Coding{"gzip", 1, nil}},
			"compress, gzip"},

		{Codings{Coding{"compress", 0.5, nil}, Coding{"gzip", 1, nil}},
			" Compress ; q = 0.5 , GZIP "},

		{Codings{Coding{"gzip", 1, nil}, Coding{"identity", 0.5, nil}, Coding{"*", 0, nil}},
			"gzip;q=1.0, identity; q=0.5, *;q=0"},

		{Codings{Coding{"gzip", 0.123, nil}},
			"gzip;q=0.123"},

		{Codings{Coding{"gzip", 1, nil}},
			"gzip;q=1.123"},

		{Codings{Coding{"gzip", 0, nil}},
			"gzip;q = -1.123"},

		{Codings{Coding{"text/plain", 1, map[string]string{"format": "flowed"}}},
			"text/plain;format=flowed"},

		{Codings{Coding{"text/plain", 1, nil}},
			"text/plain"},

		{Codings{Coding{"text/html", 0.4, map[string]string{"level": "2"}}},
			"text/html;Level=2;Q=0.4"},

		{Codings{Coding{"text/html", 1, map[string]string{"charset": "utf-8"}}},
			"Text/HTML;Charset=UTF-8"},

		// quoted attribute values are not yet supported (and may not be required)
		//{Codings{Coding{"text/html", 1, map[string]string{"charset": "utf-8"}}},
		//	`Text/HTML;Charset="UTF-8"`},
	}

	for i, ex := range examples {

		v, err := Parse(ex.hdr)

		if err != nil {
			t.Errorf("Test %d: got error '%v'", i, err)
		}
		if !reflect.DeepEqual(v, ex.c) {
			t.Errorf("Test %d: got '%#v', expected '%#v'", i, v, ex.c)
		}
	}
}

func TestSortedCodingsLike(t *testing.T) {

	examples := []struct {
		hdr         string
		like        string
		preferences []string
	}{
		{"", "en", []string{}},

		{"text/plain", "text/", []string{"text/plain"}},

		{"text/plain;q=0.5, text/html", "text/", []string{"text/html", "text/plain;q=0.5"}},

		{"text/html;level=2;q=0.4", "text/", []string{"text/html;level=2;q=0.4"}},

		{"audio/*;q=0.2, audio/basic", "audio/", []string{"audio/basic", "audio/*;q=0.2"}},

		{"text/*, text/plain, text/plain;format=flowed, */*", "text/", []string{"text/plain;format=flowed", "text/plain", "text/*", "*/*"}},
	}

	for i, ex := range examples {
		cs, err := Parse(ex.hdr)
		if err != nil {
			t.Errorf("Test %d: got error '%v'", i, err)
		}

		v := cs.Sorted().Like(ex.like)

		if len(v) != len(ex.preferences) {
			t.Errorf("Test %d: got %d [%v], expected %d %v", i, len(v), v, len(ex.preferences), ex.preferences)
		} else {
			for j, x := range v {
				if x.String() != ex.preferences[j] {
					t.Errorf("Test %d.%d: got %v, expected %s", i, j, v, ex.preferences[j])
				}
			}
		}
	}
}

func TestCodingsIfAccepted(t *testing.T) {

	examples := []struct {
		hdr         string
		preferences []string
	}{
		{"", []string{}},

		{"gzip;q=1.0, identity; q=0.5, *;q=0", []string{"gzip", "identity"}},
		{"*;q=0, gzip;q=1.0, identity; q=0", []string{"gzip"}},
	}

	for i, ex := range examples {
		cs, err := Parse(ex.hdr)
		if err != nil {
			t.Errorf("Test %d: got error '%v'", i, err)
		}

		names := cs.IfAccepted().Names()

		if !reflect.DeepEqual(names, ex.preferences) {
			t.Errorf("Test %d: got '%#v', expected '%#v'", i, names, ex.preferences)
		}
	}
}

func ExampleCodings_Sorted() {
	accept := "text/*, text/plain, text/plain;format=flowed, */*"

	v, _ := Parse(accept)
	sorted := v.Sorted()

	fmt.Printf("%s: %s\n", Accept, accept)
	fmt.Printf("sorted = %v\n", sorted)

	// Output: Accept: text/*, text/plain, text/plain;format=flowed, */*
	// sorted = text/plain;format=flowed, text/plain, text/*, */*
}
