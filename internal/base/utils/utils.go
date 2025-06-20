package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ExtractUserIDFromToken(tokenString string) (uuid.UUID, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("token parsing failed: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token claims")
	}

	userIDStr, ok := claims["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("user_id not found in token")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, errors.New("invalid user ID in token")
	}

	return userID, nil
}

func FormatEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// FormatMobileNumber formats a mobile number to the standard Saudi Arabia format (966xxxxxxxxx)
// It handles the following cases:
// - 5006054839 -> 966506054839 (add country code)
// - 0506054839 -> 966506054839 (replace leading 0 with country code)
// - +966506054839 -> 966506054839 (remove + sign)
func FormatMobileNumber(mobile string) string {
	mobile = strings.TrimSpace(mobile)

	mobile = strings.TrimPrefix(mobile, "+")

	if strings.HasPrefix(mobile, "0") {
		mobile = mobile[1:]
	}
	if !strings.HasPrefix(mobile, "966") {
		if len(mobile) >= 9 && (strings.HasPrefix(mobile, "5") || strings.HasPrefix(mobile, "4") || strings.HasPrefix(mobile, "3")) {
			mobile = "966" + mobile
		}
	}

	return mobile
}
