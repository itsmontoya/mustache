package mustache

import (
	"io"
	"os"
	"path"

	"github.com/itsmontoya/buffer"
	"github.com/missionMeteora/toolkit/errors"
)

const (
	charLCurly      = '{'
	charRCurly      = '}'
	charPound       = '#'
	charFSlash      = '/'
	charSpace       = ' '
	charNewline     = '\n'
	charTab         = '\t'
	charPeriod      = '.'
	charCarrot      = '^'
	charGreaterThan = '>'

	lwrCaseStart = 'a'
	lwrCaseEnd   = 'z'
	uprCaseStart = 'A'
	uprCaseEnd   = 'Z'
)

const (
	stateRootStart uint8 = iota
	stateContainerStart
	stateContainerOpen
	stateContainerEnd

	stateValueOpen
	stateValueEnd
	stateValueClosing
	stateValueClosed

	stateUnescapedValueStart
	stateUnescapedValueOpen
	stateUnescapedValueEnd
	stateUnescapedValueClosingA
	stateUnescapedValueClosingB
	stateUnescapedValueClosed

	stateSectionStart
	stateSectionOpen
	stateSectionEnd
	stateSectionClosing
	stateSectionClosed

	stateInvertedSectionStart
	stateInvertedSectionOpen
	stateInvertedSectionEnd
	stateInvertedSectionClosing
	stateInvertedSectionClosed

	stateTmplStart
	stateTmplOpen
	stateTmplEnd
	stateTmplClosing

	stateRootEnd

	stateError
)

const (
	// ErrInvalidSyntax is returned when syntax is invalid
	ErrInvalidSyntax = errors.Error("invalid syntax")

	// ErrForEachSet is returned when ForEach is called more than once for a particular parser
	ErrForEachSet = errors.Error("ForEach has already been called for this parser")

	// ErrUnsupportedType is returned when an upsupported type is provided
	ErrUnsupportedType = errors.Error("unsupported type provided")
)

var bp = buffer.NewPool(32)

//var bp = newPool()

// Render is a one-shot render
func Render(a, b, c interface{}) error {
	return nil
}

// Parse will parse a byteslice template and return a mustache Template
func Parse(tmpl []byte, filePath string) (t *Template, err error) {
	var tkns tokens
	if tkns, err = parse(tmpl, filePath); err != nil {
		return
	}

	t = newTemplate(tmpl, tkns)
	return
}

func parse(tmpl []byte, fp string) (tkns tokens, err error) {
	p := parser{
		kbuf: bp.Get(),
		tmpl: tmpl,
		fp:   fp,
	}

	if err = p.parse(); err != nil {
		return
	}

	tkns = p.tkns
	return
}

type parser struct {
	kbuf *buffer.Buffer

	tmpl  []byte
	state uint8

	idx    int
	start  int
	kstart int

	tkns tokens

	fp string // Filepath
}

func (p *parser) parse() (err error) {
	var v byte
	for ; p.idx < len(p.tmpl); p.idx++ {
		v = p.tmpl[p.idx]
		switch p.state {
		case stateRootStart:
			p.rootStart(v)

		case stateContainerStart:
			p.containerStart(v)
		case stateContainerOpen:
			p.containerOpen(v)

		case stateValueOpen:
			p.valueOpen(v)
		case stateValueEnd:
			p.valueEnd(v)
		case stateValueClosing:
			p.valueClosing(v, true)

		case stateUnescapedValueStart:
			p.unescapedValueStart(v)
		case stateUnescapedValueOpen:
			p.unescapedValueOpen(v)
		case stateUnescapedValueEnd:
			p.unescapedValueEnd(v)
		case stateUnescapedValueClosingA:
			p.unescapedValueClosing(v)
		case stateUnescapedValueClosingB:
			p.valueClosing(v, false)

		case stateSectionStart:
			p.sectionStart(v)
		case stateSectionOpen:
			p.sectionOpen(v)
		case stateSectionEnd:
			p.sectionEnd(v)
		case stateSectionClosing:
			p.sectionClosing(v)

		case stateInvertedSectionStart:
			p.invertedSectionStart(v)
		case stateInvertedSectionOpen:
			p.invertedSectionOpen(v)
		case stateInvertedSectionEnd:
			p.invertedSectionEnd(v)
		case stateInvertedSectionClosing:
			p.invertedSectionClosing(v)

		case stateTmplStart:
			p.tmplStart(v)
		case stateTmplOpen:
			p.tmplOpen(v)
		case stateTmplEnd:
			p.tmplEnd(v)
		case stateTmplClosing:
			p.tmplClosing(v)

		case stateRootEnd:
			break

		case stateError:
			err = ErrInvalidSyntax
			goto END
		}
	}

	if p.start > -1 {
		p.tkns = append(p.tkns, tmplToken{
			start: p.start,
			end:   len(p.tmpl),
		})
	}

END:
	bp.Put(p.kbuf)
	p.kbuf = nil
	p.tmpl = nil
	return
}

