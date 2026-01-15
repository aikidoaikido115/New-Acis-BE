package controllers

import (
    "bytes"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/aikidoaikido115/New-Acis-BE/modules/entities"
    "github.com/aikidoaikido115/New-Acis-BE/modules/user/usecases"
    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func setupTestController() (*UserController, *usecases.MockUserUsecase, *fiber.App) {
    mockUsecase := new(usecases.MockUserUsecase)
    controller := NewUserController(mockUsecase)
    app := fiber.New()
    return controller, mockUsecase, app
}

// ========== REGISTER HANDLER TESTS ==========

func TestRegisterHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/register", controller.RegisterHandler)

    user := &entities.User{
        ID:       "user-1",
        Username: "testuser",
        Email:    "test@example.com",
    }

    mockUsecase.On("Register", mock.AnythingOfType("*entities.User"), "Medical Staff").Return(user, nil)

    reqBody := map[string]string{
        "username":  "testuser",
        "email":     "test@example.com",
        "password":  "password123",
        "role_name": "Medical Staff",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Success", response["status"])
    mockUsecase.AssertExpectations(t)
}

func TestRegisterHandler_MissingUsername(t *testing.T) {
    controller, _, app := setupTestController()

    app.Post("/register", controller.RegisterHandler)

    reqBody := map[string]string{
        "email":     "test@example.com",
        "password":  "password123",
        "role_name": "Medical Staff",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Username is missing", response["message"])
}

func TestRegisterHandler_UsernameExists(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/register", controller.RegisterHandler)

    mockUsecase.On("Register", mock.AnythingOfType("*entities.User"), "Medical Staff").
        Return(nil, errors.New("username already exists"))

    reqBody := map[string]string{
        "username":  "existinguser",
        "email":     "test@example.com",
        "password":  "password123",
        "role_name": "Medical Staff",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
    mockUsecase.AssertExpectations(t)
}

// ========== LOGIN HANDLER TESTS ==========

func TestLoginHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/login", controller.LoginHandler)

    user := &entities.User{
        ID:       "user-1",
        Username: "testuser",
        Email:    "test@example.com",
    }

    mockUsecase.On("Login", "testuser", "password123").
        Return("jwt-token-string", user, nil)

    reqBody := map[string]string{
        "username": "testuser",
        "password": "password123",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Success", response["status"])
    assert.NotNil(t, response["result"])
    mockUsecase.AssertExpectations(t)
}

func TestLoginHandler_MissingPassword(t *testing.T) {
    controller, _, app := setupTestController()

    app.Post("/login", controller.LoginHandler)

    reqBody := map[string]string{
        "username": "testuser",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/login", controller.LoginHandler)

    mockUsecase.On("Login", "testuser", "wrongpassword").
        Return("", (*entities.User)(nil), errors.New("invalid password"))

    reqBody := map[string]string{
        "username": "testuser",
        "password": "wrongpassword",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
    mockUsecase.AssertExpectations(t)
}

// ========== RESET PASSWORD HANDLER TESTS ==========

func TestResetPasswordHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Put("/resetpassword", func(c *fiber.Ctx) error {
        c.Locals("user_id", "user-1")
        return controller.ResetPasswordHandler(c)
    })

    mockUsecase.On("ResetPassword", "user-1", "oldpassword", "newpassword").Return(nil)

    reqBody := map[string]string{
        "email":        "test@example.com",
        "old_password": "oldpassword",
        "new_password": "newpassword",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPut, "/resetpassword", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)
    mockUsecase.AssertExpectations(t)
}

func TestResetPasswordHandler_Unauthorized(t *testing.T) {
    controller, _, app := setupTestController()

    app.Put("/resetpassword", controller.ResetPasswordHandler)

    reqBody := map[string]string{
        "old_password": "oldpassword",
        "new_password": "newpassword",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPut, "/resetpassword", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== FORGOT PASSWORD HANDLER TESTS ==========

func TestForgotPasswordHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/forgotpassword", controller.ForgotPasswordHandler)

    mockUsecase.On("ForgotPassword", "test@example.com").Return(nil)

    reqBody := map[string]string{
        "email": "test@example.com",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/forgotpassword", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Sent OTP successfully", response["message"])
    mockUsecase.AssertExpectations(t)
}

func TestForgotPasswordHandler_MissingEmail(t *testing.T) {
    controller, _, app := setupTestController()

    app.Post("/forgotpassword", controller.ForgotPasswordHandler)

    reqBody := map[string]string{}
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/forgotpassword", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
}

// ========== VERIFY OTP HANDLER TESTS ==========

func TestVerifyOTPHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/verify-otp", controller.VerifyOTPHandler)

    mockUsecase.On("VerifyOTP", "test@example.com", "123456").Return(nil)

    reqBody := map[string]string{
        "email": "test@example.com",
        "otp":   "123456",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "OTP is correct", response["message"])
    mockUsecase.AssertExpectations(t)
}

func TestVerifyOTPHandler_InvalidOTP(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Post("/verify-otp", controller.VerifyOTPHandler)

    mockUsecase.On("VerifyOTP", "test@example.com", "wrong-otp").
        Return(errors.New("invalid OTP code"))

    reqBody := map[string]string{
        "email": "test@example.com",
        "otp":   "wrong-otp",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest(http.MethodPost, "/verify-otp", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
    mockUsecase.AssertExpectations(t)
}

// ========== GET USER BY ID HANDLER TESTS ==========

func TestGetUserByIDHandler_Success(t *testing.T) {
    controller, mockUsecase, app := setupTestController()

    app.Get("/user", func(c *fiber.Ctx) error {
        c.Locals("user_id", "user-1")
        return controller.GetUserByIDHandler(c)
    })

    user := &entities.User{
        ID:       "user-1",
        Username: "testuser",
        Email:    "test@example.com",
    }

    mockUsecase.On("GetUserByID", "user-1").Return(user, nil)

    req := httptest.NewRequest(http.MethodGet, "/user", nil)
    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&response)
    assert.Equal(t, "Success", response["status"])
    mockUsecase.AssertExpectations(t)
}

func TestGetUserByIDHandler_Unauthorized(t *testing.T) {
    controller, _, app := setupTestController()

    app.Get("/user", controller.GetUserByIDHandler)

    req := httptest.NewRequest(http.MethodGet, "/user", nil)
    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
}

// ========== LOGOUT HANDLER TESTS ==========

func TestLogoutHandler_Success(t *testing.T) {
    controller, _, app := setupTestController()

    app.Post("/logout", controller.LogoutHandler)

    req := httptest.NewRequest(http.MethodPost, "/logout", nil)
    resp, err := app.Test(req)

    assert.NoError(t, err)
    assert.Equal(t, fiber.StatusOK, resp.StatusCode)

    var response map[string]interface{}
    bodyBytes, _ := io.ReadAll(resp.Body)
    json.Unmarshal(bodyBytes, &response)
    assert.Equal(t, "Logout successful", response["message"])
}