package sitecontent

import (
	"webtyp.com/input"
	"webtyp.com/model"
)

var StatModel = model.Definition{
	Name: "stat",
	Fields: model.Fields{
		{Name: "Value", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Extra: []rune("+%.,")}},
		{Name: "Label", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(".,-")}},
	},
}
