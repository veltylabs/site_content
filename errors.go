package sitecontent

import (
	"github.com/tinywasm/fmt"
)

const (
	ErrRequired      = "required"
	ErrInvalidSlug   = "site_content: slug invalido %q: solo minusculas, numeros y guiones"
	ErrDuplicateAttr = "site_content: dos servicios comparten %s: %q"
	ErrInvalidColor  = "site_content: color primario invalido: debe ser #rrggbb"
)

var (
	ErrNotFound = fmt.Err("site_content", "not found")
)
