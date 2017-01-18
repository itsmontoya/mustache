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
	var (
		s       section
		ok      bool
		invalid bool
	)

	if s, ok, invalid = getSection(nil, data); invalid {
		return ErrUnsupportedType
	} else if !ok {
		if data != nil {
			return
		}

		s = nil
	}

	buf := t.bp.Get()

	switch st := s.(type) {
	case Aficionado:
		err = t.render(st, buf)
	case nil:
		err = t.renderList(nil, buf)
	case []Aficionado:
		err = t.renderList(st, buf)
	default:
		err = ErrUnsupportedType
		goto END
	}

	fn(buf.Bytes())

END:
	t.bp.Put(buf)
	return
}

// Render will render a template with the provided data
func (t *Template) render(a Aficionado, buf *buffer.Buffer) (err error) {
	r := rp.Get()
	r.t = t
	r.buf = buf
	r.a = a

	if err = a.MarshalMustache(r); err != nil {
		return
	}

	r.t = nil
	r.buf = nil
	r.a = nil
	r.get = nil

	rp.Put(r)
	return
}

// Render will render a template with the provided data
func (t *Template) renderList(as []Aficionado, buf *buffer.Buffer) (err error) {
	r := rp.Get()
	r.t = t
	r.buf = buf
	r.as = as

	r.render()

	r.t = nil
	r.buf = nil
	r.as = nil
	r.get = nil

	rp.Put(r)
	return
}
