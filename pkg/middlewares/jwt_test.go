package middlewares

import (
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/aikidoaikido115/New-Acis-BE/configs"
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"
)

func generateTestToken(secret string, userID string, exp time.Time) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "exp":     exp.Unix(),
        "iat":     time.Now().Unix(),
        "jti":     uuid.New().String(),
    })

    tokenString, _ := token.SignedString([]byte(secret))
    return tokenString
}

func setupTestApp(config configs.JWT) *fiber.App {
    app := fiber.New()
    
    app.Use("/protected", JWTMiddleware(config))
    app.Get("/protected", func(c *fiber.Ctx) error {
        userID := c.Locals("user_id")
        return c.JSON(fiber.Map{
            "message": "Success",
            "user_id": userID,
        })
    })
    
    return app
}

// ========== VALID TOKEN TESTS ==========

func TestJWTMiddleware_ValidToken(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    token := generateTestToken(config.Secret, "user-123", time.Now().Add(30*time.Minute))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}

func TestJWTMiddleware_ValidToken_UserIDInLocals(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := fiber.New()

    var capturedUserID interface{}
    app.Use("/protected", JWTMiddleware(config))
    app.Get("/protected", func(c *fiber.Ctx) error {
        capturedUserID = c.Locals("user_id")
        return c.SendStatus(fiber.StatusOK)
    })

    token := generateTestToken(config.Secret, "user-456", time.Now().Add(30*time.Minute))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    app.Test(req)

    assert.Equal(t, "user-456", capturedUserID)
}

// ========== MISSING TOKEN TESTS ==========

func TestJWTMiddleware_MissingToken(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_EmptyAuthorization(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_InvalidAuthorizationFormat(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "InvalidFormat")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== EXPIRED TOKEN TESTS ==========

func TestJWTMiddleware_ExpiredToken(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    // Token that expired 1 hour ago
    token := generateTestToken(config.Secret, "user-123", time.Now().Add(-1*time.Hour))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== INVALID TOKEN TESTS ==========

func TestJWTMiddleware_InvalidToken(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer invalid.token.here")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_WrongSecret(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    // Token signed with different secret
    token := generateTestToken("wrong-secret-key", "user-123", time.Now().Add(30*time.Minute))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

func TestJWTMiddleware_MalformedToken(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    tests := []string{
        "Bearer ",
        "Bearer invalidtoken",
        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9",
        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.invalidpayload",
    }

    for _, tokenStr := range tests {
        t.Run(tokenStr, func(t *testing.T) {
            req := httptest.NewRequest(http.MethodGet, "/protected", nil)
            req.Header.Set("Authorization", tokenStr)

            resp, err := app.Test(req)

            assert.NoError(t, err)
            assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
        })
    }
}

// ========== TOKEN WITHOUT BEARER PREFIX TESTS ==========

func TestJWTMiddleware_TokenWithoutBearer(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    token := generateTestToken(config.Secret, "user-123", time.Now().Add(30*time.Minute))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", token) // Without "Bearer " prefix

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== DIFFERENT SIGNING METHODS TESTS ==========

func TestJWTMiddleware_WrongSigningMethod(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    // Create token with RS256 instead of HS256 (will fail validation)
    token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
        "user_id": "user-123",
        "exp":     time.Now().Add(30 * time.Minute).Unix(),
    })

    tokenString, _ := token.SignedString([]byte(config.Secret))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+tokenString)

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== TOKEN WITH MISSING CLAIMS TESTS ==========

func TestJWTMiddleware_TokenWithoutUserID(t *testing.T) {
    config := configs.JWT{Secret: "test-secret-key"}
    app := setupTestApp(config)

    // Token without user_id claim
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "exp": time.Now().Add(30 * time.Minute).Unix(),
        "iat": time.Now().Unix(),
    })

    tokenString, _ := token.SignedString([]byte(config.Secret))

    req := httptest.NewRequest(http.MethodGet, "/protected", nil)
    req.Header.Set("Authorization", "Bearer "+tokenString)

    resp, err := app.Test(req)

    assert.NoError(t, err)
    // Should still pass middleware (user_id will be nil in Locals)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}