package auth

import (
	"net/http"

	"github.com/arihant-dev/anef-tracker/pkg/session"
)

type CurlSession = session.Session

func SaveSession(sess *session.Session) error {
	return session.SaveSession(sess)
}

func LoadSession() (*session.Session, error) {
	return session.LoadSession()
}

func ParseCurl(curlCmd string) (*session.Session, error) {
	return session.ParseCurl(curlCmd)
}

func InjectAuthHeaders(req *http.Request, sess *session.Session) {
	session.InjectAuthHeaders(req, sess)
}

func AuthenticateViaBrowser(portalURL string) (*session.Session, error) {
	return session.AuthenticateViaBrowser(portalURL)
}

func AuthenticateWithCredentials(httpClient *http.Client, username, password string) (*session.Session, error) {
	return session.AuthenticateWithCredentials(httpClient, username, password)
}
