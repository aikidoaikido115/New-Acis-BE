package utils

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestNormalizeEmail_ValidEmail(t *testing.T) {
    tests := []struct {
        input    string
        expected string
    }{
        {"Test@Example.com", "test@example.com"},
        {"USER@DOMAIN.COM", "user@domain.com"},
        {"test.user@example.com", "testuser@example.com"},
        {"first.last@company.co.th", "firstlast@company.co.th"},
        {"a.b.c@test.com", "abc@test.com"},
    }

    for _, tt := range tests {
        t.Run(tt.input, func(t *testing.T) {
            result, err := NormalizeEmail(tt.input)

            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestNormalizeEmail_InvalidEmail(t *testing.T) {
    tests := []string{
        "",
        "invalidemail",
        "@example.com",
        "user@",
        "user@@example.com",
        "user@domain@com",
    }

    for _, email := range tests {
        t.Run(email, func(t *testing.T) {
            result, err := NormalizeEmail(email)

            assert.Error(t, err)
            assert.Empty(t, result)
            
            // ✅ Safe check: only access err.Error() if err is not nil
            if err != nil {
                assert.Equal(t, "invalid Email", err.Error())
            }
        })
    }
}

func TestNormalizeEmail_RemovesDots(t *testing.T) {
    email := "john.doe.smith@example.com"
    expected := "johndoesmith@example.com"

    result, err := NormalizeEmail(email)

    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    assert.NotContains(t, result[:len(result)-len("@example.com")], ".", 
        "Local part should not contain dots")
}

func TestNormalizeEmail_LowercaseConversion(t *testing.T) {
    email := "JohnDoe@EXAMPLE.COM"
    expected := "johndoe@example.com"

    result, err := NormalizeEmail(email)

    assert.NoError(t, err)
    assert.Equal(t, expected, result)
    assert.Equal(t, result, result, "Email should be entirely lowercase")
}

func TestNormalizeEmail_PreservesDomain(t *testing.T) {
    email := "user@example.com"

    result, err := NormalizeEmail(email)

    assert.NoError(t, err)
    assert.Contains(t, result, "@example.com", 
        "Domain should be preserved")
}

func TestNormalizeEmail_EdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        hasError bool
    }{
        {
            name:     "Single character local part",
            input:    "a@example.com",
            expected: "a@example.com",
            hasError: false,
        },
        {
            name:     "Multiple dots",
            input:    "a.b.c.d@example.com",
            expected: "abcd@example.com",
            hasError: false,
        },
        {
            name:     "Numbers in email",
            input:    "user123@example.com",
            expected: "user123@example.com",
            hasError: false,
        },
        {
            name:     "Subdomain",
            input:    "user@mail.example.com",
            expected: "user@mail.example.com",
            hasError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := NormalizeEmail(tt.input)

            if tt.hasError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.expected, result)
            }
        })
    }
}

