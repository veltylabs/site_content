package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

// Body, Mission y Vision son prosa: llevan textarea. El resto del charset se
// declara aqui, no en el panel, para que un renombre de campo no degrade el
// widget en silencio.
var AboutModel = model.Definition{
	Name: "about",
	Fields: model.Fields{
		{Name: "Title", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Body", Type: input.Textarea()},
		{Name: "Image", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "Mission", Type: input.Textarea()},
		{Name: "Vision", Type: input.Textarea()},
	},
}
