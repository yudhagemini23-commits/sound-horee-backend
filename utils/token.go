package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a JWT token valid for 30 days.
func GenerateToken(uid string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["uid"] = uid
	// Kasih 72 jam (3 hari) biar Mas gak perlu login ulang tiap 15 menit pas lagi ngoding
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["iat"] = time.Now().Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// NowMillis returns current timestamp in Milliseconds.
// TAMBAHAN: Ini fungsi yang tadi error (undefined).
func NowMillis() int64 {
	return time.Now().UnixMilli()
}
