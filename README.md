# site_content

Esquema del contenido editable de un sitio de cliente dentro del producto `misitio` (`misitio.velty.cl`).

## Inicio Rápido

```go
import (
	"github.com/tinywasm/orm"
	"github.com/tinywasm/storage/mem"
	"github.com/veltylabs/site_content"
)

db := orm.New(mem.New())
m, err := sitecontent.New(sitecontent.Deps{
	DB: db,
	IDs: myIDGenerator,
})

content := &sitecontent.Content{
	SiteId: "sitio-demo",
	Brand: sitecontent.Brand{Name: "Empresa Demo", PrimaryColor: "#0055ff"},
	Contact: sitecontent.Contact{Phone: "+56900000000", Email: "info@demo.cl", Address: "Santiago"},
	Seo: sitecontent.SEO{Description: "Sitio web de demostracion"},
}

err = m.Save(content)
```

## Operaciones

| Op | Recurso | Acción | Descripción |
|---|---|---|---|
| `get` | `site_content` | `model.Read` | Obtiene el contenido del sitio por `SiteID`. |
| `save` | `site_content` | `model.Create \| model.Update` | Valida y guarda el contenido del sitio. |

## Archivos Clave

- `content.go`: Definición del modelo raíz `Content`.
- `validate.go`: Reglas de validación avanzadas (unicidades de servicio, slugs, formato de color).
- `module.go`: Implementación del módulo y de la interfaz `router.OpModule`.
- `docs/ARCHITECTURE.md`: Documentación de arquitectura y ciclo de vida.
- `docs/diagrams/database.md`: Diagrama de entidades en Mermaid.
