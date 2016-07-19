package mustache

import "strconv"

func isChar(b byte) bool {
	return (b >= lwrCaseStart && b <= lwrCaseEnd) || (b >= uprCaseStart && b <= uprCaseEnd)
}

func isWhiteSpace(b byte) bool {
	return b == charSpace || b == charNewline || b == charTab
}

// Aficionado is someone who really appreciates the Mustache
type Aficionado interface {
	MarshalMustache(*Renderer) error
}

// StringMap is a common map[string]string, has the func needed to be an Aficionado
type StringMap map[string]string

// MarshalMustache is what makes us one of the best, baby!
func (m StringMap) MarshalMustache(r *Renderer) (err error) {
	r.ForEach(m.Get)
	return
}

// Get will get a value by key
func (m StringMap) Get(key string) (val interface{}) {
	return []byte(m[key])
}

// InterfaceMap is a common map[string]string, has the func needed to be an Aficionado
type InterfaceMap map[string]interface{}

// MarshalMustache is what makes us one of the best, baby!
func (m InterfaceMap) MarshalMustache(r *Renderer) (err error) {
	r.ForEach(m.Get)
	return
}

// Get will get a value by key
func (m InterfaceMap) Get(key string) (val interface{}) {
	return m[key]
}

// BytesMap is a common map[string][]byte, has the func needed to be an Aficionado
type BytesMap map[string][]byte

// MarshalMustache is what makes us one of the best, baby!
func (m BytesMap) MarshalMustache(r *Renderer) (err error) {
	r.ForEach(m.Get)
	return
}

// Get will get a value by key
func (m BytesMap) Get(key string) (val interface{}) {
	return m[key]
}

func getValueBytes(v interface{}) (b []byte, ok bool) {
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

func getAficionado(pa Aficionado, v interface{}) (a Aficionado, ok bool) {
	ok = true
	switch nv := v.(type) {
	case Aficionado:
		a = nv
	case bool:
		if nv {
			a = pa
		}
	case map[string]string:
		a = StringMap(nv)
	case map[string]interface{}:
		a = InterfaceMap(nv)
	case map[string][]byte:
		a = BytesMap(nv)

	default:
		ok = false
	}

	return
}

func getAficionados(pa Aficionado, v interface{}) (as []Aficionado, ok bool) {
	var a Aficionado
	if a, ok = getAficionado(pa, v); ok {
		as = []Aficionado{a}
		return
	}

	ok = true
	switch nv := v.(type) {
	case []Aficionado:
		as = nv
	case []interface{}:
		for _, i := range nv {
			if a, ok = getAficionado(pa, i); !ok {
				return
			}

			as = append(as, a)
		}

	default:
		ok = false
	}

	return
}
