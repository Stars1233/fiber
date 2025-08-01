---
id: whats_new
title: 🆕 What's New in v3
sidebar_position: 2
toc_max_heading_level: 4
---

[//]: # (https://github.com/gofiber/fiber/releases/tag/v3.0.0-beta.2)

## 🎉 Welcome

We are excited to announce the release of Fiber v3! 🚀

In this guide, we'll walk you through the most important changes in Fiber `v3` and show you how to migrate your existing Fiber `v2` applications to Fiber `v3`.

Here's a quick overview of the changes in Fiber `v3`:

- [🚀 App](#-app)
- [🎣 Hooks](#-hooks)
- [🚀 Listen](#-listen)
- [🗺️ Router](#-router)
- [🧠 Context](#-context)
- [📎 Binding](#-binding)
- [🔄️ Redirect](#-redirect)
- [🌎 Client package](#-client-package)
- [🧰 Generic functions](#-generic-functions)
- [🥡 Services](#-services)
- [📃 Log](#-log)
- [📦 Storage Interface](#-storage-interface)
- [🧬 Middlewares](#-middlewares)
  - [Important Change for Accessing Middleware Data](#important-change-for-accessing-middleware-data)
  - [Adaptor](#adaptor)
  - [BasicAuth](#basicauth)
  - [Cache](#cache)
  - [CORS](#cors)
  - [CSRF](#csrf)
  - [Compression](#compression)
  - [EncryptCookie](#encryptcookie)
  - [Filesystem](#filesystem)
  - [Healthcheck](#healthcheck)
  - [KeyAuth](#keyauth)
  - [Logger](#logger)
  - [Monitor](#monitor)
  - [Proxy](#proxy)
  - [Session](#session)
- [🔌 Addons](#-addons)
- [📋 Migration guide](#-migration-guide)

## Drop for old Go versions

Fiber `v3` drops support for Go versions below `1.24`. We recommend upgrading to Go `1.24` or higher to use Fiber `v3`.

## 🚀 App

We have made several changes to the Fiber app, including:

- **Listen**: The `Listen` method has been unified with the configuration, allowing for more streamlined setup.
- **Static**: The `Static` method has been removed and its functionality has been moved to the [static middleware](./middleware/static.md).
- **app.Config properties**: Several properties have been moved to the listen configuration:
  - `DisableStartupMessage`
  - `EnablePrefork` (previously `Prefork`)
  - `EnablePrintRoutes`
  - `ListenerNetwork` (previously `Network`)
- **Trusted Proxy Configuration**: The `EnabledTrustedProxyCheck` has been moved to `app.Config.TrustProxy`, and `TrustedProxies` has been moved to `TrustProxyConfig.Proxies`.
- **XMLDecoder Config Property**: The `XMLDecoder` property has been added to allow usage of 3rd-party XML libraries in XML binder.

### New Methods

- **RegisterCustomBinder**: Allows for the registration of custom binders.
- **RegisterCustomConstraint**: Allows for the registration of custom constraints.
- **NewWithCustomCtx**: Initialize an app with a custom context in one step.
- **State**: Provides a global state for the application, which can be used to store and retrieve data across the application. Check out the [State](./api/state) method for further details.
- **NewErrorf**: Allows variadic parameters when creating formatted errors.

#### Custom Route Constraints

Custom route constraints enable you to define your own validation rules for route parameters.
Use `RegisterCustomConstraint` to add a constraint type that implements the `CustomConstraint` interface.

<details>
<summary>Example</summary>

```go
type UlidConstraint struct {
    fiber.CustomConstraint
}

func (*UlidConstraint) Name() string {
    return "ulid"
}

func (*UlidConstraint) Execute(param string, args ...string) bool {
    _, err := ulid.Parse(param)
    return err == nil
}

app.RegisterCustomConstraint(&UlidConstraint{})

app.Get("/login/:id<ulid>", func(c fiber.Ctx) error {
    return c.SendString("User " + c.Params("id"))
})
```

</details>

### Removed Methods

- **Mount**: Use `app.Use()` instead.
- **ListenTLS**: Use `app.Listen()` with `tls.Config`.
- **ListenTLSWithCertificate**: Use `app.Listen()` with `tls.Config`.
- **ListenMutualTLS**: Use `app.Listen()` with `tls.Config`.
- **ListenMutualTLSWithCertificate**: Use `app.Listen()` with `tls.Config`.
- **Context()**: Removed. `Ctx` now directly implements `context.Context`, so you can pass `c` anywhere a `context.Context` is required.
- **SetContext()**: Removed. Attach additional context information using `Locals` or middleware if needed.

### Method Changes

- **Test**: The `Test` method has replaced the timeout parameter with a configuration parameter. `0` or lower represents no timeout.
- **Listen**: Now has a configuration parameter.
- **Listener**: Now has a configuration parameter.

### Custom Ctx Interface in Fiber v3

Fiber v3 introduces a customizable `Ctx` interface, allowing developers to extend and modify the context to fit their needs. This feature provides greater flexibility and control over request handling.

#### Idea Behind Custom Ctx Classes

The idea behind custom `Ctx` classes is to give developers the ability to extend the default context with additional methods and properties tailored to the specific requirements of their application. This allows for better request handling and easier implementation of specific logic.

#### NewWithCustomCtx

`NewWithCustomCtx` creates the application and sets the custom context factory at initialization time.

```go title="Signature"
func NewWithCustomCtx(fn func(app *App) CustomCtx, config ...Config) *App
```

<details>
<summary>Example</summary>

```go
package main

import (
    "log"
    "github.com/gofiber/fiber/v3"
)

type CustomCtx struct {
    fiber.Ctx
}

func (c *CustomCtx) CustomMethod() string {
    return "custom value"
}

func main() {
    app := fiber.NewWithCustomCtx(func(app *fiber.App) fiber.Ctx {
        return &CustomCtx{
            Ctx: *fiber.NewCtx(app),
        }
    })

    app.Get("/", func(c fiber.Ctx) error {
        customCtx := c.(*CustomCtx)
        return c.SendString(customCtx.CustomMethod())
    })

    log.Fatal(app.Listen(":3000"))
}
```

This example creates a `CustomCtx` with an extra `CustomMethod` and initializes the app with `NewWithCustomCtx`.

</details>

### Configurable TLS Minimum Version

We have added support for configuring the TLS minimum version. This field allows you to set the TLS minimum version for TLSAutoCert and the server listener.

```go
app.Listen(":444", fiber.ListenConfig{TLSMinVersion: tls.VersionTLS12})
```

#### TLS AutoCert support (ACME / Let's Encrypt)

We have added native support for automatic certificates management from Let's Encrypt and any other ACME-based providers.

```go
// Certificate manager
certManager := &autocert.Manager{
    Prompt: autocert.AcceptTOS,
    // Replace with your domain name
    HostPolicy: autocert.HostWhitelist("example.com"),
    // Folder to store the certificates
    Cache: autocert.DirCache("./certs"),
}

app.Listen(":444", fiber.ListenConfig{
    AutoCertManager:    certManager,
})
```

### MIME Constants

`MIMEApplicationJavaScript` and `MIMEApplicationJavaScriptCharsetUTF8` are deprecated. Use `MIMETextJavaScript` and `MIMETextJavaScriptCharsetUTF8` instead.

## 🎣 Hooks

We have made several changes to the Fiber hooks, including:

- Added new shutdown hooks to provide better control over the shutdown process:
  - `OnPreShutdown` - Executes before the server starts shutting down
  - `OnPostShutdown` - Executes after the server has shut down, receives any shutdown error
- Deprecated `OnShutdown` in favor of the new pre/post shutdown hooks
- Improved shutdown hook execution order and reliability
- Added mutex protection for hook registration and execution

Important: When using shutdown hooks, ensure app.Listen() is called in a separate goroutine:

```go
// Correct usage
go app.Listen(":3000")
// ... register shutdown hooks
app.Shutdown()

// Incorrect usage - hooks won't work
app.Listen(":3000") // This blocks
app.Shutdown()      // Never reached
```

## 🚀 Listen

We have made several changes to the Fiber listen, including:

- Removed `OnShutdownError` and `OnShutdownSuccess` from `ListenerConfig` in favor of using `OnPostShutdown` hook which receives the shutdown error

```go
app := fiber.New()

// Before - using ListenerConfig callbacks
app.Listen(":3000", fiber.ListenerConfig{
    OnShutdownError: func(err error) {
        log.Printf("Shutdown error: %v", err)
    },
    OnShutdownSuccess: func() {
        log.Println("Shutdown successful")
    },
})

// After - using OnPostShutdown hook
app.Hooks().OnPostShutdown(func(err error) error {
    if err != nil {
        log.Printf("Shutdown error: %v", err)
    } else {
        log.Println("Shutdown successful")
    }
    return nil
})
go app.Listen(":3000")
```

This change simplifies the shutdown handling by consolidating the shutdown callbacks into a single hook that receives the error status.

- Added support for Unix domain sockets via `ListenerNetwork` and `UnixSocketFileMode`

```go
// v2 - Requires manual deletion of old file and permissions change
app := fiber.New(fiber.Config{
    Network: "unix",
})

os.Remove("app.sock")
app.Hooks().OnListen(func(fiber.ListenData) error {
    return os.Chmod("app.sock", 0770)
})
app.Listen("app.sock")

// v3 - Fiber does it for you
app := fiber.New()
app.Listen("app.sock", fiber.ListenerConfig{
    ListenerNetwork:    fiber.NetworkUnix,
    UnixSocketFileMode: 0770,
})
```

## 🗺 Router

We have slightly adapted our router interface

### HTTP method registration

In `v2` one handler was already mandatory when the route has been registered, but this was checked at runtime and was not correctly reflected in the signature, this has now been changed in `v3` to make it more explicit.

```diff
-    Get(path string, handlers ...Handler) Router
+    Get(path string, handler Handler, middleware ...Handler) Router
-    Head(path string, handlers ...Handler) Router
+    Head(path string, handler Handler, middleware ...Handler) Router
-    Post(path string, handlers ...Handler) Router
+    Post(path string, handler Handler, middleware ...Handler) Router
-    Put(path string, handlers ...Handler) Router
+    Put(path string, handler Handler, middleware ...Handler) Router
-    Delete(path string, handlers ...Handler) Router
+    Delete(path string, handler Handler, middleware ...Handler) Router
-    Connect(path string, handlers ...Handler) Router
+    Connect(path string, handler Handler, middleware ...Handler) Router
-    Options(path string, handlers ...Handler) Router
+    Options(path string, handler Handler, middleware ...Handler) Router
-    Trace(path string, handlers ...Handler) Router
+    Trace(path string, handler Handler, middleware ...Handler) Router
-    Patch(path string, handlers ...Handler) Router
+    Patch(path string, handler Handler, middleware ...Handler) Router
-    All(path string, handlers ...Handler) Router
+    All(path string, handler Handler, middleware ...Handler) Router
```

### Route chaining

The route method is now like [`Express`](https://expressjs.com/de/api.html#app.route) which gives you the option of a different notation and allows you to concatenate the route declaration.

```diff
-    Route(prefix string, fn func(router Router), name ...string) Router
+    Route(path string) Register
```

<details>
<summary>Example</summary>

```go
app.Route("/api").Route("/user/:id?")
    .Get(func(c fiber.Ctx) error {
        // Get user
        return c.JSON(fiber.Map{"message": "Get user", "id": c.Params("id")})
    })
    .Post(func(c fiber.Ctx) error {
        // Create user
        return c.JSON(fiber.Map{"message": "User created"})
    })
    .Put(func(c fiber.Ctx) error {
        // Update user
        return c.JSON(fiber.Map{"message": "User updated", "id": c.Params("id")})
    })
    .Delete(func(c fiber.Ctx) error {
        // Delete user
        return c.JSON(fiber.Map{"message": "User deleted", "id": c.Params("id")})
    })
```

</details>

You can find more information about `app.Route` in the [API documentation](./api/app#route).

### Middleware registration

We have aligned our method for middlewares closer to [`Express`](https://expressjs.com/de/api.html#app.use) and now also support the [`Use`](./api/app#use) of multiple prefixes.

Registering a subapp is now also possible via the [`Use`](./api/app#use) method instead of the old `app.Mount` method.

<details>
<summary>Example</summary>

```go
// register multiple prefixes
app.Use(["/v1", "/v2"], func(c fiber.Ctx) error {
    // Middleware for /v1 and /v2
    return c.Next()
})

// define subapp
api := fiber.New()
api.Get("/user", func(c fiber.Ctx) error {
    return c.SendString("User")
})
// register subapp
app.Use("/api", api)
```

</details>

To enable the routing changes above we had to slightly adjust the signature of the `Add` method.

```diff
-    Add(method, path string, handlers ...Handler) Router
+    Add(methods []string, path string, handler Handler, middleware ...Handler) Router
```

### Test Config

The `app.Test()` method now allows users to customize their test configurations:

<details>
<summary>Example</summary>

```go
// Create a test app with a handler to test
app := fiber.New()
app.Get("/", func(c fiber.Ctx) {
    return c.SendString("hello world")
})

// Define the HTTP request and custom TestConfig to test the handler
req := httptest.NewRequest(MethodGet, "/", nil)
testConfig := fiber.TestConfig{
    Timeout:       0,
    FailOnTimeout: false,
}

// Test the handler using the request and testConfig
resp, err := app.Test(req, testConfig)
```

</details>

To provide configurable testing capabilities, we had to change
the signature of the `Test` method.

```diff
-    Test(req *http.Request, timeout ...time.Duration) (*http.Response, error)
+    Test(req *http.Request, config ...fiber.TestConfig) (*http.Response, error)
```

The `TestConfig` struct provides the following configuration options:

- `Timeout`: The duration to wait before timing out the test. Use 0 for no timeout.
- `FailOnTimeout`: Controls the behavior when a timeout occurs:
  - When true, the test will return an `os.ErrDeadlineExceeded` if the test exceeds the `Timeout` duration.
  - When false, the test will return the partial response received before timing out.

If a custom `TestConfig` isn't provided, then the following will be used:

```go
testConfig := fiber.TestConfig{
    Timeout:       time.Second,
    FailOnTimeout: true,
}
```

**Note:** Using this default is **NOT** the same as providing an empty `TestConfig` as an argument to `app.Test()`.

An empty `TestConfig` is the equivalent of:

```go
testConfig := fiber.TestConfig{
    Timeout:       0,
    FailOnTimeout: false,
}
```

## 🧠 Context

### New Features

- Cookie now allows Partitioned cookies for [CHIPS](https://developers.google.com/privacy-sandbox/3pcd/chips) support. CHIPS (Cookies Having Independent Partitioned State) is a feature that improves privacy by allowing cookies to be partitioned by top-level site, mitigating cross-site tracking.
- Cookie automatic security enforcement: When setting a cookie with `SameSite=None`, Fiber automatically sets `Secure=true` as required by RFC 6265bis and modern browsers (Chrome, Firefox, Safari). This ensures compliance with the "None" SameSite policy. See [Mozilla docs](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie#none) and [Chrome docs](https://developers.google.com/search/blog/2020/01/get-ready-for-new-samesitenone-secure) for details.
- Context now implements [context.Context](https://pkg.go.dev/context#Context).

### New Methods

- **AutoFormat**: Similar to Express.js, automatically formats the response based on the request's `Accept` header.
- **Deadline**: For implementing `context.Context`.
- **Done**: For implementing `context.Context`.
- **Err**: For implementing `context.Context`.
- **Host**: Similar to Express.js, returns the host name of the request.
- **Port**: Similar to Express.js, returns the port number of the request.
- **IsProxyTrusted**: Checks the trustworthiness of the remote IP.
- **Reset**: Resets context fields for server handlers.
- **Schema**: Similar to Express.js, returns the schema (HTTP or HTTPS) of the request.
- **SendStream**: Similar to Express.js, sends a stream as the response.
- **SendStreamWriter**: Sends a stream using a writer function.
- **SendString**: Similar to Express.js, sends a string as the response.
- **String**: Similar to Express.js, converts a value to a string.
- **Value**: For implementing `context.Context`. Returns request-scoped value from Locals.
- **ViewBind**: Binds data to a view, replacing the old `Bind` method.
- **CBOR**: Introducing [CBOR](https://cbor.io/) binary encoding format for both request & response body. CBOR is a binary data serialization format which is both compact and efficient, making it ideal for use in web applications.
- **MsgPack**: Introducing [MsgPack](https://msgpack.org/) binary encoding format for both request & response body. MsgPack is a binary serialization format that is more efficient than JSON, making it ideal for high-performance applications.
- **Drop**: Terminates the client connection silently without sending any HTTP headers or response body. This can be used for scenarios where you want to block certain requests without notifying the client, such as mitigating DDoS attacks or protecting sensitive endpoints from unauthorized access.
- **End**: Similar to Express.js, immediately flushes the current response and closes the underlying connection.

### Removed Methods

- **AllParams**: Use `c.Bind().URI()` instead.
- **ParamsInt**: Use `Params` with generic types.
- **QueryBool**: Use `Query` with generic types.
- **QueryFloat**: Use `Query` with generic types.
- **QueryInt**: Use `Query` with generic types.
- **BodyParser**: Use `c.Bind().Body()` instead.
- **CookieParser**: Use `c.Bind().Cookie()` instead.
- **ParamsParser**: Use `c.Bind().URI()` instead.
- **RedirectToRoute**: Use `c.Redirect().Route()` instead.
- **RedirectBack**: Use `c.Redirect().Back()` instead.
- **ReqHeaderParser**: Use `c.Bind().Header()` instead.

### Changed Methods

- **Bind**: Now used for binding instead of view binding. Use `c.ViewBind()` for view binding.
- **Format**: Parameter changed from `body interface{}` to `handlers ...ResFmt`.
- **Redirect**: Use `c.Redirect().To()` instead.
- **SendFile**: Now supports different configurations using a config parameter.
- **Attachment and Download**: Non-ASCII filenames now use `filename*` as
  specified by [RFC 6266](https://www.rfc-editor.org/rfc/rfc6266) and
  [RFC 8187](https://www.rfc-editor.org/rfc/rfc8187).
- **Context**: Renamed to `RequestCtx` to correspond with the FastHTTP Request Context.
- **UserContext**: Renamed to `Context`, which returns a `context.Context` object.
- **SetUserContext**: Renamed to `SetContext`.

### SendStreamWriter

In v3, we introduced support for buffered streaming with the addition of the `SendStreamWriter` method:

```go
func (c Ctx) SendStreamWriter(streamWriter func(w *bufio.Writer))
```

With this new method, you can implement:

- Server-Side Events (SSE)
- Large file downloads
- Live data streaming

```go
app.Get("/sse", func(c fiber.Ctx) {
    c.Set("Content-Type", "text/event-stream")
    c.Set("Cache-Control", "no-cache")
    c.Set("Connection", "keep-alive")
    c.Set("Transfer-Encoding", "chunked")

    return c.SendStreamWriter(func(w *bufio.Writer) {
        for {
            fmt.Fprintf(w, "event: my-event\n")
            fmt.Fprintf(w, "data: Hello SSE\n\n")

            if err := w.Flush(); err != nil {
                log.Print("Client disconnected!")
                return
            }
        }
    })
})
```

You can find more details about this feature in [/docs/api/ctx.md](./api/ctx.md).

### Drop

In v3, we introduced support to silently terminate requests through `Drop`.

```go
func (c Ctx) Drop()
```

With this method, you can:

- Block certain requests without notifying the client to mitigate DDoS attacks
- Protect sensitive endpoints from unauthorized access without leaking errors.

:::caution
While this feature adds the ability to drop connections, it is still **highly recommended** to use additional
measures (such as **firewalls**, **proxies**, etc.) to further protect your server endpoints by blocking
malicious connections before the server establishes a connection.
:::

```go
app.Get("/", func(c fiber.Ctx) error {
    if c.IP() == "192.168.1.1" {
        return c.Drop()
    }

    return c.SendString("Hello World!")
})
```

You can find more details about this feature in [/docs/api/ctx.md](./api/ctx.md).

### End

In v3, we introduced a new method to match the Express.js API's `res.end()` method.

```go
func (c Ctx) End()
```

With this method, you can:

- Stop middleware from controlling the connection after a handler further up the method chain
  by immediately flushing the current response and closing the connection.
- Use `return c.End()` as an alternative to `return nil`

```go
app.Use(func (c fiber.Ctx) error {
    err := c.Next()
    if err != nil {
        log.Println("Got error: %v", err)
        return c.SendString(err.Error()) // Will be unsuccessful since the response ended below
    }
    return nil
})

app.Get("/hello", func (c fiber.Ctx) error {
    query := c.Query("name", "")
    if query == "" {
        c.SendString("You don't have a name?")
        c.End() // Closes the underlying connection
        return errors.New("No name provided")
    }
    return c.SendString("Hello, " + query + "!")
})
```

---

## 🌎 Client package

The Gofiber client has been completely rebuilt. It includes numerous new features such as Cookiejar, request/response hooks, and more.
You can take a look to [client docs](./client/rest.md) to see what's new with the client.

## 📎 Binding

Fiber v3 introduces a new binding mechanism that simplifies the process of binding request data to structs. The new binding system supports binding from various sources such as URL parameters, query parameters, headers, and request bodies. This unified approach makes it easier to handle different types of request data in a consistent manner.

### New Features

- Unified binding from URL parameters, query parameters, headers, and request bodies.
- Support for custom binders and constraints.
- Improved error handling and validation.
- Support multipart file binding for `*multipart.FileHeader`, `*[]*multipart.FileHeader`, and `[]*multipart.FileHeader` field types.
- Support for unified binding (`Bind().All()`) with defined precedence order: (URI -> Body -> Query -> Headers -> Cookies). [Learn more](./api/bind.md#all).
- Support MsgPack binding for request body.

<details>
<summary>Example</summary>

```go
type User struct {
    ID    int    `params:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

app.Post("/user/:id", func(c fiber.Ctx) error {
    var user User
    if err := c.Bind().Body(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(user)
})
```

In this example, the `Bind` method is used to bind the request body to the `User` struct. The `Body` method of the `Bind` class performs the actual binding.

</details>

## 🔄 Redirect

Fiber v3 enhances the redirect functionality by introducing new methods and improving existing ones. The new redirect methods provide more flexibility and control over the redirection process.

### New Methods

- `Redirect().To()`: Redirects to a specific URL.
- `Redirect().Route()`: Redirects to a named route.
- `Redirect().Back()`: Redirects to the previous URL.

<details>
<summary>Example</summary>

```go
app.Get("/old", func(c fiber.Ctx) error {
    return c.Redirect().To("/new")
})

app.Get("/new", func(c fiber.Ctx) error {
    return c.SendString("Welcome to the new route!")
})
```

</details>

### Changed behavior

:::info

The default redirect status code has been updated from `302 Found` to `303 See Other` to ensure more consistent behavior across different browsers.

:::

## 🧰 Generic functions

Fiber v3 introduces new generic functions that provide additional utility and flexibility for developers. These functions are designed to simplify common tasks and improve code readability.

### New Generic Functions

- **Convert**: Converts a value with a specified converter function and default value.
- **Locals**: Retrieves or sets local values within a request context.
- **Params**: Retrieves route parameters and can handle various types of route parameters.
- **Query**: Retrieves the value of a query parameter from the request URI and can handle various types of query parameters.
- **GetReqHeader**: Returns the HTTP request header specified by the field and can handle various types of header values.

### Example

<details>
<summary>Convert</summary>

```go
package main

import (
    "strconv"
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Get("/convert", func(c fiber.Ctx) error {
        value, err := fiber.Convert[int](c.Query("value"), strconv.Atoi, 0)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).SendString(err.Error())
        }
        return c.JSON(value)
    })

    app.Listen(":3000")
}
```

```sh
curl "http://localhost:3000/convert?value=123"
# Output: 123

curl "http://localhost:3000/convert?value=abc"
# Output: "failed to convert: strconv.Atoi: parsing \"abc\": invalid syntax"
```

</details>

<details>
<summary>Locals</summary>

```go
package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Use("/user/:id", func(c fiber.Ctx) error {
        // ask database for user
        // ...
        // set local values from database
        fiber.Locals[string](c, "user", "john")
        fiber.Locals[int](c, "age", 25)
        // ...

        return c.Next()
    })

    app.Get("/user/*", func(c fiber.Ctx) error {
        // get local values
        name := fiber.Locals[string](c, "user")
        age := fiber.Locals[int](c, "age")
        // ...
        return c.JSON(fiber.Map{"name": name, "age": age})
    })

    app.Listen(":3000")
}
```

```sh
curl "http://localhost:3000/user/5"
# Output: {"name":"john","age":25}
```

</details>

<details>
<summary>Params</summary>

```go
package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Get("/params/:id", func(c fiber.Ctx) error {
        id := fiber.Params[int](c, "id", 0)
        return c.JSON(id)
    })

    app.Listen(":3000")
}
```

```sh
curl "http://localhost:3000/params/123"
# Output: 123

curl "http://localhost:3000/params/abc"
# Output: 0
```

</details>

<details>
<summary>Query</summary>

```go
package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Get("/query", func(c fiber.Ctx) error {
        age := fiber.Query[int](c, "age", 0)
        return c.JSON(age)
    })

    app.Listen(":3000")
}

```

```sh
curl "http://localhost:3000/query?age=25"
# Output: 25

curl "http://localhost:3000/query?age=abc"
# Output: 0
```

</details>

<details>
<summary>GetReqHeader</summary>

```go
package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    app.Get("/header", func(c fiber.Ctx) error {
        userAgent := fiber.GetReqHeader[string](c, "User-Agent", "Unknown")
        return c.JSON(userAgent)
    })

    app.Listen(":3000")
}
```

```sh
curl -H "User-Agent: CustomAgent" "http://localhost:3000/header"
# Output: "CustomAgent"

curl "http://localhost:3000/header"
# Output: "Unknown"
```

</details>

## 🥡 Services

Fiber v3 introduces a new feature called Services. This feature allows developers to quickly start services that the application depends on, removing the need to manually provision things like database servers, caches, or message brokers, to name a few.

### Example

<details>
<summary>Adding a service</summary>

```go
package main

import (
    "strconv"
    "github.com/gofiber/fiber/v3"
)

type myService struct {
    img string
    // ...
}

// Start initializes and starts the service. It implements the [fiber.Service] interface.
func (s *myService) Start(ctx context.Context) error {
    // start the service
    return nil
}

// String returns a string representation of the service.
// It is used to print a human-readable name of the service in the startup message.
// It implements the [fiber.Service] interface.
func (s *myService) String() string {
    return s.img
}

// State returns the current state of the service.
// It implements the [fiber.Service] interface.
func (s *myService) State(ctx context.Context) (string, error) {
    return "running", nil
}

// Terminate stops and removes the service. It implements the [fiber.Service] interface.
func (s *myService) Terminate(ctx context.Context) error {
    // stop the service
    return nil
}

func main() {
    cfg := &fiber.Config{}

    cfg.Services = append(cfg.Services, &myService{img: "postgres:latest"})
    cfg.Services = append(cfg.Services, &myService{img: "redis:latest"})

    app := fiber.New(*cfg)

    // ...
}
```

</details>

<details>
<summary>Output</summary>

```sh
$ go run . -v

    _______ __             
   / ____(_) /_  ___  _____
  / /_  / / __ \/ _ \/ ___/
 / __/ / / /_/ /  __/ /    
/_/   /_/_.___/\___/_/          v3.0.0
--------------------------------------------------
INFO Server started on:         http://127.0.0.1:3000 (bound on host 0.0.0.0 and port 3000)
INFO Services:     2
INFO   🥡 [ RUNNING ] postgres:latest
INFO   🥡 [ RUNNING ] redis:latest
INFO Total handlers count:      2
INFO Prefork:                   Disabled
INFO PID:                       12279
INFO Total process count:       1
```

</details>

## 📃 Log

`fiber.AllLogger` interface now has a new method called `Logger`. This method can be used to get the underlying logger instance from the Fiber logger middleware. This is useful when you want to configure the logger middleware with a custom logger and still want to access the underlying logger instance.

You can find more details about this feature in [/docs/api/log.md](./api/log.md#logger).

`logger.Config` now supports a new field called `ForceColors`. This field allows you to force the logger to always use colors, even if the output is not a terminal. This is useful when you want to ensure that the logs are always colored, regardless of the output destination.

```go
package main

import "github.com/gofiber/fiber/v3/middleware/logger"

app.Use(logger.New(logger.Config{
    ForceColors: true,
}))
```

## 📦 Storage Interface

The storage interface has been updated to include new subset of methods with `WithContext` suffix. These methods allow you to pass a context to the storage operations, enabling better control over timeouts and cancellation if needed. This is particularly useful when storage implementations used outside of the Fiber core, such as in background jobs or long-running tasks.

**New Methods Signatures:**

```go
// GetWithContext gets the value for the given key with a context.
// `nil, nil` is returned when the key does not exist
GetWithContext(ctx context.Context, key string) ([]byte, error)

// SetWithContext stores the given value for the given key
// with an expiration value, 0 means no expiration.
// Empty key or value will be ignored without an error.
SetWithContext(ctx context.Context, key string, val []byte, exp time.Duration) error

// DeleteWithContext deletes the value for the given key with a context.
// It returns no error if the storage does not contain the key,
DeleteWithContext(ctx context.Context, key string) error

// ResetWithContext resets the storage and deletes all keys with a context.
ResetWithContext(ctx context.Context) error
```

## 🧬 Middlewares

### Important Change for Accessing Middleware Data

In Fiber v3, many middlewares that previously set values in `c.Locals()` using string keys (e.g., `c.Locals("requestid")`) have been updated. To align with Go's context best practices and prevent key collisions, these middlewares now store their specific data in the request's context using unexported keys of custom types.

This means that directly accessing these values via `c.Locals("some_string_key")` will no longer work for such middleware-provided data.

**How to Access Middleware Data in v3:**

Each affected middleware now provides dedicated exported functions to retrieve its specific data from the context. You should use these functions instead of relying on string-based lookups in `c.Locals()`.

Examples include:

- `requestid.FromContext(c)`
- `csrf.TokenFromContext(c)`
- `csrf.HandlerFromContext(c)`
- `session.FromContext(c)`
- `basicauth.UsernameFromContext(c)`
- `keyauth.TokenFromContext(c)`

When used with the Logger middleware, the recommended approach is to use the `CustomTags` feature of the logger, which allows you to call these specific `FromContext` functions. See the [Logger](#logger) section for more details.

### Adaptor

The adaptor middleware has been significantly optimized for performance and efficiency. Key improvements include reduced response times, lower memory usage, and fewer memory allocations. These changes make the middleware more reliable and capable of handling higher loads effectively. Enhancements include the introduction of a `sync.Pool` for managing `fasthttp.RequestCtx` instances and better HTTP request and response handling between net/http and fasthttp contexts.

| Payload Size | Metric         | V2           | V3          | Percent Change |
| ------------ | -------------- | ------------ | ----------- | -------------- |
| 100KB        | Execution Time | 1056 ns/op   | 588.6 ns/op | -44.25%        |
|              | Memory Usage   | 2644 B/op    | 254 B/op    | -90.39%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 500KB        | Execution Time | 1061 ns/op   | 562.9 ns/op | -46.94%        |
|              | Memory Usage   | 2644 B/op    | 248 B/op    | -90.62%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 1MB          | Execution Time | 1080 ns/op   | 629.7 ns/op | -41.68%        |
|              | Memory Usage   | 2646 B/op    | 267 B/op    | -89.91%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 5MB          | Execution Time | 1093 ns/op   | 540.3 ns/op | -50.58%        |
|              | Memory Usage   | 2654 B/op    | 254 B/op    | -90.43%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 10MB         | Execution Time | 1044 ns/op   | 533.1 ns/op | -48.94%        |
|              | Memory Usage   | 2665 B/op    | 258 B/op    | -90.32%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 25MB         | Execution Time | 1069 ns/op   | 540.7 ns/op | -49.42%        |
|              | Memory Usage   | 2706 B/op    | 289 B/op    | -89.32%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |
| 50MB         | Execution Time | 1137 ns/op   | 554.6 ns/op | -51.21%        |
|              | Memory Usage   | 2734 B/op    | 298 B/op    | -89.10%        |
|              | Allocations    | 16 allocs/op | 5 allocs/op | -68.75%        |

### BasicAuth

The BasicAuth middleware now validates the `Authorization` header more rigorously and sets security-focused response headers. Passwords must be provided in **hashed** form (e.g. SHA-256 or bcrypt) rather than plaintext. The default challenge includes the `charset="UTF-8"` parameter and disables caching. Responses also set a `Vary: Authorization` header to prevent caching based on credentials. Passwords are no longer stored in the request context. A `Charset` option controls the value used in the challenge header.
A new `HeaderLimit` option restricts the maximum length of the `Authorization` header (default: `8192` bytes).
The `Authorizer` function now receives the current `fiber.Ctx` as a third argument, allowing credential checks to incorporate request context.

### Cache

We are excited to introduce a new option in our caching middleware: Cache Invalidator. This feature provides greater control over cache management, allowing you to define a custom conditions for invalidating cache entries.
Additionally, the caching middleware has been optimized to avoid caching non-cacheable status codes, as defined by the [HTTP standards](https://datatracker.ietf.org/doc/html/rfc7231#section-6.1). This improvement enhances cache accuracy and reduces unnecessary cache storage usage.
Cached responses now include an RFC-compliant Age header, providing a standardized indication of how long a response has been stored in cache since it was originally generated. This enhancement improves HTTP compliance and facilitates better client-side caching strategies.

:::note
The deprecated `Store` and `Key` options have been removed in v3. Use `Storage` and `KeyGenerator` instead.
:::

### CORS

We've made some changes to the CORS middleware to improve its functionality and flexibility. Here's what's new:

#### New Struct Fields

- `Config.AllowPrivateNetwork`: This new field is a boolean that allows you to control whether private networks are allowed. This is related to the [Private Network Access (PNA)](https://wicg.github.io/private-network-access/) specification from the Web Incubator Community Group (WICG). When set to `true`, the CORS middleware will allow CORS preflight requests from private networks and respond with the `Access-Control-Allow-Private-Network: true` header. This could be useful in development environments or specific use cases, but should be done with caution due to potential security risks.

#### Updated Struct Fields

We've updated several fields from a single string (containing comma-separated values) to slices, allowing for more explicit declaration of multiple values. Here are the updated fields:

- `Config.AllowOrigins`: Now accepts a slice of strings, each representing an allowed origin.
- `Config.AllowMethods`: Now accepts a slice of strings, each representing an allowed method.
- `Config.AllowHeaders`: Now accepts a slice of strings, each representing an allowed header.
- `Config.ExposeHeaders`: Now accepts a slice of strings, each representing an exposed header.

### Compression

We've added support for `zstd` compression on top of `gzip`, `deflate`, and `brotli`.

### CSRF

The `Expiration` field in the CSRF middleware configuration has been renamed to `IdleTimeout` to better describe its functionality. Additionally, the default value has been reduced from 1 hour to 30 minutes.

### EncryptCookie

Added support for specifying Key length when using `encryptcookie.GenerateKey(length)`. This allows the user to generate keys compatible with `AES-128`, `AES-192`, and `AES-256` (Default).

### EnvVar

The `ExcludeVars` field has been removed from the EnvVar middleware configuration. When upgrading, remove any references to this field and explicitly list the variables you wish to expose using `ExportVars`.

### Filesystem

We've decided to remove filesystem middleware to clear up the confusion between static and filesystem middleware.
Now, static middleware can do everything that filesystem middleware and static do. You can check out [static middleware](./middleware/static.md) or [migration guide](#-migration-guide) to see what has been changed.

### Healthcheck

The Healthcheck middleware has been enhanced to support more than two routes, with default endpoints for liveliness, readiness, and startup checks. Here's a detailed breakdown of the changes and how to use the new features.

1. **Support for More Than Two Routes**:
    - The updated middleware now supports multiple routes beyond the default liveliness and readiness endpoints. This allows for more granular health checks, such as startup probes.

2. **Default Endpoints**:
    - Three default endpoints are now available:
        - **Liveness**: `/livez`
        - **Readiness**: `/readyz`
        - **Startup**: `/startupz`
    - These endpoints can be customized or replaced with user-defined routes.

3. **Simplified Configuration**:
    - The configuration for each health check endpoint has been simplified. Each endpoint can be configured separately, allowing for more flexibility and readability.

Refer to the [healthcheck middleware migration guide](./middleware/healthcheck.md) or the [general migration guide](#-migration-guide) to review the changes.

### KeyAuth

The keyauth middleware was updated to introduce a configurable `Realm` field for the `WWW-Authenticate` header.

### Logger

New helper function called `LoggerToWriter` has been added to the logger middleware. This function allows you to use 3rd party loggers such as `logrus` or `zap` with the Fiber logger middleware without any extra afford. For example, you can use `zap` with Fiber logger middleware like this:

<details>
<summary>Example</summary>

```go
package main

import (
    "github.com/gofiber/contrib/fiberzap/v2"
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/log"
    "github.com/gofiber/fiber/v3/middleware/logger"
)

func main() {
    // Create a new Fiber instance
    app := fiber.New()

    // Create a new zap logger which is compatible with Fiber AllLogger interface
    zap := fiberzap.NewLogger(fiberzap.LoggerConfig{
        ExtraKeys: []string{"request_id"},
    })

    // Use the logger middleware with zerolog logger
    app.Use(logger.New(logger.Config{
        Output: logger.LoggerToWriter(zap, log.LevelDebug),
    }))

    // Define a route
    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Start server on http://localhost:3000
    app.Listen(":3000")
}
```

</details>

:::note
The deprecated `TagHeader` constant was removed. Use `TagReqHeader` when you need to log request headers.
:::

#### Logging Middleware Values (e.g., Request ID)

In Fiber v3, middleware (like `requestid`) now stores values in the request context using unexported keys of custom types. This aligns with Go's context best practices to prevent key collisions between packages.

As a result, directly accessing these values using string keys with `c.Locals("your_key")` or in the logger format string with `${locals:your_key}` (e.g., `${locals:requestid}`) will no longer work for values set by such middleware.

**Recommended Solution: `CustomTags`**

The cleanest and most maintainable way to include these middleware-specific values in your logs is by using the `CustomTags` option in the logger middleware configuration. This allows you to define a custom function to retrieve the value correctly from the context.

<details>
<summary>Example: Logging Request ID with CustomTags</summary>

```go
package main

import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/logger"
    "github.com/gofiber/fiber/v3/middleware/requestid"
)

func main() {
    app := fiber.New()

    // Ensure requestid middleware is used before the logger
    app.Use(requestid.New())

    app.Use(logger.New(logger.Config{
        CustomTags: map[string]logger.LogFunc{
            "requestid": func(output logger.Buffer, c fiber.Ctx, data *logger.Data, extraParam string) (int, error) {
                // Retrieve the request ID using the middleware's specific function
                return output.WriteString(requestid.FromContext(c))
            },
        },
        // Use the custom tag in your format string
        Format: "[${time}] ${ip} - ${requestid} - ${status} ${method} ${path}\n",
    }))

    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

</details>

**Alternative: Manually Copying to `Locals`**

If you have existing logging patterns that rely on `c.Locals` or prefer to manage these values in `Locals` for other reasons, you can manually copy the value from the context to `c.Locals` in a preceding middleware:

<details>
<summary>Example: Manually setting requestid in Locals</summary>

```go
app.Use(requestid.New()) // Request ID middleware
app.Use(func(c fiber.Ctx) error {
    // Manually copy the request ID to Locals
    c.Locals("requestid", requestid.FromContext(c))
    return c.Next()
})
app.Use(logger.New(logger.Config{
    // Now ${locals:requestid} can be used, but CustomTags is generally preferred
    Format: "[${time}] ${ip} - ${locals:requestid} - ${status} ${method} ${path}\n",
}))
```

</details>

Both approaches ensure your logger can access these values while respecting Go's context practices.

The `Skip` is a function to determine if logging is skipped or written to `Stream`.

<details>
<summary>Example Usage</summary>

```go
app.Use(logger.New(logger.Config{
    Skip: func(c fiber.Ctx) bool {
        // Skip logging HTTP 200 requests
        return c.Response().StatusCode() == fiber.StatusOK
    },
}))
```

```go
app.Use(logger.New(logger.Config{
    Skip: func(c fiber.Ctx) bool {
        // Only log errors, similar to an error.log
        return c.Response().StatusCode() < 400
    },
}))
```

</details>

#### Predefined Formats

Logger provides predefined formats that you can use by name or directly by specifying the format string.
<details>

<summary>Example Usage</summary>

```go
app.Use(logger.New(logger.Config{
    Format: logger.FormatCombined,
}))
```

See more in [Logger](./middleware/logger.md#predefined-formats)
</details>

### Limiter

The limiter middleware uses a new Fixed Window Rate Limiter implementation.

:::note
Deprecated fields `Duration`, `Store`, and `Key` have been removed in v3. Use `Expiration`, `Storage`, and `KeyGenerator` instead.
:::

### Monitor

Monitor middleware is migrated to the [Contrib package](https://github.com/gofiber/contrib/tree/main/monitor) with [PR #1172](https://github.com/gofiber/contrib/pull/1172).

### Proxy

The proxy middleware has been updated to improve consistency with Go naming conventions. The `TlsConfig` field in the configuration struct has been renamed to `TLSConfig`. Additionally, the `WithTlsConfig` method has been removed; you should now configure TLS directly via the `TLSConfig` property within the `Config` struct.

### Session

The Session middleware has undergone key changes in v3 to improve functionality and flexibility. While v2 methods remain available for backward compatibility, we now recommend using the new middleware handler for session management.

#### Key Updates

### Session

The session middleware has undergone significant improvements in v3, focusing on type safety, flexibility, and better developer experience.

#### Key Changes

- **Extractor Pattern**: The string-based `KeyLookup` configuration has been replaced with a more flexible and type-safe `Extractor` function pattern.

- **New Middleware Handler**: The `New` function now returns a middleware handler instead of a `*Store`. To access the session store, use the `Store` method on the middleware, or opt for `NewStore` or `NewWithStore` for custom store integration.

- **Manual Session Release**: Session instances are no longer automatically released after being saved. To ensure proper lifecycle management, you must manually call `sess.Release()`.

- **Idle Timeout**: The `Expiration` field has been replaced with `IdleTimeout`, which handles session inactivity. If the session is idle for the specified duration, it will expire. The idle timeout is updated when the session is saved. If you are using the middleware handler, the idle timeout will be updated automatically.

- **Absolute Timeout**: The `AbsoluteTimeout` field has been added. If you need to set an absolute session timeout, you can use this field to define the duration. The session will expire after the specified duration, regardless of activity.

For more details on these changes and migration instructions, check the [Session Middleware Migration Guide](./middleware/session.md#migration-guide).

### Timeout

The timeout middleware is now configurable. A new `Config` struct allows customizing the timeout duration, defining a handler that runs when a timeout occurs, and specifying errors to treat as timeouts. The `New` function now accepts a `Config` value instead of a duration.

**Migration:** Replace calls like `timeout.New(handler, 2*time.Second)` with `timeout.New(handler, timeout.Config{Timeout: 2 * time.Second})`.

## 🔌 Addons

In v3, Fiber introduced Addons. Addons are additional useful packages that can be used in Fiber.

### Retry

The Retry addon is a new addon that implements a retry mechanism for unsuccessful network operations. It uses an exponential backoff algorithm with jitter.
It calls the function multiple times and tries to make it successful. If all calls are failed, then, it returns an error.
It adds a jitter at each retry step because adding a jitter is a way to break synchronization across the client and avoid collision.

<details>
<summary>Example</summary>

```go
package main

import (
    "fmt"

    "github.com/gofiber/fiber/v3/addon/retry"
    "github.com/gofiber/fiber/v3/client"
)

func main() {
    expBackoff := retry.NewExponentialBackoff(retry.Config{})

    // Local variables that will be used inside of Retry
    var resp *client.Response
    var err error

    // Retry a network request and return an error to signify to try again
    err = expBackoff.Retry(func() error {
        client := client.New()
        resp, err = client.Get("https://gofiber.io")
        if err != nil {
            return fmt.Errorf("GET gofiber.io failed: %w", err)
        }
        if resp.StatusCode() != 200 {
            return fmt.Errorf("GET gofiber.io did not return OK 200")
        }
        return nil
    })

    // If all retries failed, panic
    if err != nil {
        panic(err)
    }
    fmt.Printf("GET gofiber.io succeeded with status code %d\n", resp.StatusCode())
}
```

</details>

## 📋 Migration guide

- [🚀 App](#-app-1)
- [🎣 Hooks](#-hooks-1)
- [🚀 Listen](#-listen-1)
- [🗺 Router](#-router-1)
- [🧠 Context](#-context-1)
- [📎 Binding (was Parser)](#-parser)
- [🔄 Redirect](#-redirect-1)
- [🌎 Client package](#-client-package-1)
- [🧬 Middlewares](#-middlewares-1)
  - [Important Change for Accessing Middleware Data](#important-change-for-accessing-middleware-data)
  - [BasicAuth](#basicauth-1)
  - [Cache](#cache-1)
  - [CORS](#cors-1)
  - [CSRF](#csrf-1)
  - [Filesystem](#filesystem-1)
  - [EnvVar](#envvar-1)
  - [Healthcheck](#healthcheck-1)
  - [Monitor](#monitor-1)
  - [Proxy](#proxy-1)
  - [Session](#session-1)

### 🚀 App

#### Static

Since we've removed `app.Static()`, you need to move methods to static middleware like the example below:

```go
// Before
app.Static("/", "./public")
app.Static("/prefix", "./public")
app.Static("/prefix", "./public", Static{
    Index: "index.htm",
})
app.Static("*", "./public/index.html")
```

```go
// After
app.Get("/*", static.New("./public"))
app.Get("/prefix*", static.New("./public"))
app.Get("/prefix*", static.New("./public", static.Config{
    IndexNames: []string{"index.htm", "index.html"},
}))
app.Get("*", static.New("./public/index.html"))
```

:::caution
You have to put `*` to the end of the route if you don't define static route with `app.Use`.
:::

#### Trusted Proxies

We've renamed `EnableTrustedProxyCheck` to `TrustProxy` and moved `TrustedProxies` to `TrustProxyConfig`.

```go
// Before
app := fiber.New(fiber.Config{
    // EnableTrustedProxyCheck enables the trusted proxy check.
    EnableTrustedProxyCheck: true,
    // TrustedProxies is a list of trusted proxy IP ranges/addresses.
    TrustedProxies: []string{"0.8.0.0", "127.0.0.0/8", "::1/128"},
})
```

```go
// After
app := fiber.New(fiber.Config{
    // TrustProxy enables the trusted proxy check
    TrustProxy: true,
    // TrustProxyConfig allows for configuring trusted proxies.
    TrustProxyConfig: fiber.TrustProxyConfig{
        // Proxies is a list of trusted proxy IP ranges/addresses.
        Proxies: []string{"0.8.0.0"},
        // Trust all loop-back IP addresses (127.0.0.0/8, ::1/128)
        Loopback: true,
    }
})
```

### 🎣 Hooks

`OnShutdown` has been replaced by two hooks: `OnPreShutdown` and `OnPostShutdown`.
Use them to run cleanup code before and after the server shuts down. When handling
shutdown errors, register an `OnPostShutdown` hook and call `app.Listen()` in a goroutine.

```go
// Before
app.OnShutdown(func() {
    // Code to run before shutdown
})
```

```go
// After
app.OnPreShutdown(func() {
    // Code to run before shutdown
})
```

### 🚀 Listen

The `Listen` helpers (`ListenTLS`, `ListenMutualTLS`, etc.) were removed. Use
`app.Listen()` with `fiber.ListenConfig` and a `tls.Config` when TLS is required.
Options such as `ListenerNetwork` and `UnixSocketFileMode` are now configured via
this struct.

```go
// Before
app.ListenTLS(":3000", "cert.pem", "key.pem")
```

```go
// After
app.Listen(":3000", fiber.ListenConfig{
    CertFile: "./cert.pem",
    CertKeyFile: "./cert.key",
})
```

### 🗺 Router

#### Middleware Registration

The signatures for [`Add`](#middleware-registration) and [`Route`](#route-chaining) have been changed.

To migrate [`Add`](#middleware-registration) you must change the `methods` in a slice.

```go
// Before
app.Add(fiber.MethodPost, "/api", myHandler)
```

```go
// After
app.Add([]string{fiber.MethodPost}, "/api", myHandler)
```

#### Mounting

In Fiber v3, the `Mount` method has been removed. Instead, you can use the `Use` method to achieve similar functionality.

```go
// Before
app.Mount("/api", apiApp)
```

```go
// After
app.Use("/api", apiApp)
```

#### Route Chaining

Refer to the [route chaining](#route-chaining) section for details on migrating `Route`.

```go
// Before
app.Route("/api", func(apiGrp Router) {
    apiGrp.Route("/user/:id?", func(userGrp Router) {
        userGrp.Get("/", func(c fiber.Ctx) error {
            // Get user
            return c.JSON(fiber.Map{"message": "Get user", "id": c.Params("id")})
        })
        userGrp.Post("/", func(c fiber.Ctx) error {
            // Create user
            return c.JSON(fiber.Map{"message": "User created"})
        })
    })
})
```

```go
// After
app.Route("/api").Route("/user/:id?")
    .Get(func(c fiber.Ctx) error {
        // Get user
        return c.JSON(fiber.Map{"message": "Get user", "id": c.Params("id")})
    })
    .Post(func(c fiber.Ctx) error {
        // Create user
        return c.JSON(fiber.Map{"message": "User created"})
    });
```

### 🗺 RebuildTree

We introduced a new method that enables rebuilding the route tree stack at runtime. This allows you to add routes dynamically while your application is running and update the route tree to make the new routes available for use.

For more details, refer to the [app documentation](./api/app.md#rebuildtree):

#### Example Usage

```go
app.Get("/define", func(c Ctx) error {  // Define a new route dynamically
    app.Get("/dynamically-defined", func(c Ctx) error {  // Adding a dynamically defined route
        return c.SendStatus(http.StatusOK)
    })

    app.RebuildTree()  // Rebuild the route tree to register the new route

    return c.SendStatus(http.StatusOK)
})
```

In this example, a new route is defined, and `RebuildTree()` is called to ensure the new route is registered and available.

Note: Use this method with caution. It is **not** thread-safe and can be very performance-intensive. Therefore, it should be used sparingly and primarily in development mode. It should not be invoke concurrently.

## RemoveRoute

- **RemoveRoute**: Removes route by path

- **RemoveRouteByName**: Removes route by name

- **RemoveRouteFunc**: Removes route by a function having `*Route` parameter

For more details, refer to the [app documentation](./api/app.md#removeroute):

### 🧠 Context

Fiber v3 introduces several new features and changes to the Ctx interface, enhancing its functionality and flexibility.

- **ParamsInt**: Use `Params` with generic types.
- **QueryBool**: Use `Query` with generic types.
- **QueryFloat**: Use `Query` with generic types.
- **QueryInt**: Use `Query` with generic types.
- **Bind**: Now used for binding instead of view binding. Use `c.ViewBind()` for view binding.

In Fiber v3, the `Ctx` parameter in handlers is now an interface, which means the `*` symbol is no longer used. Here is an example demonstrating this change:

<details>
<summary>Example</summary>

**Before**:

```go
package main

import (
    "github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    // Route Handler with *fiber.Ctx
    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

**After**:

```go
package main

import (
    "github.com/gofiber/fiber/v3"
)

func main() {
    app := fiber.New()

    // Route Handler without *fiber.Ctx
    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    app.Listen(":3000")
}
```

**Explanation**:

In this example, the `Ctx` parameter in the handler is used as an interface (`fiber.Ctx`) instead of a pointer (`*fiber.Ctx`). This change allows for more flexibility and customization in Fiber v3.

</details>

#### 📎 Parser

The `Parser` section in Fiber v3 has undergone significant changes to improve functionality and flexibility.

##### Migration Instructions

1. **BodyParser**: Use `c.Bind().Body()` instead of `c.BodyParser()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Post("/user", func(c *fiber.Ctx) error {
        var user User
        if err := c.BodyParser(&user); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(user)
    })
    ```

    ```go
    // After
    app.Post("/user", func(c fiber.Ctx) error {
        var user User
        if err := c.Bind().Body(&user); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(user)
    })
    ```

    </details>

2. **ParamsParser**: Use `c.Bind().URI()` instead of `c.ParamsParser()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/user/:id", func(c *fiber.Ctx) error {
        var params Params
        if err := c.ParamsParser(&params); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(params)
    })
    ```

    ```go
    // After
    app.Get("/user/:id", func(c fiber.Ctx) error {
        var params Params
        if err := c.Bind().URI(&params); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(params)
    })
    ```

    </details>

3. **QueryParser**: Use `c.Bind().Query()` instead of `c.QueryParser()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/search", func(c *fiber.Ctx) error {
        var query Query
        if err := c.QueryParser(&query); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(query)
    })
    ```

    ```go
    // After
    app.Get("/search", func(c fiber.Ctx) error {
        var query Query
        if err := c.Bind().Query(&query); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(query)
    })
    ```

    </details>

4. **CookieParser**: Use `c.Bind().Cookie()` instead of `c.CookieParser()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/cookie", func(c *fiber.Ctx) error {
        var cookie Cookie
        if err := c.CookieParser(&cookie); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(cookie)
    })
    ```

    ```go
    // After
    app.Get("/cookie", func(c fiber.Ctx) error {
        var cookie Cookie
        if err := c.Bind().Cookie(&cookie); err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
        }
        return c.JSON(cookie)
    })
    ```

    </details>

#### 🔄 Redirect

Fiber v3 enhances the redirect functionality by introducing new methods and improving existing ones. The new redirect methods provide more flexibility and control over the redirection process.

##### Migration Instructions

1. **RedirectToRoute**: Use `c.Redirect().Route()` instead of `c.RedirectToRoute()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/old", func(c *fiber.Ctx) error {
        return c.RedirectToRoute("newRoute")
    })
    ```

    ```go
    // After
    app.Get("/old", func(c fiber.Ctx) error {
        return c.Redirect().Route("newRoute")
    })
    ```

    </details>

2. **RedirectBack**: Use `c.Redirect().Back()` instead of `c.RedirectBack()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/back", func(c *fiber.Ctx) error {
        return c.RedirectBack()
    })
    ```

    ```go
    // After
    app.Get("/back", func(c fiber.Ctx) error {
        return c.Redirect().Back()
    })
    ```

    </details>

3. **Redirect**: Use `c.Redirect().To()` instead of `c.Redirect()`.

    <details>
    <summary>Example</summary>

    ```go
    // Before
    app.Get("/old", func(c *fiber.Ctx) error {
        return c.Redirect("/new")
    })
    ```

    ```go
    // After
    app.Get("/old", func(c fiber.Ctx) error {
        return c.Redirect().To("/new")
    })
    ```

    </details>

### 🌎 Client package

Fiber v3 introduces a completely rebuilt client package with numerous new features such as Cookiejar, request/response hooks, and more. Here is a guide to help you migrate from Fiber v2 to Fiber v3.

#### New Features

- **Cookiejar**: Manage cookies automatically.
- **Request/Response Hooks**: Customize request and response handling.
- **Improved Error Handling**: Better error management and reporting.

#### Migration Instructions

**Import Path**:

Update the import path to the new client package.

<details>
<summary>Before</summary>

```go
import "github.com/gofiber/fiber/v2/client"
```

</details>

<details>
<summary>After</summary>

```go
import "github.com/gofiber/fiber/v3/client"
```

</details>

### 🧬 Middlewares

#### Important Change for Accessing Middleware Data

**Change:** In Fiber v2, some middlewares set data in `c.Locals()` using string keys (e.g., `c.Locals("requestid")`). In Fiber v3, to align with Go's context best practices and prevent key collisions, these middlewares now store their specific data in the request's context using unexported keys of custom types.

**Impact:** Directly accessing these middleware-provided values via `c.Locals("some_string_key")` will no longer work.

**Migration Action:**
You must update your code to use the dedicated exported functions provided by each affected middleware to retrieve its data from the context.

**Examples of new helper functions to use:**

- `requestid.FromContext(c)`
- `csrf.TokenFromContext(c)`
- `csrf.HandlerFromContext(c)`
- `session.FromContext(c)`
- `basicauth.UsernameFromContext(c)`
- `keyauth.TokenFromContext(c)`

**For logging these values:**
The recommended approach is to use the `CustomTags` feature of the Logger middleware, which allows you to call these specific `FromContext` functions. Refer to the [Logger section in "What's New"](#logger) for detailed examples.

:::note
If you were manually setting and retrieving your own application-specific values in `c.Locals()` using string keys, that functionality remains unchanged. This change specifically pertains to how Fiber's built-in (and some contrib) middlewares expose their data.
:::

#### BasicAuth

The `Authorizer` callback now receives the current request context. Update custom
functions from:

```go
Authorizer: func(user, pass string) bool {
    // v2 style
    return user == "admin" && pass == "secret"
}
```

to:

```go
Authorizer: func(user, pass string, _ fiber.Ctx) bool {
    // v3 style with access to the Fiber context
    return user == "admin" && pass == "secret"
}
```

Passwords configured for BasicAuth must now be pre-hashed. If no prefix is supplied the middleware expects a SHA-256 digest encoded in hex. Common prefixes like `{SHA256}` and `{SHA512}` and bcrypt strings are also supported. Plaintext passwords are no longer accepted. Unauthorized responses also include a `Vary: Authorization` header for correct caching behavior.

You can also set the optional `HeaderLimit` and `Charset`
options to further control authentication behavior.

#### Cache

The deprecated `Store` and `Key` fields were removed. Use `Storage` and
`KeyGenerator` instead to configure caching backends and cache keys.

#### CORS

The CORS middleware has been updated to use slices instead of strings for the `AllowOrigins`, `AllowMethods`, `AllowHeaders`, and `ExposeHeaders` fields. Here's how you can update your code:

```go
// Before
app.Use(cors.New(cors.Config{
    AllowOrigins: "https://example.com,https://example2.com",
    AllowMethods: strings.Join([]string{fiber.MethodGet, fiber.MethodPost}, ","),
    AllowHeaders: "Content-Type",
    ExposeHeaders: "Content-Length",
}))

// After
app.Use(cors.New(cors.Config{
    AllowOrigins: []string{"https://example.com", "https://example2.com"},
    AllowMethods: []string{fiber.MethodGet, fiber.MethodPost},
    AllowHeaders: []string{"Content-Type"},
    ExposeHeaders: []string{"Content-Length"},
}))
```

#### CSRF

- **Field Renaming**: The `Expiration` field in the CSRF middleware configuration has been renamed to `IdleTimeout` to better describe its functionality. Additionally, the default value has been reduced from 1 hour to 30 minutes. Update your code as follows:

```go
// Before
app.Use(csrf.New(csrf.Config{
    Expiration: 10 * time.Minute,
}))

// After
app.Use(csrf.New(csrf.Config{
    IdleTimeout: 10 * time.Minute,
}))
```

- **Session Key Removal**: The `SessionKey` field has been removed from the CSRF middleware configuration. The session key is now an unexported constant within the middleware to avoid potential key collisions in the session store.

- **KeyLookup Field Removal**: The `KeyLookup` field has been removed from the CSRF middleware configuration. This field was deprecated and is no longer needed as the middleware now uses a more secure approach for token management.

```go
// Before
app.Use(csrf.New(csrf.Config{
    KeyLookup: "header:X-Csrf-Token",
    // other config...
}))

// After - use Extractor instead
app.Use(csrf.New(csrf.Config{
    Extractor: csrf.FromHeader("X-Csrf-Token"),
    // other config...
}))
```

- **FromCookie Extractor Removal**: The `csrf.FromCookie` extractor has been intentionally removed for security reasons. Using cookie-based extraction defeats the purpose of CSRF protection by making the extracted token always match the cookie value.

```go
// Before - This was a security vulnerability
app.Use(csrf.New(csrf.Config{
    Extractor: csrf.FromCookie("csrf_token"), // ❌ Insecure!
}))

// After - Use secure extractors instead
app.Use(csrf.New(csrf.Config{
    Extractor: csrf.FromHeader("X-Csrf-Token"), // ✅ Secure
    // or
    Extractor: csrf.FromForm("_csrf"),          // ✅ Secure
    // or
    Extractor: csrf.FromQuery("csrf_token"),    // ✅ Acceptable
}))
```

**Security Note**: The removal of `FromCookie` prevents a common misconfiguration that would completely bypass CSRF protection. The middleware uses the Double Submit Cookie pattern, which requires the token to be submitted through a different channel than the cookie to provide meaningful protection.

#### Timeout

The timeout middleware now accepts a configuration struct instead of a duration.
Update your code as follows:

```go
// Before
app.Use(timeout.New(handler, 2*time.Second))

// After
app.Use(timeout.New(handler, timeout.Config{Timeout: 2 * time.Second}))
```

#### Filesystem

You need to move filesystem middleware to static middleware due to it has been removed from the core.

```go
// Before
app.Use(filesystem.New(filesystem.Config{
    Root: http.Dir("./assets"),
}))

app.Use(filesystem.New(filesystem.Config{
    Root:         http.Dir("./assets"),
    Browse:       true,
    Index:        "index.html",
    MaxAge:       3600,
}))
```

```go
// After
app.Use(static.New("", static.Config{
    FS: os.DirFS("./assets"),
}))

app.Use(static.New("", static.Config{
    FS:           os.DirFS("./assets"),
    Browse:       true,
    IndexNames:   []string{"index.html"},
    MaxAge:       3600,
}))
```

#### EnvVar

The `ExcludeVars` option has been removed. Remove any references to it and use
`ExportVars` to explicitly list environment variables that should be exposed.

#### Healthcheck

Previously, the Healthcheck middleware was configured with a combined setup for liveliness and readiness probes:

```go
//before
app.Use(healthcheck.New(healthcheck.Config{
    LivenessProbe: func(c fiber.Ctx) bool {
        return true
    },
    LivenessEndpoint: "/live",
    ReadinessProbe: func(c fiber.Ctx) bool {
        return serviceA.Ready() && serviceB.Ready() && ...
    },
    ReadinessEndpoint: "/ready",
}))
```

With the new version, each health check endpoint is configured separately, allowing for more flexibility:

```go
// after

// Default liveness endpoint configuration
app.Get(healthcheck.LivenessEndpoint, healthcheck.New(healthcheck.Config{
    Probe: func(c fiber.Ctx) bool {
        return true
    },
}))

// Default readiness endpoint configuration
app.Get(healthcheck.ReadinessEndpoint, healthcheck.New())

// New default startup endpoint configuration
// Default endpoint is /startupz
app.Get(healthcheck.StartupEndpoint, healthcheck.New(healthcheck.Config{
    Probe: func(c fiber.Ctx) bool {
        return serviceA.Ready() && serviceB.Ready() && ...
    },
}))

// Custom liveness endpoint configuration
app.Get("/live", healthcheck.New())
```

#### Monitor

Since v3 the Monitor middleware has been moved to the [Contrib package](https://github.com/gofiber/contrib/tree/main/monitor)

```go
// Before
import "github.com/gofiber/fiber/v2/middleware/monitor"

app.Use("/metrics", monitor.New())
```

You only need to change the import path to the contrib package.

```go
// After
import "github.com/gofiber/contrib/monitor"

app.Use("/metrics", monitor.New())
```

#### Proxy

In previous versions, TLS settings for the proxy middleware were set using the `WithTlsConfig` method. This method has been removed in favor of a more idiomatic configuration via the `TLSConfig` field in the `Config` struct.

#### Before (v2 usage)

```go
proxy.WithTlsConfig(&tls.Config{
    InsecureSkipVerify: true,
})

// Forward to url
app.Get("/gif", proxy.Forward("https://i.imgur.com/IWaBepg.gif"))
```

#### After (v3 usage)

```go
proxy.WithClient(&fasthttp.Client{
    TLSConfig: &tls.Config{InsecureSkipVerify: true},
})

// Forward to url
app.Get("/gif", proxy.Forward("https://i.imgur.com/IWaBepg.gif"))
```

#### Session

`session.New()` now returns a middleware handler. When using the store pattern,
create a store with `session.NewStore()` or call `Store()` on the middleware.
Sessions obtained from a store must be released manually via `sess.Release()`.
Additionally, replace the deprecated `KeyLookup` option with extractor
functions such as `session.FromCookie()` or `session.FromHeader()`. Multiple
extractors can be combined with `session.Chain()`.

```go
// Before
app.Use(session.New(session.Config{
    KeyLookup: "cookie:session_id",
    Store:     session.NewStore(),
}))
```

```go
// After
app.Use(session.New(session.Config{
    Extractor: session.FromCookie("session_id"),
    Store:     session.NewStore(),
}))
```

See the [Session Middleware Migration Guide](./middleware/session.md#migration-guide)
for complete details.
