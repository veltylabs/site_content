package tests

import (
	"testing"

	"github.com/tinywasm/input"
	"github.com/tinywasm/model"
	"github.com/veltylabs/site_content"
)

// Cada campo escalar de cada Definition tiene que declarar un widget de
// form/input, no un model.Kind pelado. Es lo que permite que un panel arme el
// formulario desde el esquema: sin esto form.New falla con "form has no
// renderable field" y el consumidor se ve tentado a adivinar el widget por el
// nombre del campo, que desincroniza formulario y validacion.
func TestEveryScalarFieldDeclaresAWidget(t *testing.T) {
	defs := []*model.Definition{
		&sitecontent.BrandModel,
		&sitecontent.ContactModel,
		&sitecontent.HeroModel,
		&sitecontent.AboutModel,
		&sitecontent.ServiceModel,
		&sitecontent.StatModel,
		&sitecontent.ScheduleModel,
		&sitecontent.ImageRefModel,
		&sitecontent.ImageItemModel,
		&sitecontent.LinkModel,
		&sitecontent.SEOModel,
		&sitecontent.MapModel,
	}

	for _, def := range defs {
		for _, f := range def.Fields {
			if f.Type == nil {
				t.Errorf("%s.%s: sin Kind", def.Name, f.Name)
				continue
			}
			// Las composiciones (Struct/StructSlice) no llevan widget: el
			// formulario desciende al Definition anidado.
			switch f.Type.Storage() {
			case model.FieldStruct, model.FieldStructSlice:
				continue
			}
			if _, ok := f.Type.(input.Input); !ok {
				t.Errorf("%s.%s: Type es %T (%q), no un input.Input — declaralo con input.Text(), input.Email(), …",
					def.Name, f.Name, f.Type, f.Type.Name())
			}
		}
	}
}
