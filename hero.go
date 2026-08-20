package sitecontent

import "github.com/tinywasm/model"

var ImageItemModel = model.Definition{
	Name: "image_item",
	Fields: model.Fields{
		{Name: "Key", Type: model.Text()},
	},
}

var LinkModel = model.Definition{
	Name: "link",
	Fields: model.Fields{
		{Name: "Text", Type: model.Text()},
		{Name: "URL", Type: model.Text()},
	},
}

var HeroModel = model.Definition{
	Name: "hero",
	Fields: model.Fields{
		{Name: "Title", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Subtitle", Type: model.Text()},
		{Name: "CTAs", Type: model.StructSlice(&LinkModel)},
		{Name: "Images", Type: model.StructSlice(&ImageItemModel)},
	},
}
