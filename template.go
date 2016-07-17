package mustache

// Template is a parsed template
type Template struct {
	tmpl []byte
	tkns tokens
}

// Render will render a template with the provided data
func (t *Template) Render(data interface{}, fn func([]byte)) (err error) {
	var a Aficionado
	switch nd := data.(type) {
	case Aficionado:
		a = nd
	case map[string]string:
		a = StringMap(nd)
	case map[string]interface{}:
		a = InterfaceMap(nd)

	default:
		err = ErrUnsupportedType
		return
	}

	r := newRenderer(t, a)
	if err = a.MarshalMustache(r); err != nil {
		return
	}

	fn(r.buf.Bytes())
	bp.Put(r.buf)
	r.buf = nil
	return
}
