package importer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type JWTClaims struct {
	Login             string `json:"login"`
	Sub               string `json:"sub"`
	PreferredUsername string `json:"preferred_username"`
	NumEtranger        string `json:"numero_etranger"`
	Iat               int64  `json:"iat"`
	Exp               int64  `json:"exp"`
	Iss               string `json:"iss"`
}

// DecodeJWT decodes unverified JWT token claims without validating signature.
func DecodeJWT(tokenStr string) (*JWTClaims, error) {
	tokenStr = strings.TrimSpace(tokenStr)
	if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
		tokenStr = strings.TrimSpace(tokenStr[7:])
	} else if strings.HasPrefix(strings.ToLower(tokenStr), "token ") {
		tokenStr = strings.TrimSpace(tokenStr[6:])
	}

	parts := strings.Split(tokenStr, ".")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid jwt format")
	}

	payload := parts[1]
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	data, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		data, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return nil, err
		}
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	claims := &JWTClaims{}
	if val, ok := raw["login"].(string); ok {
		claims.Login = val
	}
	if val, ok := raw["sub"].(string); ok {
		claims.Sub = val
	}
	if val, ok := raw["preferred_username"].(string); ok {
		claims.PreferredUsername = val
	}
	if val, ok := raw["numero_etranger"].(string); ok {
		claims.NumEtranger = val
	}
	if val, ok := raw["iss"].(string); ok {
		claims.Iss = val
	}

	if val, ok := raw["iat"].(float64); ok {
		claims.Iat = int64(val)
	}
	if val, ok := raw["exp"].(float64); ok {
		claims.Exp = int64(val)
	}

	if claims.Login == "" {
		if claims.PreferredUsername != "" {
			claims.Login = claims.PreferredUsername
		} else if claims.NumEtranger != "" {
			claims.Login = claims.NumEtranger
		} else if claims.Sub != "" {
			claims.Login = claims.Sub
		}
	}

	return claims, nil
}

func GetUserFromJWT(token string) string {
	claims, err := DecodeJWT(token)
	if err != nil || claims == nil {
		return ""
	}
	return claims.Login
}

func GetExpiryFromJWT(token string) time.Time {
	claims, err := DecodeJWT(token)
	if err != nil || claims == nil || claims.Exp == 0 {
		return time.Time{}
	}
	return time.Unix(claims.Exp, 0)
}
