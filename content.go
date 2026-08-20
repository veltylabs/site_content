package sitecontent

import "github.com/tinywasm/model"

var ContentModel = model.Definition{
	Name: "content",
	Fields: model.Fields{
		{Name: "SiteID", Type: model.Text(), NotNull: true, DB: &model.FieldDB{PK: true}},
		{Name: "Brand", Type: model.Struct(&BrandModel)},
		{Name: "Contact", Type: model.Struct(&ContactModel)},
		{Name: "Hero", Type: model.Struct(&HeroModel)},
		{Name: "About", Type: model.Struct(&AboutModel)},
		{Name: "Services", Type: model.StructSlice(&ServiceModel)},
		{Name: "Stats", Type: model.StructSlice(&StatModel)},
		{Name: "Hours", Type: model.StructSlice(&ScheduleModel)},
		{Name: "Map", Type: model.Struct(&MapModel)},
		{Name: "SEO", Type: model.Struct(&SEOModel)},
		{Name: "Images", Type: model.StructSlice(&ImageRefModel)},
	},
}
