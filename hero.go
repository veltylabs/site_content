package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var ImageItemModel = model.Definition{
	Name: "image_item",
	Fields: model.Fields{
		{Name: "Key", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
	},
}

var LinkModel = model.Definition{
	Name: "link",
	Fields: model.Fields{
		{Name: "Text", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,-")}},
		{Name: "URL", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.:#?=&+%")}},
	},
}

var HeroModel = model.Definition{
	Name: "hero",
	Fields: model.Fields{
		{Name: "Title", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Subtitle", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "CTAs", Type: model.StructSlice(&LinkModel)},
		{Name: "Images", Type: model.StructSlice(&ImageItemModel)},
	},
}
