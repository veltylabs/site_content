package sitecontent

import (
	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
)

var ScheduleModel = model.Definition{
	Name: "schedule",
	Fields: model.Fields{
		{Name: "Days", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune("-,")}},
		{Name: "Hours", Type: input.Text(), Permitted: model.Permitted{Letters: true, Numbers: true, Spaces: true, Extra: []rune(":-.")}},
	},
}
