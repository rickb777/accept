package accept

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
)

// Coding combines a name (such as "gzip") with a possible value or a quality factor. The quality
// factor, if present, is a number between 0 and 1, inclusive.
type Coding struct {
	Name       string
	Weight     float64
	Attributes map[string]string
}

func (c Coding) IsAccepted() bool {
	return c.Weight > 0
}

func (c Coding) IsIdentity() bool {
	return c.Weight > 0
}

func (c Coding) String() string {
	if c.Name == "" {
		return ""
	} else if c.Weight == 0 && len(c.Attributes) == 0 {
		return ""
	}

	buf := &bytes.Buffer{}
	buf.WriteString(c.Name)

	for k, v := range c.Attributes {
		buf.WriteByte(';')
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(v)
	}

	if 0 < c.Weight && c.Weight < 1 {
		// no more than three decimal places are allowed (https://tools.ietf.org/html/rfc7231#section-5.3.1)
		fmt.Fprintf(buf, ";q=%.3g", c.Weight)
	}

	return buf.String()
}

func (c Coding) biasedWeight() float64 {
	// note that the biases rely on the 3 decimal places accuracy of the q weight values
	weight := c.Weight
	if strings.HasPrefix(c.Name, "*") {
		// "*" 0r "*/*" has much less weight
		weight -= 0.0002
	} else if strings.HasSuffix(c.Name, "*") {
		// "text/*" has less weight
		weight -= 0.0001
	}
	return weight
}

//-------------------------------------------------------------------------------------------------

// Codings holds a list of codings.
type Codings []Coding

// Like finds codings that have names beginning with a given prefix.
func (cs Codings) Like(prefix string) Codings {
	result := make(Codings, 0)
	for _, v := range cs {
		if v.Weight > 0 && (strings.HasPrefix(v.Name, prefix) || strings.HasPrefix(v.Name, "*")) {
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

// Sorted sorts the codings by quality factor, highest first. Returns cs, which has been sorted.
// After sorting, the first items in the list are the most preferred. Sorting also takes into
// account "*" wildcards (these diminish the weight), and any attributes that make a coding more
// specific (these boost the weight).
func (cs Codings) Sorted() Codings {
	sort.Slice(cs, func(i, j int) bool {
		wi := cs[i].biasedWeight()
		wj := cs[j].biasedWeight()
		if wi == wj {
			// more attributes means more specificity
			return len(cs[i].Attributes) > len(cs[j].Attributes)
		}
		return wi > wj
	})
	return cs
}

func (cs Codings) String() string {
	str := make([]string, 0, len(cs))
	for _, c := range cs {
		str = append(str, c.String())
	}
	return strings.Join(str, ", ")
}
