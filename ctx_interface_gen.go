// Code generated by ifacemaker; DO NOT EDIT.

package fiber

import (
	"bufio"
	"crypto/tls"
	"io"
	"mime/multipart"
	"time"

	"github.com/valyala/fasthttp"
)

// Ctx represents the Context which hold the HTTP request and response.
// It has methods for the request query string, parameters, body, HTTP headers and so on.
type Ctx interface {
	// App returns the *App reference to the instance of the Fiber application
	App() *App
	// BaseURL returns (protocol + host + base path).
	BaseURL() string
	// RequestCtx returns *fasthttp.RequestCtx that carries a deadline
	// a cancellation signal, and other values across API boundaries.
	RequestCtx() *fasthttp.RequestCtx
	// Deadline returns the time when work done on behalf of this context
	// should be canceled. Deadline returns ok==false when no deadline is
	// set. Successive calls to Deadline return the same results.
	//
	// Due to current limitations in how fasthttp works, Deadline operates as a nop.
	// See: https://github.com/valyala/fasthttp/issues/965#issuecomment-777268945
	Deadline() (time.Time, bool)
	// Done returns a channel that's closed when work done on behalf of this
	// context should be canceled. Done may return nil if this context can
	// never be canceled. Successive calls to Done return the same value.
	// The close of the Done channel may happen asynchronously,
	// after the cancel function returns.
	//
	// Due to current limitations in how fasthttp works, Done operates as a nop.
	// See: https://github.com/valyala/fasthttp/issues/965#issuecomment-777268945
	Done() <-chan struct{}
	// If Done is not yet closed, Err returns nil.
	// If Done is closed, Err returns a non-nil error explaining why:
	// context.DeadlineExceeded if the context's deadline passed,
	// or context.Canceled if the context was canceled for some other reason.
	// After Err returns a non-nil error, successive calls to Err return the same error.
	//
	// Due to current limitations in how fasthttp works, Err operates as a nop.
	// See: https://github.com/valyala/fasthttp/issues/965#issuecomment-777268945
	Err() error
	// Request return the *fasthttp.Request object
	// This allows you to use all fasthttp request methods
	// https://godoc.org/github.com/valyala/fasthttp#Request
	Request() *fasthttp.Request
	// Response return the *fasthttp.Response object
	// This allows you to use all fasthttp response methods
	// https://godoc.org/github.com/valyala/fasthttp#Response
	Response() *fasthttp.Response
	// Get returns the HTTP request header specified by field.
	// Field names are case-insensitive
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	Get(key string, defaultValue ...string) string
	// GetHeaders returns the HTTP request headers.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetHeaders() map[string][]string
	// GetReqHeaders returns the HTTP request headers.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetReqHeaders() map[string][]string
	// GetRespHeader returns the HTTP response header specified by field.
	// Field names are case-insensitive
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetRespHeader(key string, defaultValue ...string) string
	// GetRespHeaders returns the HTTP response headers.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	GetRespHeaders() map[string][]string
	// ClientHelloInfo return CHI from context
	ClientHelloInfo() *tls.ClientHelloInfo
	// Next executes the next method in the stack that matches the current route.
	Next() error
	// RestartRouting instead of going to the next handler. This may be useful after
	// changing the request path. Note that handlers might be executed again.
	RestartRouting() error
	// OriginalURL contains the original request URL.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	OriginalURL() string
	// Path returns the path part of the request URL.
	// Optionally, you could override the path.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Path(override ...string) string
	// Req returns a convenience type whose API is limited to operations
	// on the incoming request.
	Req() Req
	// Res returns a convenience type whose API is limited to operations
	// on the outgoing response.
	Res() Res
	// Redirect returns the Redirect reference.
	// Use Redirect().Status() to set custom redirection status code.
	// If status is not specified, status defaults to 303 See Other.
	// You can use Redirect().To(), Redirect().Route() and Redirect().Back() for redirection.
	Redirect() *Redirect
	// ViewBind Add vars to default view var map binding to template engine.
	// Variables are read by the Render method and may be overwritten.
	ViewBind(vars Map) error
	// Route returns the matched Route struct.
	Route() *Route
	// SaveFile saves any multipart file to disk.
	SaveFile(fileheader *multipart.FileHeader, path string) error
	// SaveFileToStorage saves any multipart file to an external storage system.
	SaveFileToStorage(fileheader *multipart.FileHeader, path string, storage Storage) error
	// Secure returns whether a secure connection was established.
	Secure() bool
	// Status sets the HTTP status for the response.
	// This method is chainable.
	Status(status int) Ctx
	// String returns unique string representation of the ctx.
	//
	// The returned value may be useful for logging.
	String() string
	// Value makes it possible to retrieve values (Locals) under keys scoped to the request
	// and therefore available to all following routes that match the request.
	Value(key any) any
	// XHR returns a Boolean property, that is true, if the request's X-Requested-With header field is XMLHttpRequest,
	// indicating that the request was issued by a client library (such as jQuery).
	XHR() bool
	// configDependentPaths set paths for route recognition and prepared paths for the user,
	// here the features for caseSensitive, decoded paths, strict paths are evaluated
	configDependentPaths()
	// Reset is a method to reset context fields by given request when to use server handlers.
	Reset(fctx *fasthttp.RequestCtx)
	// Release is a method to reset context fields when to use ReleaseCtx()
	release()
	renderExtensions(bind any)
	// Bind You can bind body, cookie, headers etc. into the map, map slice, struct easily by using Binding method.
	// It gives custom binding support, detailed binding options and more.
	// Replacement of: BodyParser, ParamsParser, GetReqHeaders, GetRespHeaders, AllParams, QueryParser, ReqHeaderParser
	Bind() *Bind
	// Methods to use with next stack.
	getMethodInt() int
	setMethodInt(methodInt int)
	getIndexRoute() int
	getTreePathHash() int
	getDetectionPath() string
	getValues() *[maxParams]string
	getMatched() bool
	setIndexHandler(handler int)
	setIndexRoute(route int)
	setMatched(matched bool)
	setRoute(route *Route)
	keepOriginalPath()
	getPathOriginal() string
	// Accepts checks if the specified extensions or content types are acceptable.
	Accepts(offers ...string) string
	// AcceptsCharsets checks if the specified charset is acceptable.
	AcceptsCharsets(offers ...string) string
	// AcceptsEncodings checks if the specified encoding is acceptable.
	AcceptsEncodings(offers ...string) string
	// AcceptsLanguages checks if the specified language is acceptable.
	AcceptsLanguages(offers ...string) string
	// BodyRaw contains the raw body submitted in a POST request.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	BodyRaw() []byte
	tryDecodeBodyInOrder(originalBody *[]byte, encodings []string) ([]byte, uint8, error)
	// Body contains the raw body submitted in a POST request.
	// This method will decompress the body if the 'Content-Encoding' header is provided.
	// It returns the original (or decompressed) body data which is valid only within the handler.
	// Don't store direct references to the returned data.
	// If you need to keep the body's data later, make a copy or use the Immutable option.
	Body() []byte
	// Cookies are used for getting a cookie value by key.
	// Defaults to the empty string "" if the cookie doesn't exist.
	// If a default value is given, it will return that value if the cookie doesn't exist.
	// The returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Cookies(key string, defaultValue ...string) string
	// FormFile returns the first file by key from a MultipartForm.
	FormFile(key string) (*multipart.FileHeader, error)
	// FormValue returns the first value by key from a MultipartForm.
	// Search is performed in QueryArgs, PostArgs, MultipartForm and FormFile in this particular order.
	// Defaults to the empty string "" if the form value doesn't exist.
	// If a default value is given, it will return that value if the form value does not exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting instead.
	FormValue(key string, defaultValue ...string) string
	// Fresh returns true when the response is still “fresh” in the client's cache,
	// otherwise false is returned to indicate that the client cache is now stale
	// and the full response should be sent.
	// When a client sends the Cache-Control: no-cache request header to indicate an end-to-end
	// reload request, this module will return false to make handling these requests transparent.
	// https://github.com/jshttp/fresh/blob/master/index.js#L33
	Fresh() bool
	// Host contains the host derived from the X-Forwarded-Host or Host HTTP header.
	// Returned value is only valid within the handler. Do not store any references.
	// In a network context, `Host` refers to the combination of a hostname and potentially a port number used for connecting,
	// while `Hostname` refers specifically to the name assigned to a device on a network, excluding any port information.
	// Example: URL: https://example.com:8080 -> Host: example.com:8080
	// Make copies or use the Immutable setting instead.
	// Please use Config.TrustProxy to prevent header spoofing, in case when your app is behind the proxy.
	Host() string
	// Hostname contains the hostname derived from the X-Forwarded-Host or Host HTTP header using the c.Host() method.
	// Returned value is only valid within the handler. Do not store any references.
	// Example: URL: https://example.com:8080 -> Hostname: example.com
	// Make copies or use the Immutable setting instead.
	// Please use Config.TrustProxy to prevent header spoofing, in case when your app is behind the proxy.
	Hostname() string
	// Port returns the remote port of the request.
	Port() string
	// IP returns the remote IP address of the request.
	// If ProxyHeader and IP Validation is configured, it will parse that header and return the first valid IP address.
	// Please use Config.TrustProxy to prevent header spoofing, in case when your app is behind the proxy.
	IP() string
	// extractIPsFromHeader will return a slice of IPs it found given a header name in the order they appear.
	// When IP validation is enabled, any invalid IPs will be omitted.
	extractIPsFromHeader(header string) []string
	// extractIPFromHeader will attempt to pull the real client IP from the given header when IP validation is enabled.
	// currently, it will return the first valid IP address in header.
	// when IP validation is disabled, it will simply return the value of the header without any inspection.
	// Implementation is almost the same as in extractIPsFromHeader, but without allocation of []string.
	extractIPFromHeader(header string) string
	// IPs returns a string slice of IP addresses specified in the X-Forwarded-For request header.
	// When IP validation is enabled, only valid IPs are returned.
	IPs() []string
	// Is returns the matching content type,
	// if the incoming request's Content-Type HTTP header field matches the MIME type specified by the type parameter
	Is(extension string) bool
	// Locals makes it possible to pass any values under keys scoped to the request
	// and therefore available to all following routes that match the request.
	//
	// All the values are removed from ctx after returning from the top
	// RequestHandler. Additionally, Close method is called on each value
	// implementing io.Closer before removing the value from ctx.
	Locals(key any, value ...any) any
	// Method returns the HTTP request method for the context, optionally overridden by the provided argument.
	// If no override is given or if the provided override is not a valid HTTP method, it returns the current method from the context.
	// Otherwise, it updates the context's method and returns the overridden method as a string.
	Method(override ...string) string
	// MultipartForm parse form entries from binary.
	// This returns a map[string][]string, so given a key the value will be a string slice.
	MultipartForm() (*multipart.Form, error)
	// Params is used to get the route parameters.
	// Defaults to empty string "" if the param doesn't exist.
	// If a default value is given, it will return that value if the param doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Params(key string, defaultValue ...string) string
	// Scheme contains the request protocol string: http or https for TLS requests.
	// Please use Config.TrustProxy to prevent header spoofing, in case when your app is behind the proxy.
	Scheme() string
	// Protocol returns the HTTP protocol of request: HTTP/1.1 and HTTP/2.
	Protocol() string
	// Query returns the query string parameter in the url.
	// Defaults to empty string "" if the query doesn't exist.
	// If a default value is given, it will return that value if the query doesn't exist.
	// Returned value is only valid within the handler. Do not store any references.
	// Make copies or use the Immutable setting to use the value outside the Handler.
	Query(key string, defaultValue ...string) string
	// Queries returns a map of query parameters and their values.
	//
	// GET /?name=alex&wanna_cake=2&id=
	// Queries()["name"] == "alex"
	// Queries()["wanna_cake"] == "2"
	// Queries()["id"] == ""
	//
	// GET /?field1=value1&field1=value2&field2=value3
	// Queries()["field1"] == "value2"
	// Queries()["field2"] == "value3"
	//
	// GET /?list_a=1&list_a=2&list_a=3&list_b[]=1&list_b[]=2&list_b[]=3&list_c=1,2,3
	// Queries()["list_a"] == "3"
	// Queries()["list_b[]"] == "3"
	// Queries()["list_c"] == "1,2,3"
	//
	// GET /api/search?filters.author.name=John&filters.category.name=Technology&filters[customer][name]=Alice&filters[status]=pending
	// Queries()["filters.author.name"] == "John"
	// Queries()["filters.category.name"] == "Technology"
	// Queries()["filters[customer][name]"] == "Alice"
	// Queries()["filters[status]"] == "pending"
	Queries() map[string]string
	// Range returns a struct containing the type and a slice of ranges.
	Range(size int) (Range, error)
	// Subdomains returns a slice of subdomains from the host, excluding the last `offset` components.
	// If the offset is negative or exceeds the number of subdomains, an empty slice is returned.
	// If the offset is zero every label (no trimming) is returned.
	Subdomains(offset ...int) []string
	// Stale returns the inverse of Fresh, indicating if the client's cached response is considered stale.
	Stale() bool
	// IsProxyTrusted checks trustworthiness of remote ip.
	// If Config.TrustProxy false, it returns true
	// IsProxyTrusted can check remote ip by proxy ranges and ip map.
	IsProxyTrusted() bool
	// IsFromLocal will return true if request came from local.
	IsFromLocal() bool
	getBody() []byte
	// Append the specified value to the HTTP response header field.
	// If the header is not already set, it creates the header with the specified value.
	Append(field string, values ...string)
	// Attachment sets the HTTP response Content-Disposition header field to attachment.
	Attachment(filename ...string)
	// ClearCookie expires a specific cookie by key on the client side.
	// If no key is provided it expires all cookies that came with the request.
	ClearCookie(key ...string)
	// Cookie sets a cookie by passing a cookie struct.
	Cookie(cookie *Cookie)
	// Download transfers the file from path as an attachment.
	// Typically, browsers will prompt the user for download.
	// By default, the Content-Disposition header filename= parameter is the filepath (this typically appears in the browser dialog).
	// Override this default with the filename parameter.
	Download(file string, filename ...string) error
	// Format performs content-negotiation on the Accept HTTP header.
	// It uses Accepts to select a proper format and calls the matching
	// user-provided handler function.
	// If no accepted format is found, and a format with MediaType "default" is given,
	// that default handler is called. If no format is found and no default is given,
	// StatusNotAcceptable is sent.
	Format(handlers ...ResFmt) error
	// AutoFormat performs content-negotiation on the Accept HTTP header.
	// It uses Accepts to select a proper format.
	// The supported content types are text/html, text/plain, application/json, and application/xml.
	// For more flexible content negotiation, use Format.
	// If the header is not specified or there is no proper format, text/plain is used.
	AutoFormat(body any) error
	// JSON converts any interface or string to JSON.
	// Array and slice values encode as JSON arrays,
	// except that []byte encodes as a base64-encoded string,
	// and a nil slice encodes as the null JSON value.
	// If the ctype parameter is given, this method will set the
	// Content-Type header equal to ctype. If ctype is not given,
	// The Content-Type header will be set to application/json; charset=utf-8.
	JSON(data any, ctype ...string) error
	// MsgPack converts any interface or string to MessagePack encoded bytes.
	// If the ctype parameter is given, this method will set the
	// Content-Type header equal to ctype. If ctype is not given,
	// The Content-Type header will be set to application/vnd.msgpack.
	MsgPack(data any, ctype ...string) error
	// CBOR converts any interface or string to CBOR encoded bytes.
	// If the ctype parameter is given, this method will set the
	// Content-Type header equal to ctype. If ctype is not given,
	// The Content-Type header will be set to application/cbor.
	CBOR(data any, ctype ...string) error
	// JSONP sends a JSON response with JSONP support.
	// This method is identical to JSON, except that it opts-in to JSONP callback support.
	// By default, the callback name is simply callback.
	JSONP(data any, callback ...string) error
	// XML converts any interface or string to XML.
	// This method also sets the content header to application/xml; charset=utf-8.
	XML(data any) error
	// Links joins the links followed by the property to populate the response's Link HTTP header field.
	Links(link ...string)
	// Location sets the response Location HTTP header to the specified path parameter.
	Location(path string)
	// getLocationFromRoute get URL location from route using parameters
	getLocationFromRoute(route Route, params Map) (string, error)
	// GetRouteURL generates URLs to named routes, with parameters. URLs are relative, for example: "/user/1831"
	GetRouteURL(routeName string, params Map) (string, error)
	// Render a template with data and sends a text/html response.
	// We support the following engines: https://github.com/gofiber/template
	Render(name string, bind any, layouts ...string) error
	// Send sets the HTTP response body without copying it.
	// From this point onward the body argument must not be changed.
	Send(body []byte) error
	// SendFile transfers the file from the specified path.
	// By default, the file is not compressed. To enable compression, set SendFile.Compress to true.
	// The Content-Type response HTTP header field is set based on the file's extension.
	// If the file extension is missing or invalid, the Content-Type is detected from the file's format.
	SendFile(file string, config ...SendFile) error
	// SendStatus sets the HTTP status code and if the response body is empty,
	// it sets the correct status message in the body.
	SendStatus(status int) error
	// SendString sets the HTTP response body for string types.
	// This means no type assertion, recommended for faster performance
	SendString(body string) error
	// SendStream sets response body stream and optional body size.
	SendStream(stream io.Reader, size ...int) error
	// SendStreamWriter sets response body stream writer
	SendStreamWriter(streamWriter func(*bufio.Writer)) error
	// Set sets the response's HTTP header field to the specified key, value.
	Set(key, val string)
	setCanonical(key, val string)
	// Type sets the Content-Type HTTP header to the MIME type specified by the file extension.
	Type(extension string, charset ...string) Ctx
	// Vary adds the given header field to the Vary response header.
	// This will append the header, if not already listed, otherwise leaves it listed in the current location.
	Vary(fields ...string)
	// Write appends p into response body.
	Write(p []byte) (int, error)
	// Writef appends f & a into response body writer.
	Writef(f string, a ...any) (int, error)
	// WriteString appends s to response body.
	WriteString(s string) (int, error)
	// Drop closes the underlying connection without sending any response headers or body.
	// This can be useful for silently terminating client connections, such as in DDoS mitigation
	// or when blocking access to sensitive endpoints.
	Drop() error
	// End immediately flushes the current response and closes the underlying connection.
	End() error
}
