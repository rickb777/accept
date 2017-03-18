package accept

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAcceptsEncoding(t *testing.T) {

	examples := []struct {
		hdr string
		ok  bool
	}{
		{"gzip", true},
		{"*", false},
		{"", false},
		{"compress, gzip, br", true},
		{" compress ; q = 0.5 , gzip ", true},
		{"gzip;q=0.123", true},
		{"gzip;q=1.123", true},
		{"gzip;q = -1.123", false},
		{" compress ; q = 0.5 , gzip;q=0 ", false},
		{"this is not valid", false},
	}

	for i, ex := range examples {
		h := make(http.Header)
		h.Set(AcceptEncoding, ex.hdr)
		v := AcceptsEncoding(h, "gzip")
		if v != ex.ok {
			t.Errorf("Test %d: got %v but wanted %v (%s)", i, v, ex.ok, ex.hdr)
		}
	}
}

func TestPreferredLikeEdgeCases(t *testing.T) {

	examples := []struct {
		hdr  string
		like string
	}{
		{"", ""},
		{" compress ; q = fail , gzip ", "gzip"},
		{"compress; q=0.5 , gzip;q=0", "gzip"},
		{"this is not valid", "gzip"},
	}

	for i, ex := range examples {
		v := PreferredLike(ex.hdr, ex.like)
		if v != "" {
			t.Errorf("Test %d: got %s but wanted %s", i, v, ex.hdr)
		}
	}
}

func ExampleAcceptsEncoding() {
	h := make(http.Header)
	ae := "compress; q=0.5, gzip, *;q=0"
	h.Set(AcceptEncoding, ae)

	v := AcceptsEncoding(h, "gzip")

	fmt.Printf("%s: %s\n", AcceptEncoding, ae)
	fmt.Printf("accepted = %v\n", v)

	// Output: Accept-Encoding: compress; q=0.5, gzip, *;q=0
	// accepted = true
}

func ExampleParse_acceptEncoding() {
	ae := "compress;q=0.5, gzip, *;q=0"

	v, _ := Parse(ae)
	accepted := v.Get("gzip").IsAccepted()

	fmt.Printf("%s: %s\n", AcceptEncoding, ae)
	fmt.Printf("accepted = %v\n", accepted)

	// Output: Accept-Encoding: compress;q=0.5, gzip, *;q=0
	// accepted = true
}

func ExamplePreferredContentTypeLike() {
	accept := "text/html;q=1.0, text/*;q=0.8, image/gif;q=0.6, image/jpeg;q=0.6, image/*;q=0.5, */*;q=0.1"
	h := make(http.Header)
	h.Set(Accept, accept)

	preferred := PreferredContentTypeLike(h, "text/")

	fmt.Printf("%s: %s\n", Accept, accept)
	fmt.Printf("preferred = %v\n", preferred)

	// Output: Accept: text/html;q=1.0, text/*;q=0.8, image/gif;q=0.6, image/jpeg;q=0.6, image/*;q=0.5, */*;q=0.1
	// preferred = text/html
}

func ExamplePreferredCharsetLike() {
	accept := "iso-8859-5, unicode-1-1;q=0.8, iso-8859-1;q=0.1"
	h := make(http.Header)
	h.Set(AcceptCharset, accept)

	preferred := PreferredCharsetLike(h, "iso-8859-")

	fmt.Printf("%s: %s\n", AcceptCharset, accept)
	fmt.Printf("preferred = %v\n", preferred)

	// Output: Accept-Charset: iso-8859-5, unicode-1-1;q=0.8, iso-8859-1;q=0.1
	// preferred = iso-8859-5
}

func ExamplePreferredLanguageLike() {
	accept := "da, en-gb;q=0.8, en;q=0.7"
	h := make(http.Header)
	h.Set(AcceptLanguage, accept)

	preferredEN := PreferredLanguageLike(h, "en")

	fmt.Printf("%s: %s\n", AcceptLanguage, accept)
	fmt.Printf("preferredEN = %v\n", preferredEN)

	// Output: Accept-Language: da, en-gb;q=0.8, en;q=0.7
	// preferredEN = en-gb
}

func ExampleParse_acceptLanguage() {
	accept := "da, en-gb;q=0.8, en;q=0.7, pt;q=0"

	// Parse the header value
	v, _ := Parse(accept)
	// sort the codings according to weighting rules
	sorted := v.Sorted()
	// the first coding is the most preferred
	mostPreferred := sorted[0].Name
	// Like filters the codings that are start with a prefix or are "*"
	preferredEN := sorted.Like("en")[0].Name

	fmt.Printf("%s: %s\n", AcceptLanguage, accept)
	fmt.Printf("mostPreferred = %v\n", mostPreferred)
	fmt.Printf("preferredEN   = %v\n", preferredEN)

	// Output: Accept-Language: da, en-gb;q=0.8, en;q=0.7, pt;q=0
	// mostPreferred = da
	// preferredEN   = en-gb
}

func ExampleParse_accept() {
	accept := "text/html;q=1.0, text/*;q=0.8, image/gif;q=0.6, image/jpeg;q=0.7, image/*;q=0.5, */*;q=0.1"

	// Parse the header value
	v, _ := Parse(accept)
	// sort the codings according to weighting rules
	sorted := v.Sorted()
	// the first coding is the most preferred
	matchingText := sorted.Like("text/").Names()
	matchingImage := sorted.Like("image/").Names()

	fmt.Printf("%s: %s\n", Accept, accept)
	fmt.Printf("matchingText  = %v\n", matchingText)
	fmt.Printf("matchingImage = %v\n", matchingImage)

	// Output: Accept: text/html;q=1.0, text/*;q=0.8, image/gif;q=0.6, image/jpeg;q=0.7, image/*;q=0.5, */*;q=0.1
	// matchingText  = [text/html text/* */*]
	// matchingImage = [image/jpeg image/gif image/* */*]
}
