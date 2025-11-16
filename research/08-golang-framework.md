# Golang Backend Framework Selection Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: Gin (Recommended for this project)

---

## Recommendation Summary

**Gin Web Framework** chosen because:
1. **Mature and battle-tested** - 86.6k GitHub stars, used by 290k+ projects
2. **Simple, Express-like API** - Familiar patterns from frontend frameworks
3. **Excellent performance** - ~40x faster than Martini, comparable to Echo/Fiber
4. **Built-in JSON binding** - Auto-parse request bodies with validation
5. **Easy file upload handling** - `c.FormFile()` for multipart/form-data
6. **Comprehensive middleware ecosystem** - JWT, CORS, logging, recovery built-in
7. **Extensive documentation** - Multiple languages, official Go tutorial
8. **Large community** - 485 contributors, active maintenance

**Alternative if specific needs arise**: 
- **Echo** - If you need more explicit HTTP handler control
- **Fiber** - If you need absolute maximum performance (v3 in beta)

**Avoid**: Fiber v3 for production (currently beta), Beego (older patterns)

---

## Core Requirements for Budget App Backend

Based on spike document:
1. **REST API endpoints** - CRUD for budgets, categories, transactions
2. **CSV file upload** - Up to 10 files per month, multipart/form-data
3. **JSON request/response** - All data exchange in JSON format
4. **PostgreSQL integration** - Database queries and transactions
5. **Authentication** - (Future) User accounts and session management
6. **CORS support** - Frontend will be separate React app
7. **Error handling** - Consistent error responses
8. **Middleware** - Logging, recovery, request validation

---

## Gin Web Framework - Express-inspired Go Framework

**Source**: https://github.com/gin-gonic/gin, https://gin-gonic.com/

### Key Features

- ✅ **Express-style routing** - `r.GET()`, `r.POST()`, `r.PUT()`, `r.DELETE()`
- ✅ **JSON binding & validation** - `c.ShouldBindJSON(&struct)` with struct tags
- ✅ **File upload support** - `c.FormFile("file")` for multipart handling
- ✅ **Route grouping** - `v1 := r.Group("/api/v1")` for API versioning
- ✅ **Middleware support** - `r.Use(middleware)` for CORS, auth, logging
- ✅ **Query/path parameters** - `c.Query("param")`, `c.Param(":id")`
- ✅ **Built-in JSON rendering** - `c.JSON(200, data)` with automatic marshaling
- ✅ **Error handling** - `c.AbortWithStatusJSON(400, gin.H{"error": "message"})`
- ✅ **Static file serving** - `r.Static("/static", "./public")`
- ✅ **HTTP/2 support** - TLS configuration for production
- ✅ **Testing support** - `httptest` integration for unit tests
- ✅ **Model binding** - Query strings, JSON, XML, YAML, protobuf

### Hello World Example

```go
package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default() // Includes Logger and Recovery middleware
    
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "message": "pong",
        })
    })
    
    r.Run(":8080") // Listen on 0.0.0.0:8080
}
```

### Budget API Example

```go
type CategoryRequest struct {
    Name      string  `json:"name" binding:"required"`
    Projected float64 `json:"projected" binding:"required,gte=0"`
    ParentID  *string `json:"parentId"` // Optional for hierarchy
}

func main() {
    r := gin.Default()
    
    // Middleware
    r.Use(cors.Default()) // CORS for frontend
    
    // API v1 group
    v1 := r.Group("/api/v1")
    {
        // Categories
        v1.GET("/categories", listCategories)
        v1.POST("/categories", createCategory)
        v1.PUT("/categories/:id", updateCategory)
        v1.DELETE("/categories/:id", deleteCategory)
        
        // CSV Upload
        v1.POST("/upload", uploadCSV)
        
        // Budgets
        v1.GET("/budgets/:month", getBudget)
        v1.PUT("/budgets/:month", saveBudget)
    }
    
    r.Run(":8080")
}

func createCategory(c *gin.Context) {
    var req CategoryRequest
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Database insert logic here
    category := insertCategory(req)
    
    c.JSON(http.StatusCreated, category)
}

func uploadCSV(c *gin.Context) {
    file, err := c.FormFile("csv")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }
    
    // Validate file type
    if !strings.HasSuffix(file.Filename, ".csv") {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Only CSV files allowed"})
        return
    }
    
    // Save or process file
    dst := fmt.Sprintf("uploads/%s", file.Filename)
    if err := c.SaveUploadedFile(file, dst); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"filename": file.Filename})
}
```

### Middleware Examples

