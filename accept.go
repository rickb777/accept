package accept

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const (
	// Accept is the "Accept" header.
	Accept = "Accept"

	// AcceptCharset is the "Accept-Charset" header.
	AcceptCharset = "Accept-Charset"

	// AcceptEncoding is the "Accept-Encoding" header.
	AcceptEncoding = "Accept-Encoding"

	// AcceptLanguage is the "Accept-Language" header.
	AcceptLanguage = "Accept-Language"

	// Identity is the coding name used for Accept-Encoding to specify settings for un-encoded behaviour,
	// e.g. "identity;q=0" excludes un-encoded behaviour.
	Identity = "identity"
)

//-------------------------------------------------------------------------------------------------

// AcceptsGzipEncoding determines whether there is a request header that indicates that the client
// accepts "gzip" content coding.
func AcceptsGzipEncoding(hdr http.Header) bool {
	return AcceptsEncoding(hdr, "gzip")
}

// AcceptsEncoding determines whether there is a request header that indicates that it will
// accept a particular coding for the response. The most used coding is "gzip".
func AcceptsEncoding(hdr http.Header, coding string) bool {
	acceptedEncodings, _ := ParseCodings(hdr.Get(AcceptEncoding))
	return acceptedEncodings.Get(coding).IsAccepted()
}

// ParseCodings returns the codings from a given header value. They are returned in the order they
// appear; use the Sorted() method on the result if you need them in preference order instead.
func ParseCodings(acceptValue string) (Codings, error) {
	c := make(Codings, 0)

	parts := strings.Split(acceptValue, ",")

	for _, part := range parts {
		coding, err := parsePart(part)

		if err != nil {
			return c, err
		} else if coding.Name != "" {
			c = append(c, coding)
		}
	}

	return c, nil
}

func parsePart(s string) (Coding, error) {
	fields := strings.Split(s, ";")
	var err error
	coding := Coding{
		Name:   strings.ToLower(strings.TrimSpace(fields[0])),
		Weight: 1,
	}

	for _, f := range fields[1:] {
		f = strings.ToLower(strings.TrimSpace(f))

		p := strings.Split(f, "=")
		switch len(p) {
		case 1:
			if coding.Attributes == nil {
				coding.Attributes = make(map[string]string)
			}
			coding.Attributes[f] = ""

		case 2:
			p0 := strings.TrimSpace(p[0])
			p1 := strings.TrimSpace(p[1])
			if p0 == "q" {
				coding.Weight, err = strconv.ParseFloat(p1, 64)
				if err != nil {
					return coding, fmt.Errorf("Cannot parse q-value in '%s'; %v", s, err)
				}

				if coding.Weight < 0 {
					coding.Weight = 0
				} else if coding.Weight > 1 {
					coding.Weight = 1
				}
			} else {
				if coding.Attributes == nil {
					coding.Attributes = make(map[string]string)
				}
				coding.Attributes[p0] = p1
			}

		default:
			return coding, fmt.Errorf("Cannot parse coding in '%s'; too many values", s)
		}
	}

	return coding, nil
}
