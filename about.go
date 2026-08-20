package sitecontent

import "github.com/tinywasm/model"

var AboutModel = model.Definition{
	Name: "about",
	Fields: model.Fields{
		{Name: "Title", Type: model.Text()},
		{Name: "Body", Type: model.Text()},
		{Name: "Image", Type: model.Text()},
		{Name: "Mission", Type: model.Text()},
		{Name: "Vision", Type: model.Text()},
	},
}
