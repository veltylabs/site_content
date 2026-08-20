package tests

import (
	"testing"

	"github.com/tinywasm/orm"
	"github.com/tinywasm/storage/mem"
	"github.com/veltylabs/site_content"
)

type mockIDGen struct{}

func (m mockIDGen) NewID() string {
	return "test-id"
}

func setupModule(t *testing.T) *sitecontent.Module {
	db := orm.New(mem.New())
	m, err := sitecontent.New(sitecontent.Deps{
		DB:  db,
		IDs: mockIDGen{},
	})
	if err != nil {
		t.Fatalf("failed to create module: %v", err)
	}
	return m
}

func validContent() *sitecontent.Content {
	return &sitecontent.Content{
		SiteId: "site-1",
		Brand: sitecontent.Brand{
			Name:         "Empresa Test; con punto y coma = igual {llaves} [corchetes]",
			PrimaryColor: "#112233",
		},
		Contact: sitecontent.Contact{
			Phone:   "+56912345678",
			Email:   "test@example.com",
			Address: "Calle Test 123, Depto 4B",
		},
		Seo: sitecontent.SEO{
			Description: "Descripcion SEO con ; = { } y [ ]",
		},
		Services: []sitecontent.Service{
			{
				Slug:        "servicio-uno",
				Title:       "Servicio Uno [Draft]",
				Description: "Descripcion Uno; detallada = yes",
			},
			{
				Slug:        "servicio-dos",
				Title:       "Servicio Dos",
				Description: "Descripcion Dos",
			},
		},
		Images: []sitecontent.ImageRef{
			{
				Key: "r2-key-1",
				Alt: "Imagen principal",
			},
		},
	}
}

// Case 1: Content con Contact.Phone vacío -> error NotNull
func TestValidation_MissingContactPhone(t *testing.T) {
	c := validContent()
	c.Contact.Phone = ""
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for empty Contact.Phone, got nil")
	}
}

// Case 2: dos Service con el mismo Slug -> error textual
func TestValidation_DuplicateServiceSlug(t *testing.T) {
	c := validContent()
	c.Services[1].Slug = c.Services[0].Slug
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for duplicate service slug, got nil")
	}
	expectedMsg := "site_content: dos servicios comparten slug: \"servicio-uno\""
	if err.Error() != expectedMsg {
		t.Fatalf("expected error message %q, got %q", expectedMsg, err.Error())
	}
}

// Case 3: dos Service con el mismo Title -> error
func TestValidation_DuplicateServiceTitle(t *testing.T) {
	c := validContent()
	c.Services[1].Title = c.Services[0].Title
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for duplicate service title, got nil")
	}
}

// Case 4: dos Service con la misma Description -> error
func TestValidation_DuplicateServiceDescription(t *testing.T) {
	c := validContent()
	c.Services[1].Description = c.Services[0].Description
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for duplicate service description, got nil")
	}
}

// Case 5: Slug con mayúsculas o espacios -> error de formato
func TestValidation_InvalidServiceSlug(t *testing.T) {
	c := validContent()
	c.Services[0].Slug = "Servicio-Invalido"
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for invalid service slug format, got nil")
	}
}

// Case 6: ImageRef.Alt vacío -> error NotNull
func TestValidation_MissingImageAlt(t *testing.T) {
	c := validContent()
	c.Images[0].Alt = ""
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for empty ImageRef.Alt, got nil")
	}
}

// Case 7: PrimaryColor fuera de #rrggbb -> error
func TestValidation_InvalidPrimaryColor(t *testing.T) {
	c := validContent()
	c.Brand.PrimaryColor = "blue"
	err := sitecontent.Validate(c)
	if err == nil {
		t.Fatal("expected error for invalid PrimaryColor, got nil")
	}
}

// Case 8: contenido válido completo -> Validate devuelve nil
func TestValidation_ValidContent(t *testing.T) {
	c := validContent()
	err := sitecontent.Validate(c)
	if err != nil {
		t.Fatalf("expected nil for valid content, got: %v", err)
	}
}

// Case 9: guardar y releer un Content -> ida y vuelta idéntico
func TestModule_SaveAndGetRoundtrip(t *testing.T) {
	m := setupModule(t)
	c := validContent()

	if err := m.Save(c); err != nil {
		t.Fatalf("failed to save content: %v", err)
	}

	read, err := m.Get(c.SiteId)
	if err != nil {
		t.Fatalf("failed to get content: %v", err)
	}

	if read.SiteId != c.SiteId {
		t.Errorf("expected SiteId %q, got %q", c.SiteId, read.SiteId)
	}
	if read.Brand.Name != c.Brand.Name {
		t.Errorf("expected Brand.Name %q, got %q", c.Brand.Name, read.Brand.Name)
	}
	if read.Contact.Phone != c.Contact.Phone {
		t.Errorf("expected Contact.Phone %q, got %q", c.Contact.Phone, read.Contact.Phone)
	}
	if read.Seo.Description != c.Seo.Description {
		t.Errorf("expected Seo.Description %q, got %q", c.Seo.Description, read.Seo.Description)
	}
	if len(read.Services) != len(c.Services) {
		t.Fatalf("expected %d services, got %d", len(c.Services), len(read.Services))
	}
	if read.Services[0].Title != c.Services[0].Title {
		t.Errorf("expected Service[0].Title %q, got %q", c.Services[0].Title, read.Services[0].Title)
	}
	if read.Services[0].Description != c.Services[0].Description {
		t.Errorf("expected Service[0].Description %q, got %q", c.Services[0].Description, read.Services[0].Description)
	}
}

// Case 10: Content de otro SiteID -> no aparece al consultar por sitio
func TestModule_SiteIsolation(t *testing.T) {
	m := setupModule(t)
	c1 := validContent()
	c1.SiteId = "site-A"

	if err := m.Save(c1); err != nil {
		t.Fatalf("failed to save site-A content: %v", err)
	}

	_, err := m.Get("site-B")
	if err != sitecontent.ErrNotFound {
		t.Fatalf("expected ErrNotFound for site-B, got: %v", err)
	}
}