func (p *parser) rootStart(b byte) {
	if p.start == -1 {
		p.start = p.idx
	}

	if b != charLCurly {
		return
	}

	p.tkns = append(p.tkns, tmplToken{
		start: p.start,
		end:   p.idx,
	})
	p.state = stateContainerStart
}

func (p *parser) containerStart(b byte) {
	if b == charLCurly {
		p.state = stateContainerOpen
	} else {
		p.state = stateRootStart
	}
}

func (p *parser) containerOpen(b byte) {
	switch {
	case isWhiteSpace(b):
	case isChar(b), b == charPeriod:
		p.kstart = p.idx
		p.state = stateValueOpen
	case b == charLCurly:
		p.state = stateUnescapedValueStart
	case b == charPound:
		p.state = stateSectionStart
	case b == charCarrot:
		p.state = stateInvertedSectionStart
	case b == charGreaterThan:
		p.state = stateTmplStart

	default:
		p.state = stateError
	}
}

func (p *parser) valueOpen(b byte) {
	switch {
	case isChar(b):
	case b == charPeriod:
	case isWhiteSpace(b):
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateValueEnd
	case b == charRCurly:
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateValueClosing
	default:
		p.state = stateError
	}
}

func (p *parser) valueEnd(b byte) {
	switch {
	case b == charRCurly:
		p.state = stateValueClosing
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) valueClosing(b byte, escape bool) {
	if b != charRCurly {
		p.state = stateError
		return
	}

	p.tkns = append(p.tkns, valToken{
		key:    p.kbuf.String(),
		escape: escape,
	})

	p.kbuf.Reset()
	p.start = -1
	p.kstart = -1
	p.state = stateRootStart
}

func (p *parser) unescapedValueStart(b byte) {
	switch {
	case isWhiteSpace(b):
	case isChar(b):
		p.kstart = p.idx
		p.state = stateUnescapedValueOpen
	default:
		p.state = stateError
	}
}

func (p *parser) unescapedValueOpen(b byte) {
	switch {
	case isChar(b):
	case b == charPeriod:
	case isWhiteSpace(b):
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateUnescapedValueEnd
	case b == charRCurly:
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateUnescapedValueClosingA
	default:
		p.state = stateError
	}
}

func (p *parser) unescapedValueEnd(b byte) {
	switch {
	case b == charRCurly:
		p.state = stateUnescapedValueClosingA
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) unescapedValueClosing(b byte) {
	if b == charRCurly {
		p.state = stateUnescapedValueClosingB
	} else {
		p.state = stateError
	}
}

func (p *parser) sectionStart(b byte) {
	switch {
	case isChar(b), b == charPeriod:
		p.state = stateSectionOpen
		p.kstart = p.idx
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) sectionOpen(b byte) {
	switch {
	case isChar(b):
	case b == charPeriod:
	case isWhiteSpace(b):
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateSectionEnd
	case b == charRCurly:
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateSectionClosing
	default:
		p.state = stateError
	}
}

func (p *parser) sectionEnd(b byte) {
	switch {
	case b == charRCurly:
		p.state = stateSectionClosing
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) sectionClosing(b byte) {
	if b != charRCurly {
		p.state = stateError
		return
	}

	p.idx++
	var ss, se int
	if ss, se = findSectionEnd(p.tmpl[p.idx:]); se == -1 {
		p.state = stateError
		return
	}

	var (
		st  *Template
		err error
	)

	if st, err = Parse(p.tmpl[p.idx:p.idx+ss], p.fp); err != nil {
		p.state = stateError
		return
	}

	p.tkns = append(p.tkns, sectionToken{
		key: p.kbuf.String(),
		t:   st,
	})
	p.idx += se

	p.kbuf.Reset()
	p.start = -1
	p.kstart = -1
	p.state = stateRootStart
}

func (p *parser) invertedSectionStart(b byte) {
	switch {
	case isChar(b), b == charPeriod:
		p.state = stateInvertedSectionOpen
		p.kstart = p.idx
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) invertedSectionOpen(b byte) {
	switch {
	case isChar(b):
	case b == charPeriod:
	case isWhiteSpace(b):
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateInvertedSectionEnd
	case b == charRCurly:
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateInvertedSectionClosing
	default:
		p.state = stateError
	}
}

func (p *parser) invertedSectionEnd(b byte) {
	switch {
	case b == charRCurly:
		p.state = stateInvertedSectionClosing
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) invertedSectionClosing(b byte) {
	if b != charRCurly {
		p.state = stateError
		return
	}

	p.idx++
	var ss, se int
	if ss, se = findSectionEnd(p.tmpl[p.idx:]); se == -1 {
		p.state = stateError
		return
	}

	var (
		st  *Template
		err error
	)

	if st, err = Parse(p.tmpl[p.idx:p.idx+ss], p.fp); err != nil {
		p.state = stateError
		return
	}

	p.tkns = append(p.tkns, invertedSectionToken{
		key: p.kbuf.String(),
		t:   st,
	})
	p.idx += se

	p.kbuf.Reset()
	p.start = -1
	p.kstart = -1
	p.state = stateRootStart
}

func (p *parser) tmplStart(b byte) {
	switch {
	case isChar(b), b == charPeriod:
		p.state = stateTmplOpen
		p.kstart = p.idx
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) tmplOpen(b byte) {
	switch {
	case isChar(b), b == charPeriod, b == charFSlash:
	case isWhiteSpace(b):
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateTmplEnd
	case b == charRCurly:
		p.kbuf.Write(p.tmpl[p.kstart:p.idx])
		p.state = stateTmplClosing
	default:
		p.state = stateError
	}
}

func (p *parser) tmplEnd(b byte) {
	switch {
	case b == charRCurly:
		p.state = stateTmplClosing
	case isWhiteSpace(b):
	default:
		p.state = stateError
	}
}

func (p *parser) tmplClosing(b byte) {
	if b != charRCurly {
		p.state = stateError
		return
	}

	var (
		f   *os.File
		st  sectionToken
		err error

		buf = bp.Get() // We are not going to return this to the pool until we copy the bytes properly
	)

	if f, err = os.Open(path.Join(p.fp, p.kbuf.String())); err != nil {
		p.state = stateError
		goto END
	}

	io.Copy(buf, f)

	st.key = "."
	if st.t, err = Parse(buf.Bytes(), p.fp); err != nil {
		p.state = stateError
		goto END
	}

	p.tkns = append(p.tkns, st)

END:
	if f != nil {
		f.Close()
	}

	p.kbuf.Reset()
	p.start = -1
	p.kstart = -1
	p.state = stateRootStart
}

func findSectionEnd(in []byte) (start, i int) {
	var (
		b     byte
		state uint8
		level uint8 = 1
	)

	for i, b = range in {
		switch state {
		case 0:
			if b == charLCurly {
				start = i
				state = 1
			}
		case 1:
			if b == charLCurly {
				state = 2
			} else {
				state = 0
			}
		case 2:
			if b == charPound || b == charCarrot {
				level++
				state = 0
			} else if b == charFSlash {
				state = 3
			} else if !isWhiteSpace(b) {
				state = 0
			}

		case 3:
			if isChar(b) || b == charPeriod {
				state = 4
			} else if !isWhiteSpace(b) {
				state = 0
			}

		case 4:
			if isWhiteSpace(b) {
				state = 5
			} else if b == charRCurly {
				state = 6
			} else if !isChar(b) {
				state = 0
			}
		case 5:
			if b == charRCurly {
				state = 6
			} else if !isWhiteSpace(b) {
				state = 0
			}
		case 6:
			level--
			state = 0

			if level == 0 {
				//	i++
				return
			}
		}
	}

	return -1, -1
}
