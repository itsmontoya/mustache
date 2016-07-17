package mustache

import (
	"bytes"
	"github.com/itsmontoya/escapist"
)

func newRenderer(t *Template, a Aficionado) *Renderer {
	return &Renderer{
		t:   t,
		a:   a,
		buf: bp.Get(),
	}
}

// Renderer helps render
type Renderer struct {
	t *Template
	a Aficionado

	buf *bytes.Buffer
	get func(string) interface{}
}

func (r *Renderer) render() (err error) {
	var (
		b  []byte
		ok bool
	)

	for _, tkn := range r.t.tkns {
		switch tt := tkn.(type) {
		case tmplToken:
			r.buf.Write(r.t.tmpl[tt.start:tt.end])
		case valToken:
			if b, ok = getValueBytes(r.get(tt.key)); !ok {
				return ErrUnsupportedType
			} else if b == nil {
				return
			}

			if tt.escape {
				b = escapist.Escape(b)
			}

			r.buf.Write(b)
		case sectionToken:
			r.processSection(tt)
		}
	}

	return
}

func (r *Renderer) processSection(tkn sectionToken) (err error) {
	var (
		as []Aficionado
		ok bool
	)

	if as, ok = getAficionados(r.a, r.get(tkn.key)); !ok {
		return ErrUnsupportedType
	}

	for _, a := range as {
		if err = tkn.t.Render(a, func(b []byte) {
			r.buf.Write(b)
		}); err != nil {
			return
		}
	}

	return
}

// ForEach takes in a get func
func (r *Renderer) ForEach(fn func(string) interface{}) (err error) {
	if r.get != nil {
		return ErrForEachSet
	}

	r.get = fn
	return r.render()
}
