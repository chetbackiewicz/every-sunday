# Authentication and Authorization Research

**Status**: ✅ Complete  
**Date**: October 19, 2025  
**Decision**: JWT (JSON Web Token) Authentication with bcrypt Password Hashing

---

## Recommendation Summary

**Chosen Approach: JWT Token Authentication**

After evaluating JWT vs session-based authentication, **JWT with refresh tokens** is recommended for the budget tracking application.

**Key Decisions**:
1. ✅ **JWT for authentication** - Stateless, scalable, API-friendly
2. ✅ **bcrypt for password hashing** - Industry standard, cost factor 12
3. ✅ **golang-jwt/jwt/v5** - Most popular Go JWT library (8.6k stars)
4. ✅ **Access token (15 minutes)** - Short-lived for security
5. ✅ **Refresh token (7 days)** - Long-lived for UX
6. ✅ **HttpOnly cookies for refresh tokens** - Secure storage
7. ✅ **Gin middleware for auth** - Validate JWT on protected routes

**Why JWT over Sessions**:
- **Stateless**: No server-side storage required
- **Scalable**: Easy to distribute across multiple servers
- **API-friendly**: Works perfectly with React frontend + Go backend
- **Mobile-ready**: If mobile app needed in future
- **Cross-domain**: Can work across different domains

---

## Authentication Comparison

### Option 1: JWT (JSON Web Token)

**Overview**: Token-based authentication where server generates signed tokens containing user claims.

**How it works**:
1. User logs in with email/password
2. Server validates credentials
3. Server generates JWT with user claims (id, email, exp)
4. Client stores JWT (localStorage or memory)
5. Client sends JWT in Authorization header for each request
6. Server validates JWT signature and expiration
7. Server extracts user info from token claims

**JWT Structure**:
```
Header.Payload.Signature
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c
```

Decoded:
```json
// Header
{
  "alg": "HS256",
  "typ": "JWT"
}

// Payload (Claims)
{
  "sub": "1234567890",  // Subject (user ID)
  "name": "John Doe",
  "iat": 1516239022,    // Issued at
  "exp": 1516240922     // Expiration
}

// Signature (HMAC SHA256)
HMACSHA256(
  base64UrlEncode(header) + "." + base64UrlEncode(payload),
  secret
)
```

**Pros**:
- ✅ **Stateless**: No server-side session storage
- ✅ **Scalable**: Works across multiple servers without shared session store
- ✅ **Performance**: No database lookup for each request
- ✅ **API-friendly**: Perfect for REST API + SPA architecture
- ✅ **Cross-domain**: Can work across different domains
- ✅ **Mobile-ready**: Easy to use with mobile apps
- ✅ **Self-contained**: Token contains all user info

**Cons**:
- ⚠️ **Cannot revoke**: Once issued, token is valid until expiration
- ⚠️ **Token size**: Larger than session ID (stored in header)
- ⚠️ **Secret management**: Must protect JWT secret key
- ⚠️ **XSS vulnerability**: If stored in localStorage

**Verdict**: ✅ **Best fit** for our React + Go API architecture.

---

### Option 2: Session-Based Authentication

**Overview**: Server creates and stores session after login, sends session ID as cookie to client.

**How it works**:
1. User logs in with email/password
2. Server validates credentials
3. Server creates session in database/Redis
4. Server sends session ID as httpOnly cookie
5. Client automatically sends cookie with each request
6. Server looks up session in database
7. Server validates session and retrieves user info

**Pros**:
- ✅ **Easy to revoke**: Just delete session from database
- ✅ **Secure by default**: HttpOnly cookies prevent XSS
- ✅ **Small cookie size**: Only session ID stored
- ✅ **Server control**: Can track active sessions

**Cons**:
- ❌ **Stateful**: Requires server-side session storage
- ❌ **Database lookup**: Every request needs database query
- ❌ **Scaling complexity**: Requires shared session store (Redis)
- ❌ **CORS issues**: Cookies require same-site or specific CORS setup
- ❌ **Mobile unfriendly**: Cookies don't work well with native apps

**Verdict**: ⚠️ **Not ideal** for API architecture with React frontend.

---

## Decision Matrix

