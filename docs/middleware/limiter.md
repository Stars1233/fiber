---
id: limiter
---

# Limiter

Limiter middleware for [Fiber](https://github.com/gofiber/fiber) that is used to limit repeat requests to public APIs and/or endpoints such as password reset. It is also useful for API clients, web crawling, or other tasks that need to be throttled.

:::note
This middleware uses our [Storage](https://github.com/gofiber/storage) package to support various databases through a single interface. The default configuration for this middleware saves data to memory, see the examples below for other databases.
:::

:::note
This module does not share state with other processes/servers by default.
:::

## Signatures

```go
func New(config ...Config) fiber.Handler
```

## Examples

Import the middleware package that is part of the Fiber web framework

```go
import (
    "github.com/gofiber/fiber/v3"
    "github.com/gofiber/fiber/v3/middleware/limiter"
)
```

After you initiate your Fiber app, you can use the following possibilities:

```go
// Initialize default config
app.Use(limiter.New())

// Or extend your config for customization
app.Use(limiter.New(limiter.Config{
    Next: func(c fiber.Ctx) bool {
        return c.IP() == "127.0.0.1"
    },
    Max:          20,
    MaxFunc: func(c fiber.Ctx) int {
      return 20
    },
    Expiration:     30 * time.Second,
    KeyGenerator:          func(c fiber.Ctx) string {
        return c.Get("x-forwarded-for")
    },
    LimitReached: func(c fiber.Ctx) error {
        return c.SendFile("./toofast.html")
    },
    Storage: myCustomStorage{},
}))
```

## Sliding window

Instead of using the standard fixed window algorithm, you can enable the [sliding window](https://en.wikipedia.org/wiki/Sliding_window_protocol) algorithm.

A example of such configuration is:

```go
app.Use(limiter.New(limiter.Config{
    Max:            20,
    Expiration:     30 * time.Second,
    LimiterMiddleware: limiter.SlidingWindow{},
}))
```

This means that every window will consider the previous window (if there was any). The given formula for the rate is:

```text
weightOfPreviousWindow = previous window's amount request * (whenNewWindow / Expiration)
rate = weightOfPreviousWindow + current window's amount request.
```

## Dynamic limit

You can also calculate the limit dynamically using the MaxFunc parameter. It's a function that receives the request's context as a parameter and allow you to calculate a different limit for each request separately.

Example:

```go
app.Use(limiter.New(limiter.Config{
    MaxFunc:  func(c fiber.Ctx) int {
      return getUserLimit(ctx.Param("id"))
    },
    Expiration:     30 * time.Second,
}))
```

## Config

| Property               | Type                      | Description                                                                                 | Default                                  |
|:-----------------------|:--------------------------|:--------------------------------------------------------------------------------------------|:-----------------------------------------|
| Next                   | `func(fiber.Ctx) bool`   | Next defines a function to skip this middleware when returned true.                         | `nil`                                    |
| Max                    | `int`                     | Max number of recent connections during `Expiration` seconds before sending a 429 response. | 5                                        |
| MaxFunc                | `func(fiber.Ctx) int`     | A function to calculate the max number of recent connections during `Expiration` seconds before sending a 429 response. | A function which returns the cfg.Max    |
| KeyGenerator           | `func(fiber.Ctx) string` | KeyGenerator allows you to generate custom keys, by default c.IP() is used.                 | A function using c.IP() as the default   |
| Expiration             | `time.Duration`           | Expiration is the time on how long to keep records of requests in memory.                   | 1 * time.Minute                          |
| LimitReached           | `fiber.Handler`           | LimitReached is called when a request hits the limit.                                       | A function sending 429 response          |
| SkipFailedRequests     | `bool`                    | When set to true, requests with StatusCode >= 400 won't be counted.                         | false                                    |
| SkipSuccessfulRequests | `bool`                    | When set to true, requests with StatusCode < 400 won't be counted.                          | false                                    |
| DisableHeaders         | `bool`                    | When set to true, the middleware will not include the rate limit headers (`X-RateLimit-*` and `Retry-After`) in the response. | false                                    |
| Storage                | `fiber.Storage`           | Store is used to store the state of the middleware.                                         | An in-memory store for this process only |
| LimiterMiddleware      | `LimiterHandler`          | LimiterMiddleware is the struct that implements a limiter middleware.                       | A new Fixed Window Rate Limiter          |

:::note
A custom store can be used if it implements the `Storage` interface - more details and an example can be found in `store.go`.
:::

## Default Config

```go
var ConfigDefault = Config{
    Max:        5,
    MaxFunc: func(c fiber.Ctx) int {
      return 5
    },
    Expiration: 1 * time.Minute,
    KeyGenerator: func(c fiber.Ctx) string {
        return c.IP()
    },
    LimitReached: func(c fiber.Ctx) error {
        return c.SendStatus(fiber.StatusTooManyRequests)
    },
    SkipFailedRequests: false,
    SkipSuccessfulRequests: false,
    DisableHeaders:        false,
    LimiterMiddleware: FixedWindow{},
}
```

### Custom Storage/Database

You can use any storage from our [storage](https://github.com/gofiber/storage/) package.

```go
storage := sqlite3.New() // From github.com/gofiber/storage/sqlite3

app.Use(limiter.New(limiter.Config{
    Storage: storage,
}))
```
