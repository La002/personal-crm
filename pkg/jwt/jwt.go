package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint, email string, secret string, expiryHours int) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["email"] = email
	claims["exp"] = now.Add(time.Duration(expiryHours) * time.Hour).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(tokenString string, secret string) (*jwt.MapClaims, error) {
	tok, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signed method: %s", jwtToken.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	if !tok.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return &claims, nil
}

func ExtractUserID(claims *jwt.MapClaims) (uint, error) {
	if claims == nil {
		return 0, fmt.Errorf("claims cannot be nil")
	}

	val, ok := (*claims)["user_id"]
	if !ok {
		return 0, fmt.Errorf("user_id not found in claims")
	}

	// JWT numeric claims are decoded as float64
	floatVal, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("user_id has invalid type")
	}

	return uint(floatVal), nil
}
