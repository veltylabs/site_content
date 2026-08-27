package tests

import (
	"testing"

	"github.com/tinywasm/orm"
	"github.com/tinywasm/router/mock"
	"github.com/tinywasm/storage/mem"
	"github.com/veltylabs/site_content"
)

type fakeMembers func(userID, siteID string) bool

func (f fakeMembers) CanEditContent(userID, siteID string) bool {
	return f(userID, siteID)
}

func setupModuleWithMembers(t *testing.T, members sitecontent.MemberChecker) *sitecontent.Module {
	db := orm.New(mem.New())
	m, err := sitecontent.New(sitecontent.Deps{
		DB:      db,
		IDs:     mockIDGen{},
		Members: members,
	})
	if err != nil {
		t.Fatalf("failed to create module: %v", err)
	}
	return m
}

func mockCtx(userID string, payload *sitecontent.Content) *mock.Context {
	ctx := &mock.Context{}
	ctx.SetUserID(userID)
	if payload != nil {
		_ = ctx.Encode(payload)
		ctx.InBody = ctx.ResponseBody()
	}
	return ctx
}

// TestOpGetDeniesNonMember: CanEditContent devuelve false -> 403, y Get nunca se llamó
func TestOpGetDeniesNonMember(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return false
	})
	m := setupModuleWithMembers(t, members)

	ctx := mockCtx("user-1", &sitecontent.Content{SiteId: "site-1"})

	m.OpGet(ctx)

	if ctx.Status != 403 {
		t.Fatalf("expected status 403, got %d", ctx.Status)
	}
}

// TestOpGetAllowsMember: miembro -> 200 con el contenido
func TestOpGetAllowsMember(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return userID == "user-1" && siteID == "site-1"
	})
	m := setupModuleWithMembers(t, members)

	c := validContent()
	c.SiteId = "site-1"
	if err := m.Save(c); err != nil {
		t.Fatalf("failed to save content: %v", err)
	}

	ctx := mockCtx("user-1", &sitecontent.Content{SiteId: "site-1"})

	m.OpGet(ctx)

	if ctx.Status != 200 {
		t.Fatalf("expected status 200, got %d", ctx.Status)
	}

	out := &sitecontent.Content{}
	if err := ctx.Decode(out); err != nil {
		t.Fatalf("failed to decode context output: %v", err)
	}
	if out.SiteId != "site-1" {
		t.Fatalf("expected SiteId site-1, got %s", out.SiteId)
	}
}

// TestOpGetRejectsAnonymous: sin identidad -> 401
func TestOpGetRejectsAnonymous(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return true
	})
	m := setupModuleWithMembers(t, members)

	ctx := mockCtx("", &sitecontent.Content{SiteId: "site-1"})

	m.OpGet(ctx)

	if ctx.Status != 401 {
		t.Fatalf("expected status 401, got %d", ctx.Status)
	}
}

// TestOpGetRejectsEmptySiteID: cuerpo sin SiteId -> 400
func TestOpGetRejectsEmptySiteID(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return true
	})
	m := setupModuleWithMembers(t, members)

	ctx := mockCtx("user-1", &sitecontent.Content{SiteId: ""})

	m.OpGet(ctx)

	if ctx.Status != 400 {
		t.Fatalf("expected status 400, got %d", ctx.Status)
	}
}

// TestOpSaveDeniesForeignSite: el caso T2: cuerpo con el SiteId de otro -> 403 y el contenido ajeno no cambió
func TestOpSaveDeniesForeignSite(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return userID == "user-1" && siteID == "my-site"
	})
	m := setupModuleWithMembers(t, members)

	// Guarda contenido inicial para el sitio de otro
	otherContent := validContent()
	otherContent.SiteId = "other-site"
	otherContent.Brand.Name = "Original Brand"
	if err := m.Save(otherContent); err != nil {
		t.Fatalf("failed to setup initial content: %v", err)
	}

	// Intento de sobrescribir el sitio ajeno desde OpSave
	attackerContent := validContent()
	attackerContent.SiteId = "other-site"
	attackerContent.Brand.Name = "Hacked Brand"

	ctx := mockCtx("user-1", attackerContent)

	m.OpSave(ctx)

	if ctx.Status != 403 {
		t.Fatalf("expected status 403, got %d", ctx.Status)
	}

	// El contenido ajeno no cambió
	read, err := m.Get("other-site")
	if err != nil {
		t.Fatalf("failed to get other-site content: %v", err)
	}
	if read.Brand.Name != "Original Brand" {
		t.Fatalf("expected Brand.Name %q, got %q", "Original Brand", read.Brand.Name)
	}
}

// TestOpSaveDeniesBeforeValidating: cuerpo ajeno e inválido -> 403, nunca 400
func TestOpSaveDeniesBeforeValidating(t *testing.T) {
	members := fakeMembers(func(userID, siteID string) bool {
		return false
	})
	m := setupModuleWithMembers(t, members)

	invalidForeignContent := validContent()
	invalidForeignContent.SiteId = "foreign-site"
	invalidForeignContent.Contact.Phone = "" // invalido

	ctx := mockCtx("user-1", invalidForeignContent)

	m.OpSave(ctx)

	if ctx.Status != 403 {
		t.Fatalf("expected status 403, got %d", ctx.Status)
	}
}

// TestNilMemberCheckerDenies: Deps{Members: nil} -> los dos ops responden 403
func TestNilMemberCheckerDenies(t *testing.T) {
	m := setupModule(t) // setupModule usa Deps{Members: nil}

	ctxGet := mockCtx("user-1", &sitecontent.Content{SiteId: "site-1"})
	m.OpGet(ctxGet)
	if ctxGet.Status != 403 {
		t.Fatalf("expected OpGet status 403 with nil Members, got %d", ctxGet.Status)
	}

	ctxSave := mockCtx("user-1", validContent())
	m.OpSave(ctxSave)
	if ctxSave.Status != 403 {
		t.Fatalf("expected OpSave status 403 with nil Members, got %d", ctxSave.Status)
	}
}

// TestServiceMethodsIgnoreMemberChecker: Get/Save directos siguen funcionando con Members: nil
func TestServiceMethodsIgnoreMemberChecker(t *testing.T) {
	m := setupModule(t)
	c := validContent()
	c.SiteId = "direct-site"

	if err := m.Save(c); err != nil {
		t.Fatalf("expected Save service method to succeed with Members: nil, got: %v", err)
	}

	read, err := m.Get("direct-site")
	if err != nil {
		t.Fatalf("expected Get service method to succeed with Members: nil, got: %v", err)
	}
	if read.SiteId != "direct-site" {
		t.Fatalf("expected SiteId direct-site, got %s", read.SiteId)
	}
}
