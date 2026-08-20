package sitecontent

import "github.com/tinywasm/model"

type keyVal struct {
	key string
	val any
}

type kvMap []keyVal

func (m *kvMap) get(key string) (any, bool) {
	if m == nil {
		return nil, false
	}
	for i := 0; i < len(*m); i++ {
		if (*m)[i].key == key {
			return (*m)[i].val, true
		}
	}
	return nil, false
}

func (m *kvMap) set(key string, val any) {
	for i := 0; i < len(*m); i++ {
		if (*m)[i].key == key {
			(*m)[i].val = val
			return
		}
	}
	*m = append(*m, keyVal{key: key, val: val})
}

type memFieldWriter struct {
	data kvMap
}

func newMemWriter() *memFieldWriter {
	return &memFieldWriter{}
}

func (w *memFieldWriter) String(k, v string)        { w.data.set(k, v) }
func (w *memFieldWriter) Int(k string, v int64)     { w.data.set(k, v) }
func (w *memFieldWriter) Float(k string, v float64) { w.data.set(k, v) }
func (w *memFieldWriter) Bool(k string, v bool)     { w.data.set(k, v) }
func (w *memFieldWriter) Bytes(k string, v []byte)  { w.data.set(k, string(v)) }
func (w *memFieldWriter) Null(k string)             {}
func (w *memFieldWriter) Raw(k, v string)           { w.data.set(k, v) }
func (w *memFieldWriter) Object(k string, v model.Encodable) {
	if v == nil || v.IsNil() {
		return
	}
	sub := newMemWriter()
	v.EncodeFields(sub)
	w.data.set(k, sub.data)
}
func (w *memFieldWriter) Array(k string, n int) model.ArrayWriter {
	return &memArrayWriter{w: w, key: k, items: make([]any, 0, n)}
}

type memArrayWriter struct {
	w     *memFieldWriter
	key   string
	items []any
}

func (a *memArrayWriter) Object(v model.Encodable) {
	if v == nil || v.IsNil() {
		return
	}
	sub := newMemWriter()
	v.EncodeFields(sub)
	a.items = append(a.items, sub.data)
}
func (a *memArrayWriter) String(v string) { a.items = append(a.items, v) }
func (a *memArrayWriter) Int(v int64)      { a.items = append(a.items, v) }
func (a *memArrayWriter) Float(v float64) { a.items = append(a.items, v) }
func (a *memArrayWriter) Bool(v bool)      { a.items = append(a.items, v) }
func (a *memArrayWriter) Bytes(v []byte)   { a.items = append(a.items, string(v)) }
func (a *memArrayWriter) Close()            { a.w.data.set(a.key, a.items) }

type memFieldReader struct {
	data kvMap
}

func (r *memFieldReader) String(k string) (string, bool) {
	if v, ok := r.data.get(k); ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}
func (r *memFieldReader) Int(k string) (int64, bool) {
	if v, ok := r.data.get(k); ok {
		if i, ok := v.(int64); ok {
			return i, true
		}
	}
	return 0, false
}
func (r *memFieldReader) Float(k string) (float64, bool) {
	if v, ok := r.data.get(k); ok {
		if f, ok := v.(float64); ok {
			return f, true
		}
	}
	return 0, false
}
func (r *memFieldReader) Bool(k string) (bool, bool) {
	if v, ok := r.data.get(k); ok {
		if b, ok := v.(bool); ok {
			return b, true
		}
	}
	return false, false
}
func (r *memFieldReader) Bytes(k string) ([]byte, bool) {
	if v, ok := r.data.get(k); ok {
		if b, ok := v.([]byte); ok {
			return b, true
		}
		if s, ok := v.(string); ok {
			return []byte(s), true
		}
	}
	return nil, false
}
func (r *memFieldReader) Raw(k string) (string, bool) {
	if v, ok := r.data.get(k); ok {
		if s, ok := v.(string); ok {
			return s, true
		}
	}
	return "", false
}
func (r *memFieldReader) Object(k string, dest model.Decodable) bool {
	if v, ok := r.data.get(k); ok {
		if subData, ok := v.(kvMap); ok {
			subReader := &memFieldReader{data: subData}
			dest.DecodeFields(subReader)
			return true
		}
	}
	return false
}
func (r *memFieldReader) Array(k string) (model.ArrayReader, bool) {
	if v, ok := r.data.get(k); ok {
		if items, ok := v.([]any); ok {
			return &memArrayReader{items: items}, true
		}
	}
	return nil, false
}

type memArrayReader struct {
	items []any
}

