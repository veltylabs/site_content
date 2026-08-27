package sitecontent

// MemberChecker responde si un usuario puede tocar el contenido de un sitio.
//
// Es una interfaz LECTORA ESTRECHA, declarada aquí y satisfecha
// estructuralmente por *sitemanager.Module: los módulos no se importan entre
// sí, y la raíz de composición es quien conecta las instancias (patrón
// CatalogReader del AGENTS.md canónico).
//
// nil DENIEGA. La ausencia de un comprobador no es un permiso — es el mismo
// criterio que model.Allowed aplica a un Authorizer ausente.
type MemberChecker interface {
	// CanEditContent informa si userID puede leer y escribir el contenido de
	// siteID. Una sola pregunta, no un rol: qué significa "puede" es del
	// módulo que conoce las membresías, no de este.
	CanEditContent(userID, siteID string) bool
}
