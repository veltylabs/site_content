package sitecontent

import (
	"webtyp.com/input"
	"webtyp.com/model"
)

var SEOModel = model.Definition{
	Name: "seo",
	Fields: model.Fields{
		{Name: "Description", Type: input.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "SocialImage", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("-_/.")}},
		{Name: "SchemaType", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true}},
	},
}
