package sitecontent

import "github.com/tinywasm/model"

var SEOModel = model.Definition{
	Name: "seo",
	Fields: model.Fields{
		{Name: "Description", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "SocialImage", Type: model.Text()},
		{Name: "SchemaType", Type: model.Text()},
	},
}
