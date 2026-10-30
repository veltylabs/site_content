# Arquitectura — `veltylabs/site_content`

## Resumen

`veltylabs/site_content` provee el esquema, persistencia, validación y operaciones para el contenido editable de los sitios de clientes bajo el producto `misitio` (`misitio.velty.cl`).

Este módulo existe para que haya **exactamente un** contrato de esquema entre el panel de edición (que valida antes de guardar) y el compilador del sitio (que valida antes de construir).

## Principios de Dominio

- **Composición Fija**: Los clientes editan contenido dentro de modelos de entidad definidos (`Brand`, `Contact`, `Hero`, `About`, `Service`, `Stat`, `Schedule`, `Map`, `SEO`, `ImageRef`). No eligen ni reordenan secciones de layout.
- **Agnóstico del Renderizador**: Este módulo no importa ni referencia `tinywasm/layout` ni ningún kit de UI. El mapeo del tema es manejado externamente por `veltylabs/sitetheme`.
- **Independiente del Hosting**: `ImageRef.Key` almacena claves de objeto R2, nunca URLs absolutas, evitando la dependencia de un solo proveedor de hosting.
- **Unicidad Fail-Closed**: `Service.Slug`, `Service.Title`, y `Service.Description` deben ser únicos por sitio para prevenir fallos (panics) en tiempo de publicación.
- **Compatible con TinyGo / Worker**: Diseñado para ejecutarse dentro de Cloudflare Workers sin usar `map[K]V` de Go, `reflect`, `encoding/json`, ni paquetes pesados de la biblioteca estándar.

## Entidades

- `Content`: Documento raíz que contiene `SiteID` (Clave Primaria) y modelos de dominio incrustados.
- `Brand`: Nombre del sitio, color primario, logos y texto alternativo.
- `Contact`: Información de contacto (`Phone`, `Email`, `Address`).
- `Hero`: Título de la sección hero, subtítulo, CTAs e imágenes.
- `About`: Detalles de la sección acerca de (`Title`, `Body`, `Image`, `Mission`, `Vision`).
- `Service`: Subpáginas de servicio (`Slug`, `Title`, `Description`, `Image`, `Body`).
- `Stat`: Métricas y estadísticas (`Value`, `Label`).
- `Schedule`: Horas de operación (`Days`, `Hours`).
- `Map`: Configuración del mapa incrustado (`EmbedURL`).
- `SEO`: Configuraciones de optimización para motores de búsqueda (`Description`, `SocialImage`, `SchemaType`).
- `ImageRef`: Referencia a un activo (`Key`, `Alt`, `Usage`).

## Superficies de Acceso y Modelo de Seguridad

Este módulo expone dos superficies de acceso distintas con diferentes contratos de seguridad:

1. **Métodos de Servicio (`Get`, `Save`)**:
   - **No** realizan comprobaciones de membresía ni autorización de membresía.
   - Diseñados para llamadas internas por una raíz de composición (ej. manejadores REST) que maneja la autorización independientemente.
   - `Deps.Members` puede ser `nil` sin afectar la invocación directa de los métodos de servicio.

2. **Operaciones del Módulo (`OpGet`, `OpSave`)**:
   - Imponen acceso por membresía a través de la interfaz inyectada `MemberChecker`.
   - Ejecutan verificación de membresía antes de leer o validar/escribir datos de dominio.
   - Deniegan el acceso a la operación (`403 Forbidden`) si `Deps.Members` es `nil` o si `CanEditContent(userID, siteID)` retorna `false`.

## Operaciones

| Op | Recurso | Acción | Descripción |
|---|---|---|---|
| `get` | `site_content` | `model.Read` | Obtiene el contenido del sitio por `SiteID` (requiere verificación `MemberChecker`). |
| `save` | `site_content` | `model.Create \| model.Update` | Valida y guarda el contenido del sitio (requiere verificación `MemberChecker`). |

## Ejemplo de Raíz de Composición

```go
package main

import (
	"github.com/webtyp/orm"
	"github.com/webtyp/storage/mem"
	"github.com/veltylabs/site_content"
)

type DummyIDs struct{}
func (d DummyIDs) NewID() string { return "id-1" }

func main() {
	db := orm.New(mem.New())
	m, err := sitecontent.New(sitecontent.Deps{
		DB:  db,
		IDs: DummyIDs{},
	})
	if err != nil {
		panic(err)
	}

	content := &sitecontent.Content{
		SiteId: "site-123",
		Brand: sitecontent.Brand{
			Name: "My Business",
			PrimaryColor: "#112233",
		},
		Contact: sitecontent.Contact{
			Phone: "+56912345678",
			Email: "contacto@example.com",
			Address: "Av. Providencia 123",
		},
		Seo: sitecontent.SEO{
			Description: "Sitio web oficial",
		},
	}

	if err := m.Save(content); err != nil {
		panic(err)
	}
}
```
