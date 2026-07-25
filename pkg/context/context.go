package context

import (
	"fmt"
)

type Scope struct {
	ProfileID     int64 `json:"profile_id" yaml:"profile_id"`
	ApplicationID int64 `json:"application_id" yaml:"application_id"`
}

func DefaultScope() Scope {
	return Scope{
		ProfileID:     1,
		ApplicationID: 1,
	}
}

func NewScope(profileID, applicationID int64) (Scope, error) {
	if profileID <= 0 {
		return Scope{}, fmt.Errorf("invalid profileID: must be > 0")
	}
	if applicationID <= 0 {
		return Scope{}, fmt.Errorf("invalid applicationID: must be > 0")
	}
	return Scope{
		ProfileID:     profileID,
		ApplicationID: applicationID,
	}, nil
}

func (s Scope) String() string {
	return fmt.Sprintf("Profile #%d | App #%d", s.ProfileID, s.ApplicationID)
}
