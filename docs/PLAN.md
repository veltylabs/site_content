---
PLAN: "fix: los ops leian y sobrescribian el contenido de cualquier sitio"
EXECUTOR: jules
REVIEWER: none
STATUS: running
SESSION: 5562334572438625732
---

> Este plan se despacha con el flujo CodeJob. Ver skill: agents-workflow.

# Plan — el `site_id` viene del cuerpo y no se comprueba contra nada

## El hallazgo

[`module.go`](../module.go), `OpGet` y `OpSave`:

```go
func (m *Module) OpGet(ctx router.Context) {
	var args Content
	if err := ctx.Decode(&args); err != nil { ctx.WriteStatus(400); return }
	if args.SiteId == "" { ctx.WriteStatus(400); return }
	res, err := m.Get(args.SiteId)
	...
}
```

**El `site_id` llega en el cuerpo de la petición y no se compara con nada.** Cualquier llamante
con `site_content:r` lee el contenido de **cualquier** sitio; cualquiera con
`site_content:u` **sobrescribe** el de cualquier sitio.

Es exactamente la regla que el consumidor de este módulo declara como la más importante de todo
su producto, el caso T2 de su `SPECS.md` §6:

```
El site_id efectivo SIEMPRE sale de la membresía del usuario autenticado.
NUNCA del cuerpo de la petición, NUNCA de la query.
```

Hoy no es explotable: el único consumidor es `veltylabs/misitio`, que monta REST y comprueba la
membresía en sus propios handlers. **Deja de no serlo en cuanto alguien coseche estos ops con
`mcp.HarvestOps`.** Se cierra antes de que eso pase.

## La forma de arreglarlo, y por qué es esta

Este módulo **no puede** importar `site_manager` — el `AGENTS.md` canónico lo prohíbe: *"los
módulos nunca se importan entre sí"*. Y tampoco debe inventarse su propio concepto de
membresía.

Lo que sí manda ese documento es el patrón de **wiring cruzado**: cuando el módulo A necesita
un dato de dominio del módulo B, **A declara la interfaz estrecha que necesita**, en su propio
paquete, y `*B.Module` la satisface estructuralmente sin que nadie importe a nadie. La raíz de
composición conecta las instancias. Es el patrón `CatalogReader`/`StaffReader`.

## Anti-footguns

