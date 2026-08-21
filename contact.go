package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var ContactModel = model.Definition{
	Name: "contact",
	Fields: model.Fields{
		{Name: "Phone", Type: input.Phone(), NotNull: true, Permitted: model.Permitted{Numbers: true, Extra: []rune("+-() ")}},
		{Name: "Email", Type: input.Email(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("@._-")}},
		{Name: "Address", Type: input.Address(), NotNull: true, Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,#-")}},
	},
}
