// Package httpheader defines HTTP header constants.
package httpheader

const (
	// Accept is the constant of the Accept header name.
	Accept = "Accept"
	// AcceptCharset is the constant of the Accept-Charset header name.
	AcceptCharset = "Accept-Charset"
	// AcceptEncoding is the constant of the Accept-Encoding header name.
	AcceptEncoding = "Accept-Encoding"
	// AcceptLanguage is the constant of the Accept-Language header name.
	AcceptLanguage = "Accept-Language"
	// Authorization is the constant of the Authorization header name.
	Authorization = "Authorization"
	// CacheControl is the constant of the Cache-Control header name.
	CacheControl = "Cache-Control"
	// ContentLength is the constant of the Content-Length header name.
	ContentLength = "Content-Length"
	// ContentMD5 is the constant of the Content-MD5 header name.
	ContentMD5 = "Content-MD5"
	// ContentType is the constant for the Content-Type header name.
	ContentType = "Content-Type"
	// DoNotTrack is the constant of the DNT header name.
	DoNotTrack = "DNT"
	// IfMatch is the constant of the If-Match header name.
	IfMatch = "If-Match"
	// IfModifiedSince is the constant of the If-Modified-Since header name.
	IfModifiedSince = "If-Modified-Since"
	// IfNoneMatch is the constant of the If-None-Match header name.
	IfNoneMatch = "If-None-Match"
	// IfRange is the constant of the If-Range header name.
	IfRange = "If-Range"
	// IfUnmodifiedSince is the constant of the If-Unmodified-Since header name.
	IfUnmodifiedSince = "If-Unmodified-Since"
	// MaxForwards is the constant of the Max-Forwards header name.
	MaxForwards = "Max-Forwards"
	// ProxyAuthorization is the constant of the Proxy-Authorization header name.
	ProxyAuthorization = "Proxy-Authorization"
	// Pragma is the constant of the Pragma header name.
	Pragma = "Pragma"
	// Range is the constant of the Range header name.
	Range = "Range"
	// Referer is the constant of the Referer header name.
	Referer = "Referer"
	// UserAgent is the constant of the User-Agent header name.
	UserAgent = "User-Agent"
	// TE is the constant of the TE header name.
	TE = "TE"
	// Via is the constant of the Via header name.
	Via = "Via"
	// Warning is the constant of the Warning header name.
	Warning = "Warning"
	// Cookie is the constant of the Cookie header name.
	Cookie = "Cookie"
	// Origin is the constant of the Origin header name.
	Origin = "Origin"
	// AcceptDatetime is the constant of the Accept-Datetime header name.
	AcceptDatetime = "Accept-Datetime"
	// XRequestedWith is the constant of the X-Requested-With header name.
	XRequestedWith = "X-Requested-With"
	// AccessControlAllowOrigin is the constant of the Access-Control-Allow-Origin header name.
	AccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// AccessControlAllowMethods is the constant of the Access-Control-Allow-Methods header name.
	AccessControlAllowMethods = "Access-Control-Allow-Methods"
	// AccessControlAllowHeaders is the constant of the Access-Control-Allow-Headers header name.
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// AccessControlAllowCredentials is the constant of the Access-Control-Allow-Credentials header name.
	AccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// AccessControlExposeHeaders is the constant of the Access-Control-Expose-Headers header name.
	AccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// AccessControlMaxAge is the constant of the Access-Control-Max-Age header name.
	AccessControlMaxAge = "Access-Control-Max-Age"
	// AccessControlRequestMethod is the constant of the Access-Control-Request-Method header name.
	AccessControlRequestMethod = "Access-Control-Request-Method"
	// AccessControlRequestHeaders is the constant of the Access-Control-Request-Headers header name.
	AccessControlRequestHeaders = "Access-Control-Request-Headers"
	// AcceptPatch is the constant of the Accept-Patch header name.
	AcceptPatch = "Accept-Patch"
	// AcceptRanges is the constant of the Accept-Ranges header name.
	AcceptRanges = "Accept-Ranges"
	// Allow is the constant of the Allow header name.
	Allow = "Allow"
	// ContentEncoding is the constant of the Content-Encoding header name.
	ContentEncoding = "Content-Encoding"
	// ContentLanguage is the constant of the Content-Language header name.
	ContentLanguage = "Content-Language"
	// ContentLocation is the constant of the Content-Location header name.
	ContentLocation = "Content-Location"
	// ContentDisposition is the constant of the Content-Disposition header name.
	ContentDisposition = "Content-Disposition"
	// ContentRange is the constant of the Content-Range header name.
	ContentRange = "Content-Range"
	// ETag is the constant of the ETag header name.
	ETag = "ETag"
	// Expires is the constant of the Expires header name.
	Expires = "Expires"
	// LastModified is the constant of the Last-Modified header name.
	LastModified = "Last-Modified"
	// Link is the constant of the Link header name.
	Link = "Link"
	// Location is the constant of the Location header name.
	Location = "Location"
	// P3P is the constant of the P3P header name.
	P3P = "P3P"
	// ProxyAuthenticate is the constant of the Proxy-Authenticate header name.
	ProxyAuthenticate = "Proxy-Authenticate"
	// RetryAfter is the constant of the Retry-After header name.
	RetryAfter = "Retry-After"
	// Server is the constant of the Server header name.
	Server = "Server"
	// SetCookie is the constant of the Set-Cookie header name.
	SetCookie = "Set-Cookie"
	// StrictTransportSecurity is the constant of the Strict-Transport-Security header name.
	StrictTransportSecurity = "Strict-Transport-Security"
	// TransferEncoding is the constant of the Transfer-Encoding header name.
	TransferEncoding = "Transfer-Encoding"
	// Upgrade is the constant of the Upgrade header name.
	Upgrade = "Upgrade"
	// Vary is the constant of the Vary header name.
	Vary = "Vary"
	// WWWAuthenticate is the constant of the WWW-Authenticate header name.
	WWWAuthenticate = "WWW-Authenticate"

	// XFrameOptions is the constant of the X-Frame-Options header name.
	XFrameOptions = "X-Frame-Options"
	// XXSSProtection is the constant of the X-XSS-Protection header name.
	XXSSProtection = "X-XSS-Protection"
	// ContentSecurityPolicy is the constant of the Content-Security-Policy header name.
	ContentSecurityPolicy = "Content-Security-Policy"
	// XContentSecurityPolicy is the constant of the X-Content-Security-Policy header name.
	XContentSecurityPolicy = "X-Content-Security-Policy"
	// XWebKitCSP is the constant of the X-WebKit-CSP header name.
	XWebKitCSP = "X-WebKit-CSP"
	// XContentTypeOptions is the constant of the X-Content-Type-Options header name.
	XContentTypeOptions = "X-Content-Type-Options"
	// XPoweredBy is the constant of the X-Powered-By header name.
	XPoweredBy = "X-Powered-By"
	// XUACompatible is the constant of the X-UA-Compatible header name.
	XUACompatible = "X-UA-Compatible"
	// XForwardedProto is the constant of the X-Forwarded-Proto header
	// that is used to identify the protocol (HTTP or HTTPS) that a visitor used to connect to the proxy server.
	XForwardedProto = "X-Forwarded-Proto"
	// XHTTPMethodOverride is the constant of the X-HTTP-Method-Override header.
	XHTTPMethodOverride = "X-HTTP-Method-Override"
	// XForwardedFor is the constant of the X-Forwarded-For header
	// that maintains proxy server and original visitor IP addresses.
	XForwardedFor = "X-Forwarded-For"
	// XRealIP is the constant of the X-Real-IP header.
	XRealIP = "X-Real-IP"
	// XCSRFToken is the constant of the X-CSRF-Token header.
	XCSRFToken = "X-CSRF-Token" //nolint:gosec
	// XRatelimitLimit is the constant of the X-Ratelimit-Limit header.
	XRatelimitLimit = "X-Ratelimit-Limit"
	// XRatelimitRemaining is the constant of the X-Ratelimit-Remaining header.
	XRatelimitRemaining = "X-Ratelimit-Remaining"
	// XRatelimitReset is the constant of the X-Ratelimit-Reset header.
	XRatelimitReset = "X-Ratelimit-Reset"

	// Cloudflare headers.

	// CFConnectingIP is the constant of the CF-Connecting-IP header that provides the client IP address connecting to Cloudflare to the origin web server.
	// This header will only be sent on the traffic from Cloudflare's edge to your origin web server.
	CFConnectingIP = "CF-Connecting-IP"
	// CFConnectingIPv6 is the constant of the CF-Connecting-IPv6 header.
	CFConnectingIPv6 = "CF-Connecting-IPv6"
	// CFEWVia is the constant of the CF-EW-Via header that is used for loop detection, similar to the CDN-Loop.
	CFEWVia = "CF-EW-Via"
	// TrueClientIP is the constant of the True-Client-IP header that provides the original client IP address to the origin web server.
	TrueClientIP = "True-Client-IP"
	// CFRay is the constant of the Cf-Ray header.
	// It is a hashed value that encodes information about the data center and the visitor's request.
	CFRay = "Cf-Ray"
	// CFIPCountry is the constant of the CF-IPCountry header that contains a two-character country code of the originating visitor's country.
	CFIPCountry = "CF-IPCountry"
	// CFVisitor is the constant of the CF-Visitor header that contains the scheme information.
	CFVisitor = "CF-Visitor"
	// CDNLoop is the constant of the CDN-Loop header that allows Cloudflare to specify
	// how many times a request can enter Cloudflare's network before it is blocked as a looping request.
	CDNLoop = "CDN-Loop"
	// CFConnectingO2O is the constant of the CF-Connecting-O2O header.
	// If SSL for SaaS is used for the SaaS provider-owned zone, a HTTP header will be set to cf-connecting-o2o: 1.
	CFConnectingO2O = "CF-Connecting-O2O"
	// CFWorker is the constant of the CF-Worker header.
	// It is added to an edge Worker sub-request that identifies the host that spawned the sub-request.
	CFWorker = "CF-Worker"
)

const (
	// ContentTypeJSON is the constant for the application/json content type.
	ContentTypeJSON = "application/json"
	// ContentTypeNdJSON is the constant for the application/x-ndjson content type.
	ContentTypeNdJSON = "application/x-ndjson"
	// ContentTypeXML is the constant for the application/xml content type.
	ContentTypeXML = "application/xml"
	// ContentTypeFormURLEncoded is the constant for the application/x-www-form-urlencoded content type.
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	// ContentTypeMultipartFormData is the constant for the multipart/form-data content type.
	ContentTypeMultipartFormData = "multipart/form-data"
	// ContentTypeTextPlain is the constant for the text/plain content type.
	ContentTypeTextPlain = "text/plain"
	// ContentTypeTextHTML is the constant for the text/html content type.
	ContentTypeTextHTML = "text/html"
	// ContentTypeTextXML is the constant for the text/xml content type.
	ContentTypeTextXML = "text/xml"
	// ContentTypeOctetStream is the constant for the application/octet-stream content type.
	ContentTypeOctetStream = "application/octet-stream"
)
