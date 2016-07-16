package mustache

import (
	"strconv"
)

func isChar(b byte) bool {
	return (b >= lwrCaseStart && b <= lwrCaseEnd) || (b >= uprCaseStart && b <= uprCaseEnd)
}

func isWhiteSpace(b byte) bool {
	return b == charSpace || b == charNewline || b == charTab
}

// Aficionado is someone who really appreciates the Mustache
type Aficionado interface {
	MarshalMustache(*parser) error
}

// StringMap is a common map[string]string, has the func needed to be an Aficionado
type StringMap map[string]string

// MarshalMustache is what makes us one of the best, baby!
func (m StringMap) MarshalMustache(p *parser) (err error) {
	p.ForEach(m.Get)
	return
}

// Get will get a value by key
func (m StringMap) Get(key string) (val interface{}) {
	return []byte(m[key])
}

// InterfaceMap is a common map[string]string, has the func needed to be an Aficionado
type InterfaceMap map[string]interface{}

// MarshalMustache is what makes us one of the best, baby!
func (m InterfaceMap) MarshalMustache(p *parser) (err error) {
	p.ForEach(m.Get)
	return
}

// Get will get a value by key
func (m InterfaceMap) Get(key string) (val interface{}) {
	return m[key]
}

func handleValue(v interface{}) (b []byte, ok bool) {
	ok = true
	switch nv := v.(type) {
	case []byte:
		b = nv
	case string:
		b = []byte(nv)
	case int64:
		b = strconv.AppendInt(b, nv, 10)
	case float64:
		b = strconv.AppendFloat(b, nv, 'f', -1, 64)
	case bool:
		b = strconv.AppendBool(b, nv)
	default:
		ok = false
	}

	return
}

func handleSection(a Aficionado, tmpl []byte, v interface{}) (bs []byte, ok bool) {
	ok = true
	switch nv := v.(type) {
	case Aficionado:
		parse(tmpl, nv, func(b []byte) {
			bs = make([]byte, len(b))
			copy(bs, b)
		})
	case bool:
		if nv {
			parse(tmpl, a, func(b []byte) {
				bs = make([]byte, len(b))
				copy(bs, b)
			})
		}
	case map[string]interface{}:
		parse(tmpl, InterfaceMap(nv), func(b []byte) {
			bs = make([]byte, len(b))
			copy(bs, b)
		})
	case []Aficionado:
		for _, vi := range nv {
			parse(tmpl, vi, func(b []byte) {
				bs = append(bs, b...)
			})
		}
	default:
		ok = false
	}

	return
}
