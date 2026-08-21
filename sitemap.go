package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var MapModel = model.Definition{
	Name: "map",
	Fields: model.Fields{
		{Name: "EmbedURL", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.:#?=&+%,")}},
	},
}
