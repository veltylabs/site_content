package sitecontent

import "github.com/tinywasm/model"

var ContactModel = model.Definition{
	Name: "contact",
	Fields: model.Fields{
		{Name: "Phone", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Numbers: true, Extra: []rune("+-() ")}},
		{Name: "Email", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("@._-")}},
		{Name: "Address", Type: model.Text(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,#-")}},
	},
}
