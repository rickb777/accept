// Package accept provides parsing functions for the HTTP Accept-... headers, which are vital
// for HTTP content negotiation.
//
// Content negotiation headers use a relatively-complex syntax that permits a
// wide range of header from the very simple, e.g.
//
//     Accept-Encoding: gzip
//     Accept-Language: de
//     Accept: *
//
// to much more complex, e.g.
//
//     Accept-Encoding: compress;q=0.5, gzip;q=1.0
//     Accept-Language: de; q=1.0, en; q=0.5
//     Accept: text/html; q=1.0, text/*; q=0.8, image/gif; q=0.6, image/jpeg; q=0.6, image/*; q=0.5, */*; q=0.1
//
// See https://tools.ietf.org/html/rfc7231#section-5.3
package accept
