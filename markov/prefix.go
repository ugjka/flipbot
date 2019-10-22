package markov

import "strings"

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

func NewPrefix(key string, prefixLen int) *Prefix {
	p := make(Prefix, prefixLen)
	copy(p, strings.SplitN(key, " ", prefixLen))
	return &p
}

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}