- Whitelist del [`AGENTS.md` canónico](https://github.com/veltylabs/modules/blob/main/AGENTS.md):
  `model`, `router`, `view`, `events`, `orm`, `storage`, `ddl`, `form/input`, `fmt`, `time`.
  Sin `map[K]V` **ni en tests**, sin `reflect`, sin la stdlib que el ecosistema ya reemplaza.
- **No importes `veltylabs/site_manager`.** Ni siquiera en un test. El punto entero de la etapa
  1 es que no haga falta.
- **No cambies la firma de `Get`, `Save` ni `Validate`.** `veltylabs/misitio` los llama directo
  desde sus handlers REST; este plan cambia el **gate de los ops**, no el dominio.
- **No borres la validación previa a escribir.** `Save` llama a `Validate(c)` antes de tocar la
  base, y ese orden es la garantía de que nunca se marca contenido como bueno sin serlo.
- Repo de `veltylabs`: identificadores en inglés, comentarios y docs en español.

---

## Etapa 1 — El puerto de pertenencia

Archivo nuevo: **`membership.go`**.

```go
// MemberChecker responde si un usuario puede tocar el contenido de un sitio.
//
// Es una interfaz LECTORA ESTRECHA, declarada aquí y satisfecha
// estructuralmente por *sitemanager.Module: los módulos no se importan entre
// sí, y la raíz de composición es quien conecta las instancias (patrón
// CatalogReader del AGENTS.md canónico).
//
// nil DENIEGA. La ausencia de un comprobador no es un permiso — es el mismo
// criterio que model.Allowed aplica a un Authorizer ausente.
type MemberChecker interface {
	// CanEditContent informa si userID puede leer y escribir el contenido de
	// siteID. Una sola pregunta, no un rol: qué significa "puede" es del
	// módulo que conoce las membresías, no de este.
	CanEditContent(userID, siteID string) bool
}
```

**Nombre deliberado.** No es `MemberOf` ni devuelve un rol: este módulo no debe saber que
existen `RoleViewer`/`RoleEditor`/`RoleOwner`. Pide la respuesta, no los datos para calcularla.

`Deps` gana el campo, **opcional en el tipo y obligatorio en el uso**:

```go
type Deps struct {
	DB      *orm.DB
	IDs     model.IDGenerator
	Members MemberChecker // nil DENIEGA en los ops; los métodos de servicio no lo consultan
}
```

`New` **no** falla si `Members` es nil: los métodos de servicio (`Get`, `Save`) siguen siendo
utilizables por una raíz de composición que hace el control por su cuenta —que es justo lo que
`veltylabs/misitio` hace hoy— y romperlos dejaría a su Worker sin arrancar. Lo que sí hace un
`Members` nil es **denegar todos los ops**, con un motivo legible.

## Etapa 2 — Los ops lo consultan

Archivo: **[`module.go`](../module.go)**.

```go
// requireMember es el único sitio donde se decide si un llamante puede tocar
// el contenido de un sitio. Una sola implementación: repetirla en cada op
// garantiza que algún día falte en uno.
func (m *Module) requireMember(ctx router.Context, siteID string) bool {
	userID := ctx.UserID()
	if userID == "" {
		ctx.WriteStatus(401)
		return false
	}
	if siteID == "" {
		ctx.WriteStatus(400)
		return false
	}
	// nil deniega: sin comprobador no hay permiso.
	if m.members == nil || !m.members.CanEditContent(userID, siteID) {
		ctx.WriteStatus(403)
		return false
	}
	return true
}
```

`OpGet` y `OpSave` lo llaman **inmediatamente después de decodificar y antes de tocar la
base**:

```go
func (m *Module) OpGet(ctx router.Context) {
	var args Content
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	if !m.requireMember(ctx, args.SiteId) {
		return
	}
	...
}
```

`403` y no `404`: los identificadores no son adivinables, y un `404` mentiroso hace
indistinguible un fallo de permisos de un dato borrado.

En `OpSave`, el orden es **decodificar → pertenencia → `Save` (que valida y luego escribe)**.
La pertenencia va **antes** de la validación: negarle a un extraño el detalle de qué campos son
inválidos en un sitio ajeno es información que no le corresponde.

## Etapa 3 — Tests

Bajo **`tests/`**, junto a `content_test.go`, con `orm.New(mem.New())` y `router/mock`.

El doble de `MemberChecker` es una función de test local — **no importes `site_manager`**:

```go
type fakeMembers func(userID, siteID string) bool

func (f fakeMembers) CanEditContent(userID, siteID string) bool { return f(userID, siteID) }
```

| Test | Afirma |
|---|---|
| `TestOpGetDeniesNonMember` | `CanEditContent` devuelve false → `403`, y **`Get` nunca se llamó** |
| `TestOpGetAllowsMember` | miembro → `200` con el contenido |
| `TestOpGetRejectsAnonymous` | sin identidad → `401` |
| `TestOpGetRejectsEmptySiteID` | cuerpo sin `SiteId` → `400` |
| `TestOpSaveDeniesForeignSite` | el caso T2: cuerpo con el `SiteId` de otro → `403` y el contenido ajeno **no cambió** |
| `TestOpSaveDeniesBeforeValidating` | cuerpo ajeno **e inválido** → `403`, nunca `400`: el extraño no descubre qué campos existen |
| `TestNilMemberCheckerDenies` | `Deps{Members: nil}` → los dos ops responden `403` |
| `TestServiceMethodsIgnoreMemberChecker` | `Get`/`Save` directos siguen funcionando con `Members: nil` — es como los usa la raíz de composición hoy |

El penúltimo y el último son los que impiden que alguien "simplifique" esto más adelante.

## Etapa 4 — Documentación

- **[`docs/ARCHITECTURE.md`](ARCHITECTURE.md)** — una sección con las **dos superficies** de
  este módulo y su contrato de seguridad, que no es el mismo:
  - los **métodos de servicio** (`Get`, `Save`) no comprueban nada y confían en quien llama —
    son para una raíz de composición que ya hizo el control;
  - los **ops** sí comprueban, contra el `MemberChecker` inyectado, y sin él deniegan.

  Que esa diferencia esté escrita es lo que evita que el próximo consumidor asuma la
  equivocada.
- **`README.md`** — `Deps.Members` en el ejemplo de construcción, con una línea sobre qué
  ocurre si se omite.
- **No enlaces `docs/PLAN.md`** desde ningún documento permanente.

## Criterios de aceptación

- [ ] `go build ./...`, `go vet ./...` limpios; `gotest ./...` en verde.
- [ ] `GOOS=js GOARCH=wasm go vet ./...` limpio.
- [ ] `grep -rn "site_manager\|sitemanager" .` → vacío, **tests incluidos**.
- [ ] `grep -rn "map\[" .` → vacío.
- [ ] `grep -rn "tinywasm/mcp\|tinywasm/json\|tinywasm/unixid\|tinywasm/sqlt\|tinywasm/postgres\|tinywasm/layout"`
      → vacío.
- [ ] `grep -n "requireMember" module.go` → la declaración y **dos** llamadas, una por op.
- [ ] `Get`, `Save` y `Validate` conservan su firma.
- [ ] `New` **no** devuelve error con `Members: nil`.
- [ ] `docs/ARCHITECTURE.md` describe las dos superficies y su contrato distinto.

## Fuera de alcance

`view.Presenter` (sin consumidor todavía), cambios de esquema o de codec, y cualquier cambio a
los métodos de servicio.

## Etapas

| # | Etapa | Archivos |
|---|---|---|
| 1 | El puerto `MemberChecker` | `membership.go`, `module.go` (`Deps`) |
| 2 | Los ops lo consultan | `module.go` |
| 3 | Tests | `tests/` |
| 4 | Documentación | `docs/ARCHITECTURE.md`, `README.md` |
