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
// These functions provide easy-to-use wrappers around the main API.

// PreferredContentTypeLike finds the first accepted content type in a given header value that is like a given prefix.
// The result is an empty string if there is no match or there was a parse error.
func PreferredContentTypeLike(hdr http.Header, like string) string {
	return PreferredLike(hdr.Get(Accept), like)
}

// PreferredCharsetLike finds the first accepted charset in a given header value that is like a given prefix.
// The result is an empty string if there is no match or there was a parse error.
func PreferredCharsetLike(hdr http.Header, like string) string {
	return PreferredLike(hdr.Get(AcceptCharset), like)
}

// PreferredLanguageLike finds the first accepted language in a given header value that is like a given prefix.
// The result is an empty string if there is no match or there was a parse error.
func PreferredLanguageLike(hdr http.Header, like string) string {
	return PreferredLike(hdr.Get(AcceptLanguage), like)
}

// PreferredLike finds the first item in a given header value that is like a given prefix.
// The result is an empty string if there is no match or there was a parse error.
func PreferredLike(hdr, like string) string {
	acceptedContentTypes, err := Parse(hdr)
	if err != nil {
		return ""
	}
	preferred := acceptedContentTypes.Like(like)
	if len(preferred) == 0 {
		return ""
	}
	return preferred.Sorted()[0].Name
}

// AcceptsEncoding determines whether there is a request header that indicates that it will
// accept a particular coding for the response. The most used coding is "gzip".
func AcceptsEncoding(hdr http.Header, coding string) bool {
	acceptedEncodings, err := Parse(hdr.Get(AcceptEncoding))
	return err == nil && acceptedEncodings.Get(coding).IsAccepted()
}

//-------------------------------------------------------------------------------------------------

// Parse returns the codings from a given header value. They are returned in the order they
// appear; use the Sorted() method on the result if you need them in preference order instead.
func Parse(acceptValue string) (Codings, error) {
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
		QValue: 1,
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
				coding.QValue, err = strconv.ParseFloat(p1, 64)
				if err != nil {
					return coding, fmt.Errorf("Cannot parse q-value in '%s'; %v", s, err)
				}

				if coding.QValue < 0 {
					coding.QValue = 0
				} else if coding.QValue > 1 {
					coding.QValue = 1
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
