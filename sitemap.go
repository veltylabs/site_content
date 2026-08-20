package sitecontent

import "github.com/tinywasm/model"

var MapModel = model.Definition{
	Name: "map",
	Fields: model.Fields{
		{Name: "EmbedURL", Type: model.Text()},
	},
}
