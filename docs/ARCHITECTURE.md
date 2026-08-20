# Architecture — `veltylabs/site_content`

## Overview

`veltylabs/site_content` provides the schema, persistence, validation, and operations for the editable content of customer websites under the `misitio` product (`misitio.velty.cl`).

This module exists so there is **exactly one** schema contract between the editing panel (which validates prior to saving) and the site compiler (which validates prior to building).

## Domain Principles

- **Fixed Composition**: Customers edit content within defined entity models (`Brand`, `Contact`, `Hero`, `About`, `Service`, `Stat`, `Schedule`, `Map`, `SEO`, `ImageRef`). They do not choose or reorder layout sections.
- **Renderer Agnostic**: This module does not import or reference `tinywasm/layout` or any UI kit. Theme mapping is handled externally by `veltylabs/sitetheme`.
- **Hosting Independent**: `ImageRef.Key` stores R2 object keys, never absolute URLs, avoiding hosting lock-in.
- **Fail-Closed Uniqueness**: `Service.Slug`, `Service.Title`, and `Service.Description` must be unique per site to prevent publish-time panics.
- **TinyGo / Worker Compatible**: Designed to execute within Cloudflare Workers without using Go `map[K]V`, `reflect`, `encoding/json`, or heavy stdlib packages.

## Entities

- `Content`: Root document containing `SiteID` (Primary Key) and embedded domain models.
- `Brand`: Site name, primary color, logos, and alt text.
- `Contact`: Contact information (`Phone`, `Email`, `Address`).
- `Hero`: Hero section title, subtitle, CTAs, and images.
- `About`: About section details (`Title`, `Body`, `Image`, `Mission`, `Vision`).
- `Service`: Service subpages (`Slug`, `Title`, `Description`, `Image`, `Body`).
- `Stat`: Metrics and statistics (`Value`, `Label`).
- `Schedule`: Operating hours (`Days`, `Hours`).
- `Map`: Map embed configuration (`EmbedURL`).
- `SEO`: Search engine optimization settings (`Description`, `SocialImage`, `SchemaType`).
- `ImageRef`: Asset reference (`Key`, `Alt`, `Usage`).

## Operations

| Op | Resource | Action | Description |
|---|---|---|---|
| `get` | `site_content` | `model.Read` | Fetches site content by `SiteID`. |
| `save` | `site_content` | `model.Create \| model.Update` | Validates and persists site content. |

## Composition Root Example

```go
package main

import (
	"github.com/tinywasm/orm"
	"github.com/tinywasm/storage/mem"
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
