package sitecontent

import "github.com/tinywasm/model"

var StatModel = model.Definition{
	Name: "stat",
	Fields: model.Fields{
		{Name: "Value", Type: model.Text()},
		{Name: "Label", Type: model.Text()},
	},
}