| Feature            | JWT          | Session-Based   |
| ------------------ | ------------ | --------------- |
| **Stateless**      | ✅            | ❌               |
| **Scalability**    | ✅            | ⚠️ Redis needed  |
| **Performance**    | ✅ Fast       | ⚠️ DB lookup     |
| **Revocation**     | ⚠️ Hard       | ✅ Easy          |
| **API-Friendly**   | ✅            | ⚠️               |
| **Mobile Support** | ✅            | ❌               |
| **Security**       | ✅ with HTTPS | ✅               |
| **Storage**        | Token ~200B  | Session ID ~32B |
| **Complexity**     | ⚠️ Medium     | ✅ Simple        |
| **CORS**           | ✅ Easy       | ⚠️ Tricky        |

**Winner: JWT** - Better fit for modern API + SPA architecture.

---

## JWT Implementation

### 1. Password Hashing with bcrypt

**Library**: `golang.org/x/crypto/bcrypt`

**Why bcrypt**:
- Industry standard for password hashing
- Adaptive algorithm (can increase cost factor)
- Built-in salt generation
- Resistant to rainbow table attacks
- Slow by design (prevents brute force)

**Cost Factor**:
- Recommended: **12-14** for production
- Higher = more secure but slower
- Cost 12 = ~250ms to hash
- Cost 14 = ~1 second to hash

```go
// utils/password.go
package utils

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

const DefaultCost = 12 // Recommended for production

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
    if len(password) > 72 {
        return "", fmt.Errorf("password exceeds bcrypt max length of 72 bytes")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(password), DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }

    return string(hash), nil
}

// CheckPassword compares a bcrypt hash with a plaintext password
func CheckPassword(hashedPassword, password string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        if err == bcrypt.ErrMismatchedHashAndPassword {
            return fmt.Errorf("invalid password")
        }
        return fmt.Errorf("failed to compare passwords: %w", err)
    }
    return nil
}

// ValidatePasswordStrength checks minimum password requirements
func ValidatePasswordStrength(password string) error {
    if len(password) < 8 {
        return fmt.Errorf("password must be at least 8 characters")
    }
    if len(password) > 72 {
        return fmt.Errorf("password cannot exceed 72 characters")
    }
    // Add more complexity requirements if needed
    return nil
}
```

**Usage Example**:
```go
// Register new user
password := "mySecretPassword123"

// Validate password strength
if err := ValidatePasswordStrength(password); err != nil {
    return err
}

// Hash password
hashedPassword, err := HashPassword(password)
if err != nil {
    return err
}

// Store hashedPassword in database
user, err := userRepo.Create(ctx, email, hashedPassword)

// Login - verify password
if err := CheckPassword(user.PasswordHash, password); err != nil {
    return errors.New("invalid credentials")
}
```

---

### 2. JWT Token Generation

**Library**: `github.com/golang-jwt/jwt/v5`

**Token Types**:
- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (7 days), used to get new access tokens

```go
// auth/jwt.go
package auth

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var (
    ErrInvalidToken      = errors.New("invalid token")
    ErrExpiredToken      = errors.New("token has expired")
    ErrInvalidSigningMethod = errors.New("invalid signing method")
)

type JWTConfig struct {
    SecretKey            string
    AccessTokenDuration  time.Duration
    RefreshTokenDuration time.Duration
}

// Custom claims for access token
type AccessTokenClaims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

// Custom claims for refresh token
type RefreshTokenClaims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// GenerateAccessToken creates a new access token for the user
func GenerateAccessToken(cfg JWTConfig, userID int, email string) (string, error) {
    claims := AccessTokenClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.AccessTokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            NotBefore: jwt.NewNumericDate(time.Now()),
            Issuer:    "budget-tracker-api",
            Subject:   fmt.Sprintf("%d", userID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(cfg.SecretKey))
    if err != nil {
        return "", fmt.Errorf("failed to sign access token: %w", err)
    }

    return signedToken, nil
}

// GenerateRefreshToken creates a new refresh token for the user
func GenerateRefreshToken(cfg JWTConfig, userID int) (string, error) {
    claims := RefreshTokenClaims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.RefreshTokenDuration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "budget-tracker-api",
            Subject:   fmt.Sprintf("%d", userID),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(cfg.SecretKey))
    if err != nil {
        return "", fmt.Errorf("failed to sign refresh token: %w", err)
    }

    return signedToken, nil
}

// GenerateTokenPair creates both access and refresh tokens
func GenerateTokenPair(cfg JWTConfig, userID int, email string) (*TokenPair, error) {
    accessToken, err := GenerateAccessToken(cfg, userID, email)
    if err != nil {
        return nil, err
    }

    refreshToken, err := GenerateRefreshToken(cfg, userID)
    if err != nil {
        return nil, err
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// ValidateAccessToken validates and parses an access token
func ValidateAccessToken(cfg JWTConfig, tokenString string) (*AccessTokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidSigningMethod
        }
        return []byte(cfg.SecretKey), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
    }

    if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, ErrInvalidToken
}

// ValidateRefreshToken validates and parses a refresh token
func ValidateRefreshToken(cfg JWTConfig, tokenString string) (*RefreshTokenClaims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrInvalidSigningMethod
        }
        return []byte(cfg.SecretKey), nil
    })

    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrExpiredToken
        }
        return nil, fmt.Errorf("%w: %v", ErrInvalidToken, err)
    }

    if claims, ok := token.Claims.(*RefreshTokenClaims); ok && token.Valid {
        return claims, nil
    }

    return nil, ErrInvalidToken
}
```

