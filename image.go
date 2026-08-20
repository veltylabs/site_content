package sitecontent

import "github.com/tinywasm/model"

var ImageRefModel = model.Definition{
	Name: "image_ref",
	Fields: model.Fields{
		{Name: "Key", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "Alt", Type: model.Text(), NotNull: true},
		{Name: "Usage", Type: model.Text()},
	},
}