```go
// Custom logging middleware
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        t := time.Now()
        c.Next() // Process request
        latency := time.Since(t)
        log.Printf("[%s] %s - %v", c.Request.Method, c.Request.URL.Path, latency)
    }
}

// Auth middleware (JWT example)
func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
            return
        }
        // Validate JWT token
        c.Next()
    }
}

// Usage
r.Use(Logger())
authorized := r.Group("/api")
authorized.Use(AuthRequired())
```

### Pros

- ✅ **Large community** - 86.6k stars, 290k+ dependents, 485 contributors
- ✅ **Production-ready** - Used by many companies, stable v1.x releases
- ✅ **Excellent docs** - Multi-language docs, Go official tutorial
- ✅ **Simple API** - Express-like, easy to learn
- ✅ **Built-in features** - JSON binding, validation, file upload
- ✅ **Middleware ecosystem** - gin-contrib has auth, CORS, sessions, etc.
- ✅ **Performance** - Fast enough for most apps (~43k req/s GitHub API benchmark)
- ✅ **Testing friendly** - Easy to write unit tests with httptest

### Cons

- ⚠️ **Not the absolute fastest** - Echo and Fiber slightly faster in benchmarks
- ⚠️ **Context usage** - Must pass `*gin.Context` everywhere (not `context.Context`)
- ⚠️ **Opinionated** - Less flexible than Echo's handler interface

---

## Echo - Minimalist, High-Performance Framework

**Source**: https://github.com/labstack/echo, https://echo.labstack.com/

### Key Features

- ✅ **Very lightweight** - Minimal core, plugin-based architecture
- ✅ **HTTP/2 and HTTP/3** - Built-in support
- ✅ **Automatic TLS** - Let's Encrypt integration
- ✅ **Middleware chaining** - Root, group, and route-level middleware
- ✅ **Data binding** - JSON, XML, form, query binding with validation
- ✅ **Template rendering** - Multiple template engine support
- ✅ **WebSocket support** - Built-in WebSocket handling
- ✅ **Standard `net/http` compatibility** - Can use standard middleware

### Hello World Example

```go
package main

import (
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
    "net/http"
)

func main() {
    e := echo.New()
    
    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    
    // Routes
    e.GET("/", func(c echo.Context) error {
        return c.String(http.StatusOK, "Hello, World!")
    })
    
    // Start server
    e.Start(":8080")
}
```

### Budget API Example

```go
type CategoryRequest struct {
    Name      string  `json:"name" validate:"required"`
    Projected float64 `json:"projected" validate:"required,gte=0"`
}

func main() {
    e := echo.New()
    e.Use(middleware.CORS())
    
    // API routes
    v1 := e.Group("/api/v1")
    v1.POST("/categories", createCategory)
    v1.POST("/upload", uploadCSV)
    
    e.Start(":8080")
}

func createCategory(c echo.Context) error {
    req := new(CategoryRequest)
    
    if err := c.Bind(req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    
    if err := c.Validate(req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }
    
    // Database insert
    category := insertCategory(*req)
    
    return c.JSON(http.StatusCreated, category)
}

func uploadCSV(c echo.Context) error {
    file, err := c.FormFile("csv")
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": "No file uploaded",
        })
    }
    
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()
    
    // Process file
    dst, err := os.Create("uploads/" + file.Filename)
    if err != nil {
        return err
    }
    defer dst.Close()
    
    if _, err = io.Copy(dst, src); err != nil {
        return err
    }
    
    return c.JSON(http.StatusOK, map[string]string{
        "filename": file.Filename,
    })
}
```

### Pros

- ✅ **Excellent performance** - 31.7k req/s in benchmarks (faster than Gin)
- ✅ **Clean API** - Error handling via `error` return
- ✅ **Flexible** - Less opinionated, more control over handler signatures
- ✅ **Standard library friendly** - Works with `net/http` middleware
- ✅ **Active development** - Regular updates and new features
- ✅ **Good docs** - Clear examples and API reference

### Cons

- ⚠️ **Smaller community** - 31.7k stars vs Gin's 86.6k
- ⚠️ **More setup required** - Validation needs external library (go-playground/validator)
- ⚠️ **Less built-in** - Fewer convenience functions than Gin

---

## Fiber - FastHTTP-based Framework (Express-like)

**Source**: https://github.com/gofiber/fiber, https://docs.gofiber.io/

### Key Features

- ✅ **Extremely fast** - Built on fasthttp (6-10x faster than net/http)
- ✅ **Express.js API** - Nearly 1:1 mapping to Express methods
- ✅ **Zero allocation router** - Memory-efficient routing
- ✅ **WebSocket support** - Built-in WebSocket handling
- ✅ **Server-sent events** - SSE support for real-time updates
- ✅ **Rich middleware** - 40+ built-in middleware packages
- ✅ **Low memory footprint** - Optimized for performance
- ✅ **Template engines** - Support for 9 template engines

