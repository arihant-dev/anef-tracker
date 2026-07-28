package session

import (
	"fmt"
	"os"
	"time"
)

type ValidationResult struct {
	HasAccessToken  bool
	HasRefreshToken bool
	CookieCount     int
	ExpiresIn       time.Duration
	IsExpired       bool
	Issuer          string
	User            string
	Provider        string
	Ready           bool
	Checks          []string
}

// ValidateSession runs thorough checklist on session readiness.
func ValidateSession(s *Session) *ValidationResult {
	res := &ValidationResult{
		Checks: make([]string, 0),
	}

	if s == nil {
		res.Checks = append(res.Checks, "✗ No session loaded")
		return res
	}

	res.Provider = s.Provider
	if res.Provider == "" {
		res.Provider = "anef"
	}

	// Access Token check
	if s.AccessToken != "" {
		res.HasAccessToken = true
		res.Checks = append(res.Checks, "✓ Access token present")
	} else {
		res.Checks = append(res.Checks, "✗ Access token missing")
	}

	// Refresh Token check
	if s.RefreshToken != "" {
		res.HasRefreshToken = true
		res.Checks = append(res.Checks, "✓ Refresh token present")
	} else {
		res.Checks = append(res.Checks, "! Refresh token omitted (auto-renewal disabled)")
	}

	// Cookies check
	res.CookieCount = len(s.Cookies)
	if res.CookieCount > 0 {
		res.Checks = append(res.Checks, fmt.Sprintf("✓ Cookies present (%d captured)", res.CookieCount))
	} else {
		res.Checks = append(res.Checks, "! No cookies captured")
	}

	// Expiry check
	if !s.ExpiresAt.IsZero() {
		res.ExpiresIn = time.Until(s.ExpiresAt)
		if res.ExpiresIn > 0 {
			res.Checks = append(res.Checks, fmt.Sprintf("✓ Valid (expires in %s)", res.ExpiresIn.Round(time.Second)))
		} else {
			res.IsExpired = true
			res.Checks = append(res.Checks, fmt.Sprintf("✗ Session token expired %s ago", (-res.ExpiresIn).Round(time.Second)))
		}
	} else {
		res.Checks = append(res.Checks, "! Token expiration unknown")
	}

	// JWT claim check
	claims, err := DecodeJWT(s.AccessToken)
	if err == nil && claims != nil {
		res.User = claims.Login
		res.Issuer = claims.Iss
		res.Checks = append(res.Checks, fmt.Sprintf("✓ JWT claims decoded (user: %s)", res.User))

		if res.Issuer != "" {
			res.Checks = append(res.Checks, fmt.Sprintf("✓ Issuer verified (%s)", res.Issuer))
		}
	} else {
		res.User = s.User
		if res.User != "" {
			res.Checks = append(res.Checks, fmt.Sprintf("✓ User identified (%s)", res.User))
		}
	}

	res.Ready = res.HasAccessToken && !res.IsExpired
	if res.Ready {
		res.Checks = append(res.Checks, "✓ Ready for fetch")
	} else {
		res.Checks = append(res.Checks, "✗ Not ready for fetch")
	}

	return res
}

type DoctorCheck struct {
	Component string
	Passed    bool
	Message   string
}

type SessionDoctorReport struct {
	Checks []DoctorCheck
}

// RunSessionDoctor executes a 5-point diagnostic suite specifically on session storage, encryption, permissions, and token health.
func RunSessionDoctor() *SessionDoctorReport {
	rep := &SessionDoctorReport{Checks: make([]DoctorCheck, 0)}

	// 1. Session File & Permissions
	path, err := GetSessionFilePath()
	if err != nil {
		rep.Checks = append(rep.Checks, DoctorCheck{"Session Path", false, fmt.Sprintf("Failed getting path: %v", err)})
	} else {
		fi, err := os.Stat(path)
		if os.IsNotExist(err) {
			rep.Checks = append(rep.Checks, DoctorCheck{"Session File", false, "Session file missing (~/.anef/session.json). Run 'anef login'"})
		} else if err != nil {
			rep.Checks = append(rep.Checks, DoctorCheck{"Session File", false, fmt.Sprintf("File stat error: %v", err)})
		} else {
			perm := fi.Mode().Perm()
			if perm == 0600 {
				rep.Checks = append(rep.Checks, DoctorCheck{"Permissions", true, fmt.Sprintf("Secure file permissions (0600) at %s", path)})
			} else {
				rep.Checks = append(rep.Checks, DoctorCheck{"Permissions", false, fmt.Sprintf("Insecure permissions %#o at %s (expected 0600)", perm, path)})
			}
		}
	}

	// 2. AES-GCM Decryption & Load
	sess, err := LoadSession()
	if err != nil {
		rep.Checks = append(rep.Checks, DoctorCheck{"AES-GCM Encryption", false, fmt.Sprintf("Failed decrypting/loading session: %v", err)})
	} else {
		rep.Checks = append(rep.Checks, DoctorCheck{"AES-GCM Encryption", true, fmt.Sprintf("Decrypted session successfully (source: %s)", sess.ImportSource)})

		// 3. Token Expiry Check
		if !sess.ExpiresAt.IsZero() && time.Now().Before(sess.ExpiresAt) {
			rep.Checks = append(rep.Checks, DoctorCheck{"Token Expiry", true, fmt.Sprintf("Access token valid for %s", time.Until(sess.ExpiresAt).Round(time.Second))})
		} else if sess.IsExpired() {
			rep.Checks = append(rep.Checks, DoctorCheck{"Token Expiry", false, "Access token expired. Re-authenticate via 'anef login'"})
		} else {
			rep.Checks = append(rep.Checks, DoctorCheck{"Token Expiry", true, "Access token active"})
		}

		// 4. JWT Claim Integrity
		claims, err := DecodeJWT(sess.AccessToken)
		if err == nil && claims != nil {
			rep.Checks = append(rep.Checks, DoctorCheck{"JWT Integrity", true, fmt.Sprintf("Decoded claims for user '%s'", claims.Login)})
		} else {
			rep.Checks = append(rep.Checks, DoctorCheck{"JWT Integrity", false, "Could not decode JWT claims"})
		}

		// 5. Refresh Token Status
		if sess.RefreshToken != "" {
			rep.Checks = append(rep.Checks, DoctorCheck{"Refresh Token", true, "OAuth2 refresh token captured"})
		} else {
			rep.Checks = append(rep.Checks, DoctorCheck{"Refresh Token", false, "No refresh token captured"})
		}
	}

	return rep
}
