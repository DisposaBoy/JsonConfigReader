package JsonConfigReader

import (
	"io"
)

type state struct {
	escapeC   int
	inString  bool
	inComment bool
	r         io.Reader
}

// Read acts as a proxy for the underlying reader and cleans p
// of comments and trailing commas preceeding ] and }
// comments are delimitted by // up until the end the line
func (s *state) Read(p []byte) (n int, err error) {
	n, err = s.r.Read(p)
	if err == nil {
		lpcOff := -1
		var lc, lpc byte
		for i, c := range p[0:n] {
			if s.inString {
				if s.escapeC == 0 {
					if c == '"' {
						s.inString = false
					} else if c == '\\' {
						s.escapeC += 1
					}
				} else {
					s.escapeC = 0
				}
			} else {
				if s.inComment {
					if c == '\n' || c == '\r' {
						s.inComment = false
					}
				} else {
					if lpc == ',' && (c == '}' || c == ']') {
						p[lpcOff] = ' '
					}
					if c != ' ' && c != '\t' && c != '\n' && c != '\r' && c != '/' {
						lpc = c
						lpcOff = i
					}

					if c == '"' {
						s.inString = true
					} else if c == '/' && lc == '/' {
						// don't clear a single /, let a syntax error occur
						p[i-1] = ' '
						s.inComment = true
					}
				}

				if s.inComment {
					p[i] = ' '
				}
			}
			lc = c
		}
	}
	return
}

// New returns an io.Reader acting as proxy to r
func New(r io.Reader) io.Reader {
	return &state{r: r}
}