**Configuration**:
```go
// Load from environment variables
jwtConfig := auth.JWTConfig{
    SecretKey:            os.Getenv("JWT_SECRET"), // Strong random string
    AccessTokenDuration:  15 * time.Minute,         // Short-lived
    RefreshTokenDuration: 7 * 24 * time.Hour,       // 7 days
}

// Generate JWT_SECRET:
// openssl rand -base64 32
// Example: "X7GmDPFkQ8HvN2sL9RtYwK3uJ5cA1bE6"
```

---

### 3. Authentication Handlers

```go
// handlers/auth_handler.go
package handlers

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    
    "your-app/auth"
    "your-app/repositories"
    "your-app/utils"
)

type AuthHandler struct {
    userRepo  repositories.UserRepository
    jwtConfig auth.JWTConfig
}

func NewAuthHandler(userRepo repositories.UserRepository, jwtConfig auth.JWTConfig) *AuthHandler {
    return &AuthHandler{
        userRepo:  userRepo,
        jwtConfig: jwtConfig,
    }
}

// RegisterRequest defines the request body for user registration
type RegisterRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// LoginRequest defines the request body for user login
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

// AuthResponse is the response containing tokens
type AuthResponse struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token,omitempty"` // Omit if using cookie
    User         UserResponse `json:"user"`
}

type UserResponse struct {
    ID    int    `json:"id"`
    Email string `json:"email"`
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // Validate password strength
    if err := utils.ValidatePasswordStrength(req.Password); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "WEAK_PASSWORD",
                "message": err.Error(),
            },
        })
        return
    }

    // Check if user already exists
    existingUser, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to check user existence",
            },
        })
        return
    }
    if existingUser != nil {
        c.JSON(http.StatusConflict, gin.H{
            "error": gin.H{
                "code":    "USER_EXISTS",
                "message": "Email already registered",
            },
        })
        return
    }

    // Hash password
    hashedPassword, err := utils.HashPassword(req.Password)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to process password",
            },
        })
        return
    }

    // Create user
    user, err := h.userRepo.Create(c.Request.Context(), req.Email, hashedPassword)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to create user",
            },
        })
        return
    }

    // Generate tokens
    tokens, err := auth.GenerateTokenPair(h.jwtConfig, user.ID, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to generate tokens",
            },
        })
        return
    }

    // Set refresh token as httpOnly cookie
    h.setRefreshTokenCookie(c, tokens.RefreshToken)

    // Return access token and user info
    c.JSON(http.StatusCreated, AuthResponse{
        AccessToken: tokens.AccessToken,
        User: UserResponse{
            ID:    user.ID,
            Email: user.Email,
        },
    })
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": gin.H{
                "code":    "VALIDATION_ERROR",
                "message": err.Error(),
            },
        })
        return
    }

    // Get user by email
    user, err := h.userRepo.GetByEmail(c.Request.Context(), req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to lookup user",
            },
        })
        return
    }
    if user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "INVALID_CREDENTIALS",
                "message": "Invalid email or password",
            },
        })
        return
    }

    // Verify password
    if err := utils.CheckPassword(user.PasswordHash, req.Password); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "INVALID_CREDENTIALS",
                "message": "Invalid email or password",
            },
        })
        return
    }

    // Generate tokens
    tokens, err := auth.GenerateTokenPair(h.jwtConfig, user.ID, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to generate tokens",
            },
        })
        return
    }

    // Set refresh token as httpOnly cookie
    h.setRefreshTokenCookie(c, tokens.RefreshToken)

    // Return access token and user info
    c.JSON(http.StatusOK, AuthResponse{
        AccessToken: tokens.AccessToken,
        User: UserResponse{
            ID:    user.ID,
            Email: user.Email,
        },
    })
}

