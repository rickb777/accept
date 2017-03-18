// Package accept provides parsing functions for the HTTP Accept-... headers, which are vital
// for HTTP content negotiation. The most widely used example is gzip content encoding.
//
// Content negotiation headers use a relatively-complex syntax that permits a
// wide range of header from the very simple (e.g. "Accept-Encoding: gzip") to much
// more complex (e.g. "Accept-Encoding: compress;q=0.5, gzip;q=1.0").
//
// See https://tools.ietf.org/html/rfc7231#section-5.3
package accept
