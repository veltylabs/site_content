---
PLAN: "feat: esquema del contenido editable de un sitio de cliente"
TAG: v0.1.0
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 13951118495392759686
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.

# Plan — `veltylabs/site_content` v0.1.0

Crear **el esquema del contenido editable** de un sitio de cliente de Velty,
dentro del producto `misitio` (`misitio.velty.cl`).

## Por qué existe este módulo, en una frase

Para que exista **un** esquema y no dos: el panel valida con él antes de
guardar, y el compilador del sitio valida con él antes de construir. El contrato
entre ambos lados es **la versión que sus `go.mod` declaran**, no un documento
que alguien tenga que mantener sincronizado.

Si alguna vez alguien copia estos structs a otro repositorio "para no depender",
este módulo perdió su motivo de existir.

---

## 0. Reglas de desarrollo — léelas completas antes de escribir código

Están en [`AGENTS.md`](../AGENTS.md), plantilla canónica de todo
`veltylabs/modules/*`. Lo que sigue repite lo que **no** puedes dejar de
aplicar, porque un agente sin contexto no va a ir a buscarlo.

### 0.1 Lista blanca de importaciones

Los archivos **no-test** pueden importar, de `github.com/tinywasm/*`, sólo:

`model` · `router` · `view` · `events` · `orm` · `storage` · `ddl` ·
`form/input` · `fmt` · `time`

### 0.2 Lista negra — **ni siquiera en `_test.go`**

- **`tinywasm/layout` en cualquier forma.** Es la prohibición central de este
  módulo: describe **datos**, nunca cómo se ven. Quien mapea este contenido a
  una plantilla es `veltylabs/sitetheme`, que vive fuera de
  `veltylabs/modules/*` justamente porque sí nombra un renderizador.
- Backends concretos: `tinywasm/sqlite`, `tinywasm/postgres`, `tinywasm/sqlt`,
  `tinywasm/indexdb`, cualquier driver `database/sql`.
  **Los tests usan `orm.New(mem.New())`** — `github.com/tinywasm/storage/mem`.
- Transportes concretos: `tinywasm/mcp`, `tinywasm/server`, `httpd`, `net/http`.
- `tinywasm/unixid`: recibe `model.IDGenerator` por inyección.
- Encoders concretos: `tinywasm/json`, `tinywasm/jsvalue`.
- Carpetas `internal/`.

### 0.3 Este módulo compila con **TinyGo dentro de un Worker**

- **Sin `map[K]V`.** Slices + búsqueda lineal, o structs de campos fijos.
- **Sin `fmt`, `errors`, `strconv`, `strings`, `log`** de la stdlib. Usa
  `github.com/tinywasm/fmt`.
- **Sin `encoding/json`, sin `reflect`.**

El binario que hospeda este módulo tiene un límite **duro** de 1 MB impuesto por
Cloudflare. No es un presupuesto que se pueda subir: el despliegue falla.

### 0.4 Sin strings mágicos, idioma y estructura

Todo string repetido es una constante nombrada. Código en inglés; documentación
y comentarios de prosa en español. Jerarquía plana, archivos de menos de 500
líneas, todos los tests bajo `tests/`.

---

## 1. Qué construir

Un documento por sitio. `Content` es la raíz; el resto son sus partes.

### 1.1 `Content` — `content.go`

| Campo | Tipo | Regla |
|---|---|---|
| `SiteID` | string | clave, `NotNull` |
| `Brand` | `Brand` | |
| `Contact` | `Contact` | |
| `Hero` | `Hero` | |
| `About` | `About` | |
| `Services` | `[]Service` | |
| `Stats` | `[]Stat` | |
| `Hours` | `[]Schedule` | |
| `Map` | `Map` | |
| `SEO` | `SEO` | |
| `Images` | `[]ImageRef` | |

### 1.2 Las partes

Un archivo por dominio, no todo en `content.go`:

**`brand.go` — `Brand`**

| Campo | Regla |
|---|---|
| `Name` | `NotNull` |
| `WideLogo`, `CompactLogo` | clave de imagen |
| `LogoAlt` | |
| `PrimaryColor` | `#rrggbb`; **lo único visual que el cliente controla** |

`PrimaryColor` es la única excepción a "el cliente no toca el estilo": es un
token CSS que `layout/landing` ya consume, así que no hay nada que construir y
es lo primero que pide todo cliente.

**`contact.go` — `Contact`**: `Phone`, `Email`, `Address`. Los tres `NotNull`.

**`hero.go` — `Hero`**: `Title` (`NotNull`), `Subtitle`, `CTAs []Link`,
`Images []string`.

**`about.go` — `About`**: `Title`, `Body`, `Image`, `Mission`, `Vision`.

**`service.go` — `Service`** — el tipo con más reglas:

| Campo | Regla |
|---|---|
| `Slug` | `NotNull`, sólo `[a-z0-9-]`, **único en el sitio** |
| `Title` | `NotNull`, **único en el sitio** |
| `Description` | `NotNull`, **único en el sitio** |
| `Image` | clave de imagen |
| `Body` | texto de la subpágina |

**Las tres unicidades no son un lujo.** Cada `Service` genera una tarjeta **y
una subpágina**; `layout/landing` entra en pánico ante títulos o descripciones
duplicados y ese pánico rompe el build de publicación. Validarlo aquí convierte
un build roto en un mensaje de formulario.

**`stat.go` — `Stat`**: `Value`, `Label`.