// RefreshToken handles token refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    // Get refresh token from cookie
    refreshToken, err := c.Cookie("refresh_token")
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "MISSING_REFRESH_TOKEN",
                "message": "Refresh token not found",
            },
        })
        return
    }

    // Validate refresh token
    claims, err := auth.ValidateRefreshToken(h.jwtConfig, refreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "INVALID_REFRESH_TOKEN",
                "message": err.Error(),
            },
        })
        return
    }

    // Get user from database
    user, err := h.userRepo.GetByID(c.Request.Context(), claims.UserID)
    if err != nil || user == nil {
        c.JSON(http.StatusUnauthorized, gin.H{
            "error": gin.H{
                "code":    "USER_NOT_FOUND",
                "message": "User no longer exists",
            },
        })
        return
    }

    // Generate new access token
    accessToken, err := auth.GenerateAccessToken(h.jwtConfig, user.ID, user.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": gin.H{
                "code":    "INTERNAL_ERROR",
                "message": "Failed to generate access token",
            },
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "access_token": accessToken,
    })
}

// Logout handles user logout
func (h *AuthHandler) Logout(c *gin.Context) {
    // Clear refresh token cookie
    c.SetCookie(
        "refresh_token",
        "",
        -1, // MaxAge -1 deletes the cookie
        "/",
        "",
        true,  // Secure
        true,  // HttpOnly
    )

    c.JSON(http.StatusOK, gin.H{
        "message": "Logged out successfully",
    })
}

// Me returns current user info
func (h *AuthHandler) Me(c *gin.Context) {
    // Get user from context (set by auth middleware)
    userID := c.GetInt("user_id")
    email := c.GetString("email")

    c.JSON(http.StatusOK, UserResponse{
        ID:    userID,
        Email: email,
    })
}

// setRefreshTokenCookie sets the refresh token as an httpOnly cookie
func (h *AuthHandler) setRefreshTokenCookie(c *gin.Context, refreshToken string) {
    c.SetCookie(
        "refresh_token",
        refreshToken,
        int(h.jwtConfig.RefreshTokenDuration.Seconds()),
        "/",
        "",    // Domain (empty for same domain)
        true,  // Secure (HTTPS only in production)
        true,  // HttpOnly (not accessible via JavaScript)
    )
}
```

---

### 4. Authentication Middleware

```go
// middleware/auth.go
package middleware

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    
    "your-app/auth"
)

// AuthMiddleware validates JWT access token
func AuthMiddleware(jwtConfig auth.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Get Authorization header
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "MISSING_TOKEN",
                    "message": "Authorization header required",
                },
            })
            c.Abort()
            return
        }

        // Check Bearer format
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "INVALID_TOKEN_FORMAT",
                    "message": "Authorization header must be Bearer {token}",
                },
            })
            c.Abort()
            return
        }

        tokenString := parts[1]

        // Validate token
        claims, err := auth.ValidateAccessToken(jwtConfig, tokenString)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{
                "error": gin.H{
                    "code":    "INVALID_TOKEN",
                    "message": err.Error(),
                },
            })
            c.Abort()
            return
        }

        // Set user info in context
        c.Set("user_id", claims.UserID)
        c.Set("email", claims.Email)

        c.Next()
    }
}

// OptionalAuthMiddleware validates JWT but doesn't require it
func OptionalAuthMiddleware(jwtConfig auth.JWTConfig) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.Next()
            return
        }

        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.Next()
            return
        }

        claims, err := auth.ValidateAccessToken(jwtConfig, parts[1])
        if err == nil {
            c.Set("user_id", claims.UserID)
            c.Set("email", claims.Email)
        }

        c.Next()
    }
}
```

---

### 5. Complete Route Setup

```go
// main.go
package main

import (
    "log"
    "os"
    "time"

    "github.com/gin-gonic/gin"
    
    "your-app/auth"
    "your-app/database"
    "your-app/handlers"
    "your-app/middleware"
    "your-app/repositories"
)

