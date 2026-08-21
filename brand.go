package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

// Los widgets salen del esquema, no del panel que lo consume: el formulario y
// la validacion quedan sincronizados por construccion. Ver AGENTS.md, seccion
// de la whitelist: form/input es un Kind, no un renderizador.
var BrandModel = model.Definition{
	Name: "brand",
	Fields: model.Fields{
		{Name: "Name", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "WideLogo", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "CompactLogo", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "LogoAlt", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,-")}},
		{Name: "PrimaryColor", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("#")}},
	},
}