### Hello World Example

```go
package main

import "github.com/gofiber/fiber/v3"

func main() {
    app := fiber.New()
    
    app.Get("/", func(c fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })
    
    app.Listen(":3000")
}
```

### Budget API Example

```go
type CategoryRequest struct {
    Name      string  `json:"name" validate:"required"`
    Projected float64 `json:"projected" validate:"required,min=0"`
}

func main() {
    app := fiber.New()
    
    // Middleware
    app.Use(cors.New())
    app.Use(logger.New())
    
    // API routes
    api := app.Group("/api/v1")
    api.Post("/categories", createCategory)
    api.Post("/upload", uploadCSV)
    
    app.Listen(":3000")
}

func createCategory(c fiber.Ctx) error {
    req := new(CategoryRequest)
    
    if err := c.Bind().JSON(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Validate
    if err := validator.New().Struct(req); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    
    // Database insert
    category := insertCategory(*req)
    
    return c.Status(201).JSON(category)
}

func uploadCSV(c fiber.Ctx) error {
    file, err := c.FormFile("csv")
    if err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "No file uploaded",
        })
    }
    
    if err := c.SaveFile(file, "uploads/"+file.Filename); err != nil {
        return c.Status(500).JSON(fiber.Map{
            "error": "Failed to save file",
        })
    }
    
    return c.JSON(fiber.Map{
        "filename": file.Filename,
    })
}
```

### Pros

- ✅ **Best performance** - 55k+ req/s in benchmarks (fastest of the three)
- ✅ **Express-like API** - Very familiar for JavaScript developers
- ✅ **Zero allocation** - Excellent memory efficiency
- ✅ **Rich features** - Many built-in middleware and features
- ✅ **Active development** - Frequent updates and improvements
- ✅ **Large ecosystem** - fiber-contrib has many plugins

### Cons

- ⚠️ **Fiber v3 in beta** - Latest version not production-ready (use v2.x)
- ⚠️ **fasthttp quirks** - Not standard `net/http`, some middleware incompatible
- ⚠️ **Context reuse** - Must copy values to use outside handler (zero-allocation tradeoff)
- ⚠️ **Breaking changes** - v2 → v3 has major API changes
- ⚠️ **Less mature** - Newer than Gin/Echo, smaller community

---

## Framework Comparison Matrix

| Feature                  | Gin                      | Echo                   | Fiber v2                 |
| ------------------------ | ------------------------ | ---------------------- | ------------------------ |
| **GitHub Stars**         | 86.6k                    | 31.7k                  | 38.2k                    |
| **Community Size**       | Largest                  | Medium                 | Medium                   |
| **Maturity**             | Very Mature (v1.x)       | Mature (v4.x)          | Mature (v2.x, v3 beta)   |
| **Performance**          | Fast (~43k req/s)        | Faster (~31k req/s)    | Fastest (~55k req/s)     |
| **Learning Curve**       | Easy                     | Easy                   | Easy                     |
| **API Style**            | Express-like             | Echo-specific          | Express-like             |
| **HTTP Library**         | net/http                 | net/http               | fasthttp                 |
| **JSON Binding**         | Built-in ✅               | Built-in ✅             | Built-in ✅               |
| **Validation**           | Struct tags ✅            | External library       | External library         |
| **File Upload**          | `c.FormFile()` ✅         | `c.FormFile()` ✅       | `c.FormFile()` ✅         |
| **Middleware**           | Extensive ecosystem      | Good ecosystem         | Very extensive           |
| **CORS Support**         | gin-contrib/cors ✅       | middleware.CORS() ✅    | fiber/cors ✅             |
| **Error Handling**       | `c.JSON()` return        | `error` return         | `error` return           |
| **Route Grouping**       | `r.Group()` ✅            | `e.Group()` ✅          | `app.Group()` ✅          |
| **Testing**              | Excellent                | Excellent              | Good                     |
| **Documentation**        | Excellent (multi-lang)   | Good                   | Excellent                |
| **Production Readiness** | ✅ Yes (v1.x stable)      | ✅ Yes (v4.x stable)    | ✅ Yes (v2.x), ⚠️ v3 beta  |
| **Bundle Size**          | ~10MB binary             | ~10MB binary           | ~10MB binary             |
| **Memory Usage**         | Low                      | Low                    | Very Low (zero-alloc)    |
| **Standard Lib Compat**  | ✅ Yes (`net/http`)       | ✅ Yes (`net/http`)     | ⚠️ No (fasthttp)          |
| **WebSocket Support**    | Via gorilla/websocket    | Built-in ✅             | Built-in ✅               |
| **HTTP/2 Support**       | ✅ Yes                    | ✅ Yes                  | ✅ Yes (v3)               |
| **Go Version Required**  | 1.24+                    | 1.13+                  | 1.24+                    |
| **Best For**             | Most projects, beginners | High perf, flexibility | Absolute max performance |

