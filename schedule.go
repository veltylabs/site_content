package sitecontent

import "github.com/tinywasm/model"

var ScheduleModel = model.Definition{
	Name: "schedule",
	Fields: model.Fields{
		{Name: "Days", Type: model.Text()},
		{Name: "Hours", Type: model.Text()},
	},
}
