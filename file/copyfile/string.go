package copyfile

import (
	"bytes"
	"io"
)

type ParsedString struct {
	str    string
	fields map[string]parsedStringField
}

const (
	openBrace  = '{'
	closeBrace = '}'
)

func Parse(str string) *ParsedString {
	ret := &ParsedString{
		str:    str,
		fields: make(map[string]parsedStringField),
	}
	ret.parse()
	return ret
}

func (s *ParsedString) parse() {
	r := newParser(s.str)

	isOpen := false
	start := 0
	idx := 0
	var paramName bytes.Buffer

	for {
		ch, ok := r.next()
		if !ok {
			break
		}
		isWriteChar := true
		switch {
		case ch == openBrace:
			if !isOpen {
				// check for escaping
				nch, ok := r.next()
				if !ok || nch != openBrace {
					if ok {
						r.unread()
					}
					isOpen = true
					start = idx
					paramName.Reset()
					isWriteChar = false
				} else {
					idx++
				}
			}
		case ch == closeBrace:
			if isOpen {
				// check for escaping
				nch, ok := r.next()
				if !ok || nch != closeBrace {
					if ok {
						r.unread()
					}
					s.fields[paramName.String()] = parsedStringField{
						name:  paramName.String(),
						start: start,
						end:   idx + 1,
					}
					isOpen = false
					isWriteChar = false
				} else {
					idx++
				}
			}
		}
		if isWriteChar && isOpen {
			_, _ = paramName.WriteRune(ch)
		}
		idx++
	}
}

type parsedStringField struct {
	name  string
	start int
	end   int
}

type parser struct {
	r *bytes.Reader
}

func newParser(str string) parser {
	return parser{r: bytes.NewReader([]byte(str))}
}

func (p *parser) next() (rune, bool) {
	ch, _, err := p.r.ReadRune()
	if err != nil {
		if err != io.EOF {
			panic(err)
		}
		return 0, false
	}
	return ch, true
}

func (p *parser) unread() {
	err := p.r.UnreadRune()
	if err != nil {
		panic(err)
	}
}
