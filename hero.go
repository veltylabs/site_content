package sitecontent

import "github.com/tinywasm/model"

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
		{Name: "Title", Type: model.Text()},
		{Name: "Subtitle", Type: model.Text()},
		{Name: "CTAs", Type: model.StructSlice(&LinkModel)},
		{Name: "Images", Type: model.Text()},
	},
}
