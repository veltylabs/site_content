package sitecontent

import "github.com/tinywasm/model"

var BrandModel = model.Definition{
	Name: "brand",
	Fields: model.Fields{
		{Name: "Name", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("=;{}[].,#_:-")}},
		{Name: "WideLogo", Type: model.Text()},
		{Name: "CompactLogo", Type: model.Text()},
		{Name: "LogoAlt", Type: model.Text()},
		{Name: "PrimaryColor", Type: model.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("#")}},
	},
}
