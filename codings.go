package accept

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Coding combines a name (such as "gzip") with a possible value or a quality value. The quality
// value, if present, is a number between 0 and 1, inclusive.
type Coding struct {
	Name       string
	QValue     float64
	Attributes map[string]string
}

// IsAccepted returns true if the quality is greater than zero.
func (c Coding) IsAccepted() bool {
	return c.QValue > 0
}

// IsIdentity returns true if the name is "identity".
func (c Coding) IsIdentity() bool {
	return c.Name == Identity
}

// String returns the string representation of this coding.
func (c Coding) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString(c.Name)

	for k, v := range c.Attributes {
		buf.WriteByte(';')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}

	if 0 <= c.QValue && c.QValue < 1 {
		// no more than three decimal places are allowed (https://tools.ietf.org/html/rfc7231#section-5.3.1)
		fmt.Fprintf(buf, ";q=%.3g", c.QValue)
	}

	return buf.String()
}

func (c Coding) biasedQValue() float64 {
	// note that quality values have 3 decimal places accuracy and the biases rely on this
	quality := c.QValue
	if strings.HasPrefix(c.Name, "*") {
		// "*" 0r "*/*" has much less quality
		quality -= 0.0002
	} else if strings.HasSuffix(c.Name, "*") {
		// "text/*" has less quality
		quality -= 0.0001
	}
	return quality
}

//-------------------------------------------------------------------------------------------------

// Codings holds a list of codings.
type Codings []Coding

// Like finds codings that have names beginning with a given prefix or with a "*". Codings with zero
// quality are not accepted; the result will not contain any of these.
func (cs Codings) Like(prefix string) Codings {
	result := make(Codings, 0)
	for _, v := range cs {
		if v.QValue > 0 && (strings.HasPrefix(v.Name, prefix) || strings.HasPrefix(v.Name, "*")) {
			result = append(result, v)
		}
	}
	return result
}

// Get finds a named coding. If not found, it returns the zero value, which is never 'accepted'.
func (cs Codings) Get(name string) Coding {
	for _, v := range cs {
		if v.Name == name {
			return v
		}
	}
	return Coding{}
}

// Accepts tests whether a named coding is 'accepted'.
func (cs Codings) Accepts(name string) bool {
	return cs.Get(name).IsAccepted()
}

// Sorted sorts the codings by quality factor, highest first. Returns cs, which has been sorted.
// After sorting, the first items in the list are the most preferred. Sorting also takes into
// account "*" wildcards (these diminish the quality), and any attributes that make a coding more
// specific (these boost the quality).
func (cs Codings) Sorted() Codings {
	sort.Slice(cs, func(i, j int) bool {
		wi := cs[i].biasedQValue()
		wj := cs[j].biasedQValue()
		if wi == wj {
			// more attributes means more specificity
			return len(cs[i].Attributes) > len(cs[j].Attributes)
		}
		return wi > wj
	})
	return cs
}

// IfAccepted returns only the codings that are accepted (i.e. non-zero quality).
func (cs Codings) IfAccepted() Codings {
	accepted := make(Codings, 0, len(cs))
	for _, c := range cs {
		if c.QValue > 0 {
			accepted = append(accepted, c)
		}
	}
	return accepted
}

// Names returns the like of names from the list of codings.
func (cs Codings) Names() []string {
	str := make([]string, 0, len(cs))
	for _, c := range cs {
		str = append(str, c.Name)
	}
	return str
}

func (cs Codings) String() string {
	str := make([]string, 0, len(cs))
	for _, c := range cs {
		str = append(str, c.String())
	}
	return strings.Join(str, ", ")
}
