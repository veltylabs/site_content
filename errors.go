package sitecontent

const (
	ErrRequired      = "required"
	ErrInvalidSlug   = "site_content: slug invalido %q: solo minusculas, numeros y guiones"
	ErrDuplicateAttr = "site_content: dos servicios comparten %s: %q"
	ErrInvalidColor  = "site_content: color primario invalido: debe ser #rrggbb"
)

// domainError is the concrete type of this package's sentinel errors. Code
// compares them by asserting this type and comparing the value: == between two
// error values compiles, under TinyGo, to runtime.interfaceEqual, which pulls
// internal/reflectlite into the wasm binary.
type domainError string

func (e domainError) Error() string { return string(e) }

const (
	ErrNotFound domainError = "site_content not found"
)
