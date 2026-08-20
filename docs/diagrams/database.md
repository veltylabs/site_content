flowchart TD
    Content["Content<br/>SiteID (PK)"]
    Brand["Brand<br/>Name, WideLogo, CompactLogo, LogoAlt, PrimaryColor"]
    Contact["Contact<br/>Phone, Email, Address"]
    Hero["Hero<br/>Title, Subtitle, Images"]
    Link["Link<br/>Text, URL"]
    About["About<br/>Title, Body, Image, Mission, Vision"]
    Service["Service<br/>Slug (Unique), Title (Unique), Description (Unique), Image, Body"]
    Stat["Stat<br/>Value, Label"]
    Schedule["Schedule<br/>Days, Hours"]
    Map["Map<br/>EmbedURL"]
    SEO["SEO<br/>Description, SocialImage, SchemaType"]
    ImageRef["ImageRef<br/>Key (R2), Alt, Usage"]

    Content --> Brand
    Content --> Contact
    Content --> Hero
    Hero --> Link
    Content --> About
    Content --> Service
    Content --> Stat
    Content --> Schedule
    Content --> Map
    Content --> SEO
    Content --> ImageRef
