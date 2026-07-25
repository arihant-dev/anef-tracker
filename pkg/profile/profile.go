package profile

import (
	"fmt"
	"github.com/arihant-dev/anef-tracker/pkg/db"
	"strings"
	"time"
)

type Profile struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Active    bool      `json:"active"`
}

type TrackedApplication struct {
	ID        int64     `json:"id"`
	ProfileID int64     `json:"profile_id"`
	Name      string    `json:"name"`
	AnefID    string    `json:"anef_id,omitempty"`
	Type      string    `json:"type"`
	Status    string    `json:"status,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateProfile(database *db.DB, name string) (*Profile, error) {
	if database == nil {
		return &Profile{ID: 1, Name: name, CreatedAt: time.Now(), Active: true}, nil
	}

	res, err := database.Conn.Exec("INSERT INTO profiles (name, created_at, active) VALUES (?, ?, 0)", name, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed creating profile: %w", err)
	}
	id, _ := res.LastInsertId()
	return &Profile{ID: id, Name: name, CreatedAt: time.Now(), Active: false}, nil
}

func ListProfiles(database *db.DB) ([]Profile, error) {
	if database == nil {
		return []Profile{{ID: 1, Name: "Default Profile", CreatedAt: time.Now(), Active: true}}, nil
	}

	rows, err := database.Conn.Query("SELECT id, name, created_at, active FROM profiles ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []Profile
	for rows.Next() {
		var p Profile
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.Active); err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}
	return profiles, rows.Err()
}

func SwitchProfile(database *db.DB, profileID int64) error {
	if database == nil {
		return nil
	}

	tx, err := database.Conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, _ = tx.Exec("UPDATE profiles SET active = 0")
	res, err := tx.Exec("UPDATE profiles SET active = 1 WHERE id = ?", profileID)
	if err != nil {
		return err
	}
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("profile #%d not found", profileID)
	}

	return tx.Commit()
}

func GetActiveProfile(database *db.DB) (*Profile, error) {
	if database == nil {
		return &Profile{ID: 1, Name: "Default Profile", CreatedAt: time.Now(), Active: true}, nil
	}

	var p Profile
	err := database.Conn.QueryRow("SELECT id, name, created_at, active FROM profiles WHERE active = 1 LIMIT 1").Scan(&p.ID, &p.Name, &p.CreatedAt, &p.Active)
	if err != nil {
		return &Profile{ID: 1, Name: "Default Profile", CreatedAt: time.Now(), Active: true}, nil
	}
	return &p, nil
}

func DeleteProfile(database *db.DB, profileID int64, confirmID int64) error {
	if profileID != confirmID {
		return fmt.Errorf("confirmation mismatch: profileID (%d) != confirmID (%d)", profileID, confirmID)
	}
	if database == nil {
		return nil
	}
	_, err := database.Conn.Exec("DELETE FROM profiles WHERE id = ?", profileID)
	return err
}

func FormatProfileList(profiles []Profile) string {
	var sb strings.Builder
	sb.WriteString("=== ANEF PROFILES & APPLICANT VAULTS ===\n\n")

	for _, p := range profiles {
		marker := "  "
		if p.Active {
			marker = "★ "
		}
		sb.WriteString(fmt.Sprintf("%sProfile #%d: %s (Created: %s)\n",
			marker, p.ID, p.Name, p.CreatedAt.Format("2006-01-02")))
	}

	return sb.String()
}
