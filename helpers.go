package mustache

import "strconv"

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

type Value struct {
	v interface{}
}

func (val Value) MarshalMustache(r *Renderer) (err error) {
	return r.ForEach(func(k string) (v interface{}) {
		switch k {
		case ".":
			v = val.v
		}
		return
	})
}

func getValueBytes(v interface{}) (b []byte, ok, invalid bool) {
	ok = true
	switch nv := v.(type) {
	case []byte:
		b = nv
	case string:
		b = []byte(nv)

	case int64:
		b = strconv.AppendInt(b, nv, 10)
	case int32:
		b = strconv.AppendInt(b, int64(nv), 10)
	case int16:
		b = strconv.AppendInt(b, int64(nv), 10)
	case int8:
		b = strconv.AppendInt(b, int64(nv), 10)
	case int:
		b = strconv.AppendInt(b, int64(nv), 10)

	case uint64:
		b = strconv.AppendUint(b, nv, 10)
	case uint32:
		b = strconv.AppendUint(b, uint64(nv), 10)
	case uint16:
		b = strconv.AppendUint(b, uint64(nv), 10)
	case uint8:
		b = strconv.AppendUint(b, uint64(nv), 10)
	case uint:
		b = strconv.AppendUint(b, uint64(nv), 10)

	case float64:
		b = strconv.AppendFloat(b, nv, 'f', -1, 64)
	case float32:
		b = strconv.AppendFloat(b, float64(nv), 'f', -1, 32)

	case bool:
		b = strconv.AppendBool(b, nv)
	case nil:
		ok = false

	default:
		ok = false
		invalid = true
	}

	return
}

func getAficionado(pa Aficionado, v interface{}) (a Aficionado, ok, invalid bool) {
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

	case nil:
		ok = false

	default:
		ok = false
		invalid = true
	}

	return
}

func getSection(pa Aficionado, v interface{}) (s section, ok, invalid bool) {
	if v == nil {
		return
	}

	var a Aficionado
	if a, ok, _ = getAficionado(pa, v); ok {
		s = a
		return
	}

	ok = true
	switch nv := v.(type) {
	case []Aficionado:
		s = nv
	case []interface{}:
		as := make([]Aficionado, len(nv))
		for k, v := range nv {
			if a, ok, invalid = getAficionado(pa, v); !ok || invalid {
				return
			}

			as[k] = a
		}

		s = as
	case []byte:
		s = ByteSlice(nv).Values()
	case []string:
		s = StringSlice(nv).Values()
	case []int64:
		s = Int64Slice(nv).Values()
	case []int32:
		s = Int32Slice(nv).Values()
	case []int:
		s = IntSlice(nv).Values()
	case []float64:
		s = Float64Slice(nv).Values()
	case []float32:
		s = Float32Slice(nv).Values()

	case nil:
		ok = false

	default:
		ok = false
		invalid = true
	}

	return

}

type ByteSlice []byte

func (s ByteSlice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type StringSlice []string

func (s StringSlice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type Int64Slice []int64

func (s Int64Slice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type Int32Slice []int32

func (s Int32Slice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type IntSlice []int

func (s IntSlice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type Float64Slice []float64

func (s Float64Slice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

type Float32Slice []float32

func (s Float32Slice) Values() (as []Aficionado) {
	as = make([]Aficionado, len(s))
	for k, v := range s {
		as[k] = Value{v}
	}

	return
}

func getInvertedSection(pa Aficionado, v interface{}) (s section, ok, invalid bool) {
	switch nv := v.(type) {
	case bool:
		if !nv {
			ok = true
			s = pa
		}
	case string:
		if len(nv) == 0 {
			ok = true
			s = pa
		}

	case map[string]string:
		if len(nv) == 0 {
			ok = true
			s = pa
		}
	case map[string]interface{}:
		if len(nv) == 0 {
			ok = true
			s = pa
		}
	case map[string][]byte:
		if len(nv) == 0 {
			ok = true
			s = pa
		}

	case []string:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []int64:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []int32:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []int:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []float64:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []float32:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []byte:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []interface{}:
		if len(nv) == 0 {
			s = pa
			ok = true
		}
	case []Aficionado:
		if len(nv) == 0 {
			s = pa
			ok = true
		}

	case Aficionado:
	case nil:
		ok = true

	default:
		invalid = true
	}

	return
}

func isChar(b byte) bool {
	return (b >= lwrCaseStart && b <= lwrCaseEnd) || (b >= uprCaseStart && b <= uprCaseEnd)
}

func isWhiteSpace(b byte) bool {
	return b == charSpace || b == charNewline || b == charTab
}