---

## Decision Criteria for Budget App

### Requirements Fit

**Gin** ✅ Excellent fit:
- Built-in JSON binding with validation (perfect for API)
- Simple file upload handling (`c.FormFile()`)
- Mature CORS middleware for React frontend
- Large community means better Stack Overflow support
- Express-like API familiar to team with React background

**Echo** ✅ Good fit:
- Slightly better performance (not critical for budget app)
- Clean error handling pattern
- Good middleware ecosystem

**Fiber** ⚠️ Acceptable but with caveats:
- Best performance (overkill for budget app scale)
- v3 in beta - would need to use v2
- fasthttp incompatibility with standard middleware

### Recommendation: **Gin**

**Why Gin for this project**:
1. **Team velocity** - Easiest to learn, most Stack Overflow answers
2. **Stability** - Mature v1.x releases, production-proven
3. **Built-in validation** - Struct tags for request validation
4. **Community size** - 86.6k stars, 290k+ projects using it
5. **Documentation** - Multi-language docs, Go official tutorial
6. **Ecosystem** - gin-contrib has everything needed (JWT, CORS, sessions)
7. **CSV upload** - Simple `c.FormFile()` API
8. **Future-proof** - Large community ensures long-term support

**Performance tradeoff**: Gin is "fast enough" for budget app scale. Echo/Fiber benchmarks are ~20-30% faster, but this won't matter until app has 10k+ concurrent users.

---

## Implementation Checklist

### Phase 1: Basic Setup
- [ ] Install Gin: `go get -u github.com/gin-gonic/gin`
- [ ] Create `main.go` with basic server
- [ ] Add CORS middleware for React frontend
- [ ] Set up route grouping (`/api/v1`)
- [ ] Add logging and recovery middleware

### Phase 2: API Endpoints
- [ ] Create budget CRUD endpoints
- [ ] Create category CRUD endpoints (with hierarchy support)
- [ ] Create transaction endpoints
- [ ] Implement CSV upload endpoint
- [ ] Add cell reference endpoints

### Phase 3: Middleware
- [ ] Request validation middleware
- [ ] Error handling middleware
- [ ] (Future) Authentication middleware (JWT)
- [ ] Rate limiting middleware
- [ ] Request logging

### Phase 4: Testing
- [ ] Write unit tests for handlers
- [ ] Integration tests with test database
- [ ] Load testing with ab/wrk

---

## Example Project Structure

```
backend/
├── main.go                 # Application entry point
├── config/                 # Configuration
│   └── config.go
├── controllers/            # HTTP handlers
│   ├── budget.go
│   ├── category.go
│   ├── transaction.go
│   └── upload.go
├── models/                 # Database models
│   ├── budget.go
│   ├── category.go
│   └── transaction.go
├── services/               # Business logic
│   ├── budget_service.go
│   └── category_service.go
├── middleware/             # Custom middleware
│   ├── auth.go
│   └── validator.go
├── db/                     # Database
│   ├── db.go               # Connection
│   └── migrations/         # SQL migrations
├── routes/                 # Route definitions
│   └── routes.go
└── tests/                  # Tests
    ├── budget_test.go
    └── integration_test.go
```

---

## Alternative Scenarios

**Use Echo if**:
- Need absolute maximum performance (millions of requests/day)
- Want more explicit control over handler signatures
- Prefer error-based control flow (`return error`)

**Use Fiber if**:
- Need zero-allocation performance (high-frequency trading level)
- Team has Express.js expertise
- Willing to use v2 and migrate to v3 later

**Stick with Gin if**:
- Building typical CRUD API
- Want stable, proven framework
- Value community size and resources

---

## References

- [Gin Documentation](https://gin-gonic.com/docs/)
- [Gin GitHub](https://github.com/gin-gonic/gin)
- [Go Official Tutorial: RESTful API with Gin](https://go.dev/doc/tutorial/web-service-gin)
- [Echo Documentation](https://echo.labstack.com/docs)
- [Fiber Documentation](https://docs.gofiber.io/)
- [Go Web Framework Benchmarks](https://github.com/gin-gonic/gin/blob/master/BENCHMARKS.md)

---

**Next Research**: API design patterns, database integration (GORM vs sqlx)