func (a *memArrayReader) Len() int { return len(a.items) }
func (a *memArrayReader) Object(i int, dest model.Decodable) bool {
	if i >= 0 && i < len(a.items) {
		if subData, ok := a.items[i].(kvMap); ok {
			subReader := &memFieldReader{data: subData}
			dest.DecodeFields(subReader)
			return true
		}
	}
	return false
}
func (a *memArrayReader) String(i int) string {
	if i >= 0 && i < len(a.items) {
		if v, ok := a.items[i].(string); ok {
			return v
		}
	}
	return ""
}
func (a *memArrayReader) Int(i int) int64 {
	if i >= 0 && i < len(a.items) {
		if v, ok := a.items[i].(int64); ok {
			return v
		}
	}
	return 0
}
func (a *memArrayReader) Float(i int) float64 {
	if i >= 0 && i < len(a.items) {
		if v, ok := a.items[i].(float64); ok {
			return v
		}
	}
	return 0
}
func (a *memArrayReader) Bool(i int) bool {
	if i >= 0 && i < len(a.items) {
		if v, ok := a.items[i].(bool); ok {
			return v
		}
	}
	return false
}
func (a *memArrayReader) Bytes(i int) []byte {
	if i >= 0 && i < len(a.items) {
		if v, ok := a.items[i].([]byte); ok {
			return v
		}
		if v, ok := a.items[i].(string); ok {
			return []byte(v)
		}
	}
	return nil
}

// MemoryCodec allows encoding and decoding a Content document safely in-memory.
type MemoryCodec struct{}

func (mc MemoryCodec) Encode(c *Content) string {
	w := newMemWriter()
	c.EncodeFields(w)
	return renderMap(w.data)
}

func (mc MemoryCodec) Decode(raw string, dest *Content) {
	data := parseMap(raw)
	r := &memFieldReader{data: data}
	dest.DecodeFields(r)
}

func hexChar(b byte) byte {
	if b < 10 {
		return '0' + b
	}
	return 'A' + (b - 10)
}

func hexVal(b byte) int {
	if b >= '0' && b <= '9' {
		return int(b - '0')
	}
	if b >= 'A' && b <= 'F' {
		return int(b - 'A' + 10)
	}
	if b >= 'a' && b <= 'f' {
		return int(b - 'a' + 10)
	}
	return -1
}

func encodeStr(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		b := s[i]
		if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_' || b == '-' || b == '.' {
			out = append(out, b)
		} else {
			out = append(out, '%', hexChar(b>>4), hexChar(b&0x0F))
		}
	}
	return string(out)
}

func decodeStr(s string) string {
	var out []byte
	for i := 0; i < len(s); i++ {
		if s[i] == '%' && i+2 < len(s) {
			h1 := hexVal(s[i+1])
			h2 := hexVal(s[i+2])
			if h1 != -1 && h2 != -1 {
				out = append(out, byte((h1<<4)|h2))
				i += 2
				continue
			}
		}
		out = append(out, s[i])
	}
	return string(out)
}

func renderMap(m kvMap) string {
	var out string
	for i := 0; i < len(m); i++ {
		if out != "" {
			out += ";"
		}
		out += m[i].key + "=" + renderVal(m[i].val)
	}
	return out
}

func renderVal(v any) string {
	switch x := v.(type) {
	case string:
		return "S" + encodeStr(x)
	case []byte:
		return "S" + encodeStr(string(x))
	case kvMap:
		return "{" + renderMap(x) + "}"
	case []any:
		var arr string
		for i, item := range x {
			if i > 0 {
				arr += ","
			}
			arr += renderVal(item)
		}
		return "[" + arr + "]"
	default:
		return ""
	}
}

func parseMap(s string) kvMap {
	var res kvMap
	if s == "" {
		return res
	}
	pairs := splitTop(s, ';')
	for _, p := range pairs {
		eq := findChar(p, '=')
		if eq == -1 {
			continue
		}
		k := p[:eq]
		vStr := p[eq+1:]
		res.set(k, parseVal(vStr))
	}
	return res
}

func parseVal(s string) any {
	if len(s) == 0 {
		return ""
	}
	if s[0] == 'S' {
		return decodeStr(s[1:])
	}
	if len(s) >= 2 && s[0] == '{' && s[len(s)-1] == '}' {
		return parseMap(s[1 : len(s)-1])
	}
	if len(s) >= 2 && s[0] == '[' && s[len(s)-1] == ']' {
		itemsStr := s[1 : len(s)-1]
		if itemsStr == "" {
			return []any{}
		}
		rawItems := splitTop(itemsStr, ',')
		items := make([]any, len(rawItems))
		for i, it := range rawItems {
			items[i] = parseVal(it)
		}
		return items
	}
	return decodeStr(s)
}

func splitTop(s string, sep byte) []string {
	var parts []string
	var start int
	depthBrace := 0
	depthBracket := 0
	for i := 0; i < len(s); i++ {
		b := s[i]
		if b == '{' {
			depthBrace++
		} else if b == '}' {
			depthBrace--
		} else if b == '[' {
			depthBracket++
		} else if b == ']' {
			depthBracket--
		} else if b == sep && depthBrace == 0 && depthBracket == 0 {
			parts = append(parts, s[start:i])
			start = i + 1
		}
	}
	if start <= len(s) {
		parts = append(parts, s[start:])
	}
	return parts
}

func findChar(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
