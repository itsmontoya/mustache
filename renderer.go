package mustache

import (
	"sync"

	"github.com/itsmontoya/buffer"
	"github.com/itsmontoya/escapist"
)

var rp = rendererPool{
	p: sync.Pool{
		New: func() interface{} {
			return &Renderer{}
		},
	},
}

type rendererPool struct {
	p sync.Pool
}

func (rp *rendererPool) Get() *Renderer {
	return rp.p.Get().(*Renderer)
}

func (rp *rendererPool) Put(r *Renderer) {
	rp.p.Put(r)
}

func newRenderer(t *Template, a Aficionado) *Renderer {
	return &Renderer{
		t:   t,
		a:   a,
		buf: t.bp.Get(),
	}
}

// Renderer helps render
type Renderer struct {
	t *Template
	a Aficionado

	buf *buffer.Buffer
	get func(string) interface{}
}

func (r *Renderer) render() (err error) {
	for _, tkn := range r.t.tkns {
		switch tt := tkn.(type) {
		case tmplToken:
			r.buf.Write(r.t.tmpl[tt.start:tt.end])
		case valToken:
			if err = r.processValue(tt); err != nil {
				return
			}
		case sectionToken:
			if err = r.processSection(tt); err != nil {
				return
			}
		}
	}

	return
}

func (r *Renderer) processValue(tkn valToken) (err error) {
	var (
		b  []byte
		ok bool
	)

	if b, ok = getValueBytes(r.get(tkn.key)); !ok {
		return ErrUnsupportedType
	} else if b == nil {
		return
	}

	if tkn.escape {
		b = escapist.Escape(b)
	}

	r.buf.Write(b)
	return
}

func (r *Renderer) processSection(tkn sectionToken) (err error) {
	var (
		as []Aficionado
		ok bool
	)

	if tkn.key == "." {
		as = []Aficionado{r.a}
		goto LOOP
	}

	if as, ok = getAficionados(r.a, r.get(tkn.key)); !ok {
		return ErrUnsupportedType
	}

LOOP:
	if err = tkn.t.render(as, r.buf); err != nil {
		return
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
