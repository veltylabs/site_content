package sitecontent

import (
	"webtyp.com/ddl"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/router"
)

var _ router.OpModule = (*Module)(nil)

type ContentRecord struct {
	SiteId string
	Data   string
}

func (m *ContentRecord) ModelName() string { return "content" }

var contentRecordFields = []model.Field{
	{Name: "SiteID", Type: model.Text(), DB: &model.FieldDB{PK: true}, NotNull: true},
	{Name: "Data", Type: model.Text()},
}

func (m *ContentRecord) Schema() []model.Field { return contentRecordFields }
func (m *ContentRecord) Pointers() []any       { return []any{&m.SiteId, &m.Data} }
func (m *ContentRecord) IsNil() bool          { return m == nil }

func (m *ContentRecord) EncodeFields(w model.FieldWriter) {
	w.String("SiteID", m.SiteId)
	w.String("Data", m.Data)
}

func (m *ContentRecord) DecodeFields(r model.FieldReader) {
	if v, ok := r.String("SiteID"); ok {
		m.SiteId = v
	}
	if v, ok := r.String("Data"); ok {
		m.Data = v
	}
}

func (m *ContentRecord) Validate(action byte) error {
	return model.ValidateFields(action, m)
}

type Deps struct {
	DB      *orm.DB
	IDs     model.IDGenerator
	Members MemberChecker // nil DENIEGA en los ops; los métodos de servicio no lo consultan
}

type Module struct {
	db      *orm.DB
	ids     model.IDGenerator
	members MemberChecker
}

func New(d Deps) (*Module, error) {
	if d.DB == nil {
		return nil, fmt.Err("site_content", "DB required")
	}
	if d.IDs == nil {
		return nil, fmt.Err("site_content", "IDs required")
	}

	m := &Module{
		db:      d.DB,
		ids:     d.IDs,
		members: d.Members,
	}

	if compiler, ok := d.DB.RawConn().(ddl.Compiler); ok {
		if err := ddl.New(d.DB.RawConn(), compiler).CreateTable(&ContentRecord{}); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Module) ModelName() string {
	return "site_content"
}

func (m *Module) MountOps(reg router.OpRegistry) {
	reg.Op("get", m.OpGet).
		Requires(model.Resource("site_content"), model.Read).
		Accepts(&Content{})

	reg.Op("save", m.OpSave).
		Requires(model.Resource("site_content"), model.Create|model.Update).
		Accepts(&Content{})
}

func (m *Module) Get(siteID string) (*Content, error) {
	if siteID == "" {
		return nil, fmt.Err("site_content", "siteID required")
	}
	var rec ContentRecord
	rec.SiteId = siteID
	err := m.db.Query(&rec).Where("SiteID").Eq(siteID).ReadOne()
	if err != nil {
		if orm.IsNotFound(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	res := &Content{SiteId: rec.SiteId}
	mc := MemoryCodec{}
	mc.Decode(rec.Data, res)
	res.SiteId = rec.SiteId
	return res, nil
}

func (m *Module) Save(c *Content) error {
	if err := Validate(c); err != nil {
		return err
	}

	mc := MemoryCodec{}
	encodedData := mc.Encode(c)

	rec := &ContentRecord{
		SiteId: c.SiteId,
		Data:   encodedData,
	}

	existing, err := m.Get(c.SiteId)
	if err != nil {
		if e, ok := err.(domainError); !ok || e != ErrNotFound {
			return err
		}
	}

	if existing != nil {
		return m.db.Update(rec, orm.Eq("SiteID", c.SiteId))
	}
	return m.db.Create(rec)
}

// requireMember es el único sitio donde se decide si un llamante puede tocar
// el contenido de un sitio. Una sola implementación: repetirla en cada op
// garantiza que algún día falte en uno.
func (m *Module) requireMember(ctx router.Context, siteID string) bool {
	userID := ctx.UserID()
	if userID == "" {
		ctx.WriteStatus(401)
		return false
	}
	if siteID == "" {
		ctx.WriteStatus(400)
		return false
	}
	// nil deniega: sin comprobador no hay permiso.
	if m.members == nil || !m.members.CanEditContent(userID, siteID) {
		ctx.WriteStatus(403)
		return false
	}
	return true
}

func (m *Module) OpGet(ctx router.Context) {
	var args Content
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	if !m.requireMember(ctx, args.SiteId) {
		return
	}

	res, err := m.Get(args.SiteId)
	if err != nil {
		if e, ok := err.(domainError); ok && e == ErrNotFound {
			ctx.WriteStatus(404)
			return
		}
		ctx.WriteStatus(500)
		return
	}

	ctx.WriteStatus(200)
	_ = ctx.Encode(res)
}

func (m *Module) OpSave(ctx router.Context) {
	var args Content
	if err := ctx.Decode(&args); err != nil {
		ctx.WriteStatus(400)
		return
	}
	if !m.requireMember(ctx, args.SiteId) {
		return
	}

	if err := m.Save(&args); err != nil {
		if e, ok := err.(domainError); ok && e == ErrNotFound {
			ctx.WriteStatus(404)
			return
		}
		ctx.WriteStatus(400)
		return
	}

	ctx.WriteStatus(200)
	_ = ctx.Encode(&args)
}
