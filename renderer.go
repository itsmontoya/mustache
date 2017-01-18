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
	t     *Template
	a     Aficionado
	as    []Aficionado
	isSet bool

	buf *buffer.Buffer
	get func(string) interface{}
}

func (r *Renderer) render() (err error) {
	for _, tkn := range r.t.tkns {
		switch tt := tkn.(type) {
		case tmplToken:
			r.buf.Write(r.t.tmpl[tt.start:tt.end])
		case valToken:
			err = r.processValue(tt)
		case sectionToken:
			err = r.processSection(tt)
		case invertedSectionToken:
			err = r.processInvertedSection(tt)
		}

		if err != nil {
			break
		}
	}

	return
}

func (r *Renderer) processValue(tkn valToken) (err error) {
	if r.a == nil {
		return
	}

	if b, ok, invalid := getValueBytes(r.get(tkn.key)); invalid {
		return ErrUnsupportedType
	} else if !ok {
		return
	} else {
		if tkn.escape {
			b = escapist.Escape(b)
		}

		r.buf.Write(b)
	}

	return
}

func (r *Renderer) processSection(tkn sectionToken) (err error) {
	var (
		s       section
		ok      bool
		invalid bool
	)

	if r.as != nil {
		if tkn.key == "." {
			s = r.as
		} else {
			err = ErrUnsupportedType
			return
		}
	} else {
		var v interface{}
		if tkn.key == "." {
			v = r.a != nil
		} else {
			v = r.get(tkn.key)
		}

		if s, ok, invalid = getSection(r.a, v); invalid {
			return ErrUnsupportedType
		} else if !ok {
			return
		}
	}

	switch st := s.(type) {
	case Aficionado:
		err = tkn.t.render(st, r.buf)
	case []Aficionado:
		for _, a := range st {
			err = tkn.t.render(a, r.buf)
		}

	case nil:

	default:
		err = ErrUnsupportedType
	}

	return
}

func (r *Renderer) processInvertedSection(tkn invertedSectionToken) (err error) {
	var (
		v       interface{}
		s       section
		ok      bool
		invalid bool
	)

	if r.a != nil {
		if tkn.key == "." {
			v = r.a != nil
		} else {
			v = r.get(tkn.key)
		}
	} else {
		if tkn.key != "." {
			err = ErrUnsupportedType
			return
		}

		v = r.as
	}

	if s, ok, invalid = getInvertedSection(r.a, v); invalid {
		return ErrUnsupportedType
	} else if !ok {
		return
	}

	switch st := s.(type) {
	case Aficionado:
		err = tkn.t.render(st, r.buf)
	case []Aficionado, nil:
		err = tkn.t.renderList(nil, r.buf)
		//	case nil:

	default:
		err = ErrUnsupportedType
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

type section interface{}
