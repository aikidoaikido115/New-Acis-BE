package utils

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

func TestGenerateRandomOTP_ValidLength(t *testing.T) {
    otp, err := GenerateRandomOTP(6)

    assert.NoError(t, err)
    assert.Len(t, otp, 6)
    assert.Regexp(t, "^[0-9]{6}$", otp, "OTP should contain only digits")
}

func TestGenerateRandomOTP_DifferentLengths(t *testing.T) {
    tests := []struct {
        length   int
        expected int
    }{
        {4, 4},
        {6, 6},
        {8, 8},
        {10, 10},
    }

    for _, tt := range tests {
        t.Run("Length"+string(rune(tt.length)), func(t *testing.T) {
            otp, err := GenerateRandomOTP(tt.length)

            assert.NoError(t, err)
            assert.Len(t, otp, tt.expected)
        })
    }
}

func TestGenerateRandomOTP_InvalidLength(t *testing.T) {
    tests := []int{0, -1, -5}

    for _, length := range tests {
        t.Run("InvalidLength", func(t *testing.T) {
            otp, err := GenerateRandomOTP(length)

            assert.Error(t, err)
            assert.Empty(t, otp)
            assert.Equal(t, "OTP length must be greater than 0", err.Error())
        })
    }
}

func TestGenerateRandomOTP_Uniqueness(t *testing.T) {
    // Generate multiple OTPs and check they're different
    otps := make(map[string]bool)
    iterations := 100

    for i := 0; i < iterations; i++ {
        otp, err := GenerateRandomOTP(6)
        assert.NoError(t, err)
        otps[otp] = true
    }

    // Statistical test: at least 90% should be unique
    uniqueCount := len(otps)
    assert.Greater(t, uniqueCount, iterations*9/10, 
        "OTPs should be mostly unique (at least 90%%)")
}

func TestGenerateRandomOTP_OnlyDigits(t *testing.T) {
    otp, err := GenerateRandomOTP(20)

    assert.NoError(t, err)
    
    // Check each character is a digit
    for _, char := range otp {
        assert.True(t, char >= '0' && char <= '9', 
            "Each character should be a digit (0-9)")
    }
}

func TestGenerateRandomOTP_Randomness(t *testing.T) {
    // Generate two OTPs and ensure they're different
    otp1, err1 := GenerateRandomOTP(6)
    otp2, err2 := GenerateRandomOTP(6)

    assert.NoError(t, err1)
    assert.NoError(t, err2)
    
    // While there's a tiny chance they could be the same,
    // it's statistically extremely unlikely (1 in 1,000,000)
    assert.NotEqual(t, otp1, otp2, 
        "Two consecutive OTPs should be different")
}