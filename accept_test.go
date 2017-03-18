package accept

import (
	"fmt"
	"net/http"
	"testing"
)

func TestAcceptsGzipEncoding(t *testing.T) {

	examples := []struct {
		hdr string
		ok  bool
	}{
		{"gzip", true},
		{"*", false},
		{"", false},
		{"compress, gzip", true},
		{" compress ; q = 0.5 , gzip ", true},
		{"gzip;q=0.123", true},
		{"gzip;q=1.123", true},
		{"gzip;q = -1.123", false},
		{" compress ; q = 0.5 , gzip;q=0 ", false},
	}

	for i, ex := range examples {
		h := make(http.Header)
		h.Set(AcceptEncoding, ex.hdr)
		v := AcceptsGzipEncoding(h)
		if v != ex.ok {
			t.Errorf("Test %d: got %v but wanted %v (%s)", i, v, ex.ok, ex.hdr)
		}
	}
}

func ExampleAcceptsGzipEncoding() {
	h := make(http.Header)
	ae := "compress; q=0.5, gzip, *;q=0"
	h.Set(AcceptEncoding, ae)

	v := AcceptsGzipEncoding(h)

	fmt.Printf("%s: %s\n", AcceptEncoding, ae)
	fmt.Printf("accepted = %v\n", v)

	// Output: Accept-Encoding: compress; q=0.5, gzip, *;q=0
	// accepted = true
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

func ExampleParseCodings_acceptEncoding() {
	ae := "compress;q=0.5, gzip, *;q=0"

	v, _ := ParseCodings(ae)
	accepted := v.Get("gzip").IsAccepted()

	fmt.Printf("%s: %s\n", AcceptEncoding, ae)
	fmt.Printf("accepted = %v\n", accepted)

	// Output: Accept-Encoding: compress;q=0.5, gzip, *;q=0
	// accepted = true
}

func ExampleParseCodings_accept() {
	accept := "audio/*; q=0.2, audio/basic"

	v, _ := ParseCodings(accept)
	sorted := v.Sorted()
	preferred := sorted[0]

	fmt.Printf("%s: %s\n", Accept, accept)
	fmt.Printf("preferred = %v\n", preferred)

	// Output: Accept: audio/*; q=0.2, audio/basic
	// preferred = audio/basic
}

func ExampleParseCodings_acceptLanguage() {
	accept := "da, en-gb;q=0.8, en;q=0.7"

	v, _ := ParseCodings(accept)
	sorted := v.Sorted()
	mostPreferred := sorted[0].Name
	preferredEN := sorted.Like("en")[0].Name

	fmt.Printf("%s: %s\n", AcceptLanguage, accept)
	fmt.Printf("mostPreferred = %v\n", mostPreferred)
	fmt.Printf("preferredEN   = %v\n", preferredEN)

	// Output: Accept-Language: da, en-gb;q=0.8, en;q=0.7
	// mostPreferred = da
	// preferredEN   = en-gb
}
