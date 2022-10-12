package httplam

const (
	MIMEApplicationJSON            = "application/json"
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationXML             = "application/xml"
	MIMEApplicationXMLCharsetUTF8  = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextHTML                   = "text/html"
	MIMETextHTMLCharsetUTF8        = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                  = "text/plain"
	MIMETextPlainCharsetUTF8       = MIMETextPlain + "; " + charsetUTF8
)

const (
	charsetUTF8 = "charset=UTF-8"
)

var (
	stringContentTypeTable = map[string]bool{
		MIMEApplicationJSON:            true,
		MIMEApplicationJSONCharsetUTF8: true,
		MIMEApplicationXML:             true,
		MIMEApplicationXMLCharsetUTF8:  true,
		MIMETextHTML:                   true,
		MIMETextHTMLCharsetUTF8:        true,
		MIMETextPlain:                  true,
		MIMETextPlainCharsetUTF8:       true,
	}
)

const (
	HeaderContentType        = "Content-Type"
	HeaderSetCookie          = "Set-Cookie"
	HeaderXForwardedProto    = "X-Forwarded-Proto"
	HeaderXForwardedProtocol = "X-Forwarded-Protocol"
	HeaderXForwardedSsl      = "X-Forwarded-Ssl"
	HeaderXUrlScheme         = "X-Url-Scheme"
	HeaderHost               = "Host"
)
