package importer_test

import (
	"testing"

	"github.com/arihant-dev/anef-tracker/pkg/importer"
)

func TestCorpus_ChromeMacOS(t *testing.T) {
	curlCmd := `curl 'https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour' \
  -H 'authority: administration-etrangers-en-france.interieur.gouv.fr' \
  -H 'accept: application/json, text/plain, */*' \
  -H 'authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijk5OTk5OTk5OTkiLCJleHAiOjIwMDAwMDAwMDB9.sig' \
  -H 'cookie: auth_token=token123; refresh_token=ref123' \
  --compressed`

	imp := importer.NewCurlImporter()
	sess, err := imp.Import([]byte(curlCmd))
	if err != nil {
		t.Fatalf("Chrome MacOS import failed: %v", err)
	}

	if sess.User != "9999999999" {
		t.Errorf("expected user 9999999999, got %s", sess.User)
	}
	if !sess.Capabilities.CanFetch {
		t.Errorf("expected CanFetch true")
	}
}

func TestCorpus_FirefoxLinux(t *testing.T) {
	curlCmd := `curl "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour" \
  -H "Host: administration-etrangers-en-france.interieur.gouv.fr" \
  -H "User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/115.0" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijg4ODg4ODg4ODgiLCJleHAiOjIwMDAwMDAwMDB9.sig" \
  -H "Cookie: auth_token=token456"`

	imp := importer.NewCurlImporter()
	sess, err := imp.Import([]byte(curlCmd))
	if err != nil {
		t.Fatalf("Firefox Linux import failed: %v", err)
	}

	if sess.User != "8888888888" {
		t.Errorf("expected user 8888888888, got %s", sess.User)
	}
}

func TestCorpus_PowerShellWindows(t *testing.T) {
	curlCmd := `curl "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour" ` + "`" + `
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6Ijc3Nzc3Nzc3NzciLCJleHAiOjIwMDAwMDAwMDB9.sig" ` + "`" + `
  -H "Cookie: auth_token=token789"`

	imp := importer.NewCurlImporter()
	sess, err := imp.Import([]byte(curlCmd))
	if err != nil {
		t.Fatalf("PowerShell Windows import failed: %v", err)
	}

	if sess.User != "7777777777" {
		t.Errorf("expected user 7777777777, got %s", sess.User)
	}
}

func TestCorpus_CmdWindows(t *testing.T) {
	curlCmd := `curl "https://administration-etrangers-en-france.interieur.gouv.fr/api/sejour/usager/demande_sejour" ^
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJsb2dpbiI6IjY2NjY2NjY2NjYiLCJleHAiOjIwMDAwMDAwMDB9.sig" ^
  -H "Cookie: auth_token=token101"`

	imp := importer.NewCurlImporter()
	sess, err := imp.Import([]byte(curlCmd))
	if err != nil {
		t.Fatalf("Cmd Windows import failed: %v", err)
	}

	if sess.User != "6666666666" {
		t.Errorf("expected user 6666666666, got %s", sess.User)
	}
}
