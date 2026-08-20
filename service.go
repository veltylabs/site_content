package sitecontent

import "github.com/tinywasm/model"

var ServiceModel = model.Definition{
	Name: "service",
	Fields: model.Fields{
		{Name: "Slug", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-")}},
		{Name: "Title", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Description", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "Image", Type: model.Text()},
		{Name: "Body", Type: model.Text()},
	},
}
