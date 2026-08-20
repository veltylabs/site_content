package sitecontent

import (
	"github.com/tinywasm/fmt"
	"github.com/tinywasm/model"
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

// Validate runs model.ValidateFields over the entire document and additionally
// enforces rules that the schema alone cannot express: the three uniqueness rules
// of Service, slug formatting, and PrimaryColor format.
func Validate(c *Content) error {
	if c == nil {
		return fmt.Err("site_content", "nil content")
	}

	// First validate root Content fields and nested single-struct fields via model.ValidateFields
	if err := model.ValidateFields(model.ActionCreate, c); err != nil {
		return err
	}

	// Validate Brand PrimaryColor format if present
	if c.Brand.PrimaryColor != "" {
		if !isValidHexColor(c.Brand.PrimaryColor) {
			return fmt.Err("site_content: color primario invalido: debe ser #rrggbb")
		}
	}

	// Validate slice elements
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

	// Service uniqueness checks (using pair-wise slice comparisons, no map[K]V)
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
