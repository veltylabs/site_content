package sitecontent

const (
	ErrRequired      = "required"
	ErrInvalidSlug   = "site_content: slug invalido %q: solo minusculas, numeros y guiones"
	ErrDuplicateAttr = "site_content: dos servicios comparten %s: %q"
	ErrInvalidColor  = "site_content: color primario invalido: debe ser #rrggbb"
)

// domainError es el tipo concreto de los errores centinela de este paquete. El código
// los compara afirmando este tipo y comparando el valor: == entre dos
// valores de error compila, bajo TinyGo, a runtime.interfaceEqual, lo que incluye
// internal/reflectlite en el binario wasm.
type domainError string

func (e domainError) Error() string { return string(e) }

const (
	ErrNotFound domainError = "site_content not found"
)