**`schedule.go` — `Schedule`**: `Days`, `Hours`.

**`sitemap.go` — `Map`**: `EmbedURL`.

**`seo.go` — `SEO`**: `Description` (`NotNull`), `SocialImage`,
`SchemaType` (el tipo schema.org del rubro: `LocalBusiness`, `MedicalClinic`, …).

**`image.go` — `ImageRef`**

| Campo | Regla |
|---|---|
| `Key` | **clave de R2**, `NotNull` |
| `Alt` | `NotNull` — es accesibilidad y es SEO |
| `Usage` | dónde se usa |

**`Key` es una clave de R2, jamás una URL absoluta.** La URL final la construye
el tema en el dominio del sitio; guardar una URL aquí ataría el contenido al
hosting y rompería el día que cambie.

### 1.3 Validación — `validate.go`

```go
// Validate corre model.ValidateFields sobre todo el documento y ademas las
// reglas que el esquema por si solo no expresa: las tres unicidades de
// Service y el formato de los slugs.
func Validate(c *Content) error
```

Mensajes textuales, que el panel muestra tal cual:

```
site_content: dos servicios comparten %s: %q
site_content: slug invalido %q: solo minusculas, numeros y guiones
```

**Anti-footgun (TinyGo):** la unicidad se verifica con **slices y comparación
por pares**, no con un `map[string]bool`. Son unidades de servicios, no miles, y
un mapa aquí inflaría el binario del Worker hacia un límite que no se puede
subir.

---

## 2. Estructura de archivos

```
content.go     brand.go     contact.go    hero.go
about.go       service.go   stat.go       schedule.go
sitemap.go     seo.go       image.go
validate.go
module.go      // Module, Deps, New(), Ops, CreateTable
errors.go      // errores centinela
*_orm.go       // GENERADOS por ormc — no editar a mano
docs/ARCHITECTURE.md
docs/diagrams/database.md
tests/
```

Borra `site_content.go` (el archivo vacío del andamiaje). Verificación:
`ls site_content.go` → no existe.

```go
type Deps struct {
	DB  *orm.DB
	IDs model.IDGenerator // INYECTADO — nunca construido aqui
}

func New(d Deps) (*Module, error)
```

`New` devuelve error si falta una dependencia. Chequeo de contrato junto a la
implementación: `var _ router.OpModule = (*Module)(nil)`.

---

## 3. Tests — `tests/`

Con `orm.New(mem.New())`.

| # | Caso | Espera |
|---|---|---|
| 1 | `Content` con `Contact.Phone` vacío | error de `NotNull` |
| 2 | dos `Service` con el mismo `Slug` | error con el mensaje textual de §1.3 |
| 3 | dos `Service` con el mismo `Title` | error |
| 4 | dos `Service` con la misma `Description` | error |
| 5 | `Slug` con mayúsculas o espacios | error de formato |
| 6 | `ImageRef.Alt` vacío | error de `NotNull` |
| 7 | `PrimaryColor` fuera de `#rrggbb` | error |
| 8 | contenido válido completo | `Validate` devuelve `nil` |
| 9 | guardar y releer un `Content` | ida y vuelta idéntico |
| 10 | `Content` de otro `SiteID` | no aparece al consultar por sitio |

Los casos 2, 3 y 4 son los que evitan que un build de publicación explote en
CI. Escríbelos aunque los tres parezcan el mismo test.

---

## 4. Documentación

- `docs/ARCHITECTURE.md` — el documento y sus partes, tabla de Ops, ejemplo de
  raíz de composición. **Sin código de implementación.** Debe decir
  explícitamente que la composición es fija y que este módulo no conoce
  plantillas.
- `docs/diagrams/database.md` — ERD en Mermaid. **Nunca uses `subgraph`**:
  rompe el renderizado en el TUI. Usa `flowchart TD` y `<br/>` para los saltos.
- `README.md` — inicio rápido, tabla de Ops, archivos clave. Enlaza todo lo de
  `docs/` **excepto** `PLAN.md`, que es efímero.

---

## 5. Criterios de aceptación

- [ ] `go vet ./...` limpio; `go test ./tests/...` en verde con los 10 casos.
- [ ] `grep -rn "tinywasm/layout" .` → **vacío**. Es la regla que define al módulo.
- [ ] `grep -rn "tinywasm/sqlite\|tinywasm/postgres\|tinywasm/sqlt\|net/http\|tinywasm/mcp\|tinywasm/unixid\|tinywasm/json\|tinywasm/jsvalue" .` → **vacío, tests incluidos**.
- [ ] `grep -rn "map\[" --include=*.go . | grep -v _test.go` → vacío.
- [ ] `grep -rn "encoding/json\|\"reflect\"\|\"strings\"\|\"errors\"\|\"strconv\"\|\"log\"" --include=*.go .` → vacío.
- [ ] `grep -rn "http://\|https://" --include=*.go . | grep -v _test.go` → vacío: **`ImageRef.Key` no es una URL**.
- [ ] `grep -rn "internal/" .` → vacío.
- [ ] `ls site_content.go` → no existe.
- [ ] Códecs generados con `ormc`, commiteados.
- [ ] `grep -rn "subgraph" docs/` → vacío.

## 6. Fuera de alcance

Sitios, membresías y planes (viven en `veltylabs/site_manager`); el mapeo a la
plantilla (`veltylabs/sitetheme`); cualquier ruta HTTP; y **el orden o la
presencia de las secciones**: la composición es fija y no es un dato de este
módulo.