func main() {
    // JWT configuration
    jwtConfig := auth.JWTConfig{
        SecretKey:            os.Getenv("JWT_SECRET"),
        AccessTokenDuration:  15 * time.Minute,
        RefreshTokenDuration: 7 * 24 * time.Hour,
    }

    // Database setup (from doc 13)
    db, err := database.NewPostgresDB(database.DefaultConfig())
    if err != nil {
        log.Fatalf("Database connection failed: %v", err)
    }
    defer db.Close()

    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    budgetRepo := repositories.NewBudgetRepository(db)

    // Initialize handlers
    authHandler := handlers.NewAuthHandler(userRepo, jwtConfig)
    budgetHandler := handlers.NewBudgetHandler(budgetRepo)

    // Setup Gin router
    r := gin.Default()

    // CORS middleware
    r.Use(middleware.CORS())

    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Public auth routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.POST("/refresh", authHandler.RefreshToken)
            auth.POST("/logout", authHandler.Logout)
        }

        // Protected routes
        protected := v1.Group("")
        protected.Use(middleware.AuthMiddleware(jwtConfig))
        {
            // User info
            protected.GET("/auth/me", authHandler.Me)

            // Budget routes
            budgets := protected.Group("/budgets")
            {
                budgets.GET("", budgetHandler.List)
                budgets.POST("", budgetHandler.Create)
                budgets.GET("/:month", budgetHandler.Get)
                budgets.PUT("/:month", budgetHandler.Update)
                budgets.DELETE("/:month", budgetHandler.Delete)
            }
        }
    }

    // Start server
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

---

## Security Best Practices

### 1. Secret Key Management

```bash
# Generate strong JWT secret (32 bytes)
openssl rand -base64 32

# Store in .env file
JWT_SECRET=X7GmDPFkQ8HvN2sL9RtYwK3uJ5cA1bE6

# Never commit .env to git
echo ".env" >> .gitignore
```

### 2. HTTPS Only

```go
// In production, enforce HTTPS
if gin.Mode() == gin.ReleaseMode {
    r.Use(middleware.SecureHeaders())
}

// middleware/secure.go
func SecureHeaders() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Enforce HTTPS
        c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        // Prevent MIME sniffing
        c.Header("X-Content-Type-Options", "nosniff")
        // XSS protection
        c.Header("X-Frame-Options", "DENY")
        c.Next()
    }
}
```

### 3. Token Storage (Frontend)

**Recommended Approach**:
- **Access Token**: Store in memory (React state/context) - Most secure, lost on refresh
- **Refresh Token**: HttpOnly cookie - Protected from XSS

**Alternative**:
- **Access Token**: localStorage - Persistent but XSS vulnerable
- **Refresh Token**: HttpOnly cookie

```typescript
// React: Store access token in memory
const AuthContext = React.createContext();

function AuthProvider({ children }) {
  const [accessToken, setAccessToken] = useState(null);

  // Refresh token automatically before expiration
  useEffect(() => {
    const interval = setInterval(async () => {
      if (accessToken) {
        const response = await fetch('/api/v1/auth/refresh', {
          method: 'POST',
          credentials: 'include', // Send cookies
        });
        const data = await response.json();
        setAccessToken(data.access_token);
      }
    }, 14 * 60 * 1000); // Refresh every 14 minutes (before 15min expiration)

    return () => clearInterval(interval);
  }, [accessToken]);

  return (
    <AuthContext.Provider value={{ accessToken, setAccessToken }}>
      {children}
    </AuthContext.Provider>
  );
}
```

### 4. CORS Configuration

```go
// middleware/cors.go
package middleware

import (
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "time"
)

func CORS() gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"}, // React dev server
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true, // Allow cookies
        MaxAge:           12 * time.Hour,
    })
}
```

### 5. Rate Limiting

```go
// middleware/rate_limit.go
package middleware

import (
    "net/http"
    "sync"
    "time"

    "github.com/gin-gonic/gin"
)

type visitor struct {
    lastSeen time.Time
    count    int
}

var visitors = make(map[string]*visitor)
var mu sync.Mutex

func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ip := c.ClientIP()
        
        mu.Lock()
        v, exists := visitors[ip]
        if !exists {
            visitors[ip] = &visitor{
                lastSeen: time.Now(),
                count:    1,
            }
            mu.Unlock()
            c.Next()
            return
        }

        if time.Since(v.lastSeen) > window {
            v.lastSeen = time.Now()
            v.count = 1
        } else {
            v.count++
        }

        if v.count > maxRequests {
            mu.Unlock()
            c.JSON(http.StatusTooManyRequests, gin.H{
                "error": gin.H{
                    "code":    "RATE_LIMIT_EXCEEDED",
                    "message": "Too many requests, please try again later",
                },
            })
            c.Abort()
            return
        }

        mu.Unlock()
        c.Next()
    }
}

// Apply to auth routes
auth.POST("/login", middleware.RateLimit(5, time.Minute), authHandler.Login)
auth.POST("/register", middleware.RateLimit(3, time.Hour), authHandler.Register)
```

