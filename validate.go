package sitecontent

import (
	"webtyp.com/fmt"
	"webtyp.com/model"
)

func isValidHexColor(c string) bool {
	if len(c) != 7 || c[0] != '#' {
		return false
	}
	for i := 1; i < 7; i++ {
		b := c[i]
		if !(b >= '0' && b <= '9') && !(b >= 'a' && b <= 'f') && !(b >= 'A' && b <= 'F') {
			return false
		}
	}
	return true
}

func isValidSlug(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i := 0; i < len(s); i++ {
		b := s[i]
		if !(b >= 'a' && b <= 'z') && !(b >= '0' && b <= '9') && b != '-' {
			return false
		}
	}
	return true
}

// Validate ejecuta model.ValidateFields sobre todo el documento y adicionalmente
// impone reglas que el esquema por sí solo no puede expresar: las tres reglas de
// unicidad de Service, el formato de slug, y el formato de PrimaryColor.
//
// Ejemplo: err := sitecontent.Validate(content)
func Validate(c *Content) error {
	if c == nil {
		return fmt.Err("site_content", "nil content")
	}

	// Primero valida los campos raíz de Content y los campos de estructura única anidados vía model.ValidateFields
	if err := model.ValidateFields(model.ActionCreate, c); err != nil {
		return err
	}

	// Valida que el Título del Hero no sea nulo
	if err := model.ValidateFields(model.ActionCreate, &c.Hero); err != nil {
		return err
	}

	// Valida el formato de PrimaryColor de Brand si está presente
	if c.Brand.PrimaryColor != "" {
		if !isValidHexColor(c.Brand.PrimaryColor) {
			return fmt.Err("site_content: color primario invalido: debe ser #rrggbb")
		}
	}

	// Valida los elementos del slice (arreglo)
	for i := 0; i < len(c.Services); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Services[i]); err != nil {
			return err
		}
		if !isValidSlug(c.Services[i].Slug) {
			return fmt.Errf("site_content: slug invalido %q: solo minusculas, numeros y guiones", c.Services[i].Slug)
		}
	}

	for i := 0; i < len(c.Images); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Images[i]); err != nil {
			return err
		}
	}

	for i := 0; i < len(c.Hero.CtAs); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Hero.CtAs[i]); err != nil {
			return err
		}
	}

	for i := 0; i < len(c.Hero.Images); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Hero.Images[i]); err != nil {
			return err
		}
	}

	for i := 0; i < len(c.Stats); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Stats[i]); err != nil {
			return err
		}
	}

	for i := 0; i < len(c.Hours); i++ {
		if err := model.ValidateFields(model.ActionCreate, &c.Hours[i]); err != nil {
			return err
		}
	}

	// Comprobaciones de unicidad de Service (usando comparaciones de pares de slices, sin map[K]V)
	n := len(c.Services)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if c.Services[i].Slug == c.Services[j].Slug {
				return fmt.Errf("site_content: dos servicios comparten slug: %q", c.Services[i].Slug)
			}
			if c.Services[i].Title == c.Services[j].Title {
				return fmt.Errf("site_content: dos servicios comparten titulo: %q", c.Services[i].Title)
			}
			if c.Services[i].Description == c.Services[j].Description {
				return fmt.Errf("site_content: dos servicios comparten descripcion: %q", c.Services[i].Description)
			}
		}
	}

	return nil
}
