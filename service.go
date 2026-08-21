package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var ServiceModel = model.Definition{
	Name: "service",
	Fields: model.Fields{
		{Name: "Slug", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-")}},
		{Name: "Title", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Description", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Image", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "Body", Type: input.Textarea()},
	},
}
