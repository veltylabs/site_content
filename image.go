package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var ImageRefModel = model.Definition{
	Name: "image_ref",
	Fields: model.Fields{
		{Name: "Key", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "Alt", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,-")}},
		{Name: "Usage", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_")}},
	},
}
