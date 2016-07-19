package mustache

import "github.com/itsmontoya/buffer"

func newTemplate(tmpl []byte, tkns tokens) *Template {
	return &Template{
		tmpl: tmpl,
		tkns: tkns,
		bp:   bp,
	}
}

// Template is a parsed template
type Template struct {
	tmpl []byte
	tkns tokens

	bp *buffer.Pool
}

func (t *Template) setBaseLen(baseLen int) {
	if baseLen == -1 {
		baseLen = len(t.tmpl)
	}

	if v := baseLen * 130 / 100; v < 32 {
		t.bp = bp
	} else {
		t.bp = buffer.NewPool(v)
	}
}

// Render will render a template with the provided data
func (t *Template) Render(data interface{}, fn func([]byte)) (err error) {
	as, ok := getAficionados(nil, data)
	if as == nil || !ok {
		return ErrUnsupportedType
	}

	buf := t.bp.Get()
	t.render(as, buf)
	fn(buf.Bytes())
	t.bp.Put(buf)
	return
}

// Render will render a template with the provided data
func (t *Template) render(as []Aficionado, buf *buffer.Buffer) (err error) {
	r := rp.Get()
	r.t = t
	r.buf = buf

	for _, a := range as {
		r.a = a
		if err = a.MarshalMustache(r); err != nil {
			return
		}
		r.get = nil
	}

	r.t = nil
	r.a = nil
	r.buf = nil
	rp.Put(r)
	return
}
