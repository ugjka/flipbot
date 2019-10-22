package markov

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"strings"
)

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	roots     []string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), []string{}, prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}

		key := p.String()

		startsSentence := strings.ToLower(p[0]) != p[0]
		chapterTitle := strings.ToUpper(p[0]) == p[0]
		if startsSentence && !chapterTitle {
			c.roots = append(c.roots, key)
		}

		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
//func (c *Chain) Generate(w io.Writer, n int) {
func (c *Chain) Generate(n int) string {
	p := c.randomRoot()
	//bw := bufio.NewWriter(w)
	//defer bw.Flush()
	genText := ""
	genText = genText + p.String() + " "
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			p = c.randomRoot()
		}
		next := choices[rand.Intn(len(choices))]
		p.Shift(next)
		genText = genText + next + " "
		//bw.WriteString(next)
		//bw.WriteString(" ")
	}
	if num := strings.LastIndex(genText, "."); num != -1 {
		return genText[0:num + 1]
		}                                        // Build chains from standard input.
	return genText
}

func (c *Chain) randomRoot() *Prefix {
	key := c.roots[rand.Intn(len(c.roots))]
	p := NewPrefix(key, c.prefixLen)
	return p
}