---

## Testing

### Unit Tests for Auth Functions

```go
// auth/jwt_test.go
package auth_test

import (
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    
    "your-app/auth"
)

func TestGenerateAndValidateAccessToken(t *testing.T) {
    cfg := auth.JWTConfig{
        SecretKey:           "test-secret-key",
        AccessTokenDuration: 15 * time.Minute,
    }

    // Generate token
    userID := 123
    email := "test@example.com"
    token, err := auth.GenerateAccessToken(cfg, userID, email)
    require.NoError(t, err)
    require.NotEmpty(t, token)

    // Validate token
    claims, err := auth.ValidateAccessToken(cfg, token)
    require.NoError(t, err)
    assert.Equal(t, userID, claims.UserID)
    assert.Equal(t, email, claims.Email)
}

func TestValidateExpiredToken(t *testing.T) {
    cfg := auth.JWTConfig{
        SecretKey:           "test-secret-key",
        AccessTokenDuration: -1 * time.Hour, // Already expired
    }

    token, err := auth.GenerateAccessToken(cfg, 123, "test@example.com")
    require.NoError(t, err)

    // Should fail validation
    _, err = auth.ValidateAccessToken(cfg, token)
    assert.ErrorIs(t, err, auth.ErrExpiredToken)
}
```

### Integration Tests for Auth Handler

```go
// handlers/auth_handler_test.go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    "your-app/auth"
    "your-app/handlers"
    "your-app/repositories/mocks"
)

func TestRegister_Success(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)
    mockRepo := new(mocks.MockUserRepository)
    jwtConfig := auth.JWTConfig{
        SecretKey:            "test-secret",
        AccessTokenDuration:  15 * time.Minute,
        RefreshTokenDuration: 7 * 24 * time.Hour,
    }
    handler := handlers.NewAuthHandler(mockRepo, jwtConfig)

    // Mock expectations
    mockRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", mock.Anything, "test@example.com", mock.AnythingOfType("string")).
        Return(&repositories.User{ID: 1, Email: "test@example.com"}, nil)

    // Request
    reqBody := handlers.RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
    }
    jsonBody, _ := json.Marshal(reqBody)

    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(jsonBody))
    c.Request.Header.Set("Content-Type", "application/json")

    // Execute
    handler.Register(c)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    
    var response handlers.AuthResponse
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.NotEmpty(t, response.AccessToken)
    assert.Equal(t, "test@example.com", response.User.Email)

    mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
    // Similar test structure...
}
```

---

## Pros of JWT Authentication

- ✅ **Stateless**: No server-side session storage
- ✅ **Scalable**: Works across multiple servers
- ✅ **API-friendly**: Perfect for REST API + React
- ✅ **Performance**: No database lookup per request
- ✅ **Mobile-ready**: Easy to use with native apps
- ✅ **Flexibility**: Can add custom claims (roles, permissions)
- ✅ **Industry standard**: Widely adopted, many libraries

---

## Cons/Challenges

- ⚠️ **Token revocation**: Hard to invalidate before expiration
  - Solution: Short access token lifespan (15 min)
  - Solution: Refresh token rotation
- ⚠️ **Token size**: Larger than session ID (~200 bytes)
  - Impact: Minimal for HTTP/2
- ⚠️ **Secret management**: Must protect JWT secret
  - Solution: Environment variables, secret manager
- ⚠️ **XSS vulnerability**: If stored in localStorage
  - Solution: Use memory storage + httpOnly refresh token

---

## Next Steps

1. **Implement password hashing** (utils/password.go)
2. **Create JWT functions** (auth/jwt.go)
3. **Build auth handlers** (handlers/auth_handler.go)
4. **Add auth middleware** (middleware/auth.go)
5. **Setup auth routes** (main.go)
6. **Add CORS middleware** (middleware/cors.go)
7. **Implement rate limiting** (middleware/rate_limit.go on auth routes)
8. **Write tests** (auth_test.go, handlers_test.go)

---

## References

- [JWT Introduction](https://jwt.io/introduction)
- [golang-jwt/jwt Documentation](https://pkg.go.dev/github.com/golang-jwt/jwt/v5)
- [bcrypt Package Documentation](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
- [OWASP Authentication Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html)
- [RFC 7519 - JSON Web Token](https://datatracker.ietf.org/doc/html/rfc7519)

---

**Implementation Complete**: All research documents finished (14/14). Ready to begin implementation phase using Gin, PostgreSQL, sqlx, golang-migrate, and JWT authentication.
