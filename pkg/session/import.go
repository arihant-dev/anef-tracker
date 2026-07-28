package session

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// PromptCurlPaste reads a cURL command from standard input in the terminal.
func PromptCurlPaste(r io.Reader) (*Session, error) {
	if r == nil {
		r = os.Stdin
	}

	scanner := bufio.NewScanner(r)
	fmt.Println("\n------------------------------------------------------------")
	fmt.Println("1. Log into ANEF in your browser")
	fmt.Println("2. Open DevTools (F12) -> Network tab")
	fmt.Println("3. Right click any request to ANEF -> Copy as cURL")
	fmt.Println("4. Paste below and press [ENTER] twice:")
	fmt.Println("------------------------------------------------------------")
	fmt.Print("> ")

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" && len(lines) > 0 {
			break
		}
		if trimmed != "" {
			lines = append(lines, line)
		}
	}

	rawInput := strings.Join(lines, "\n")
	if strings.TrimSpace(rawInput) == "" {
		return nil, fmt.Errorf("no input provided")
	}

	sess, err := ParseCurl(rawInput)
	if err != nil {
		return nil, fmt.Errorf("failed parsing cURL: %w", err)
	}

	sess.ImportSource = ImportCurl
	return sess, nil
}

// AuthenticateViaBrowser opens browser to ANEF portal and prompts for cURL paste.
func AuthenticateViaBrowser(portalURL string) (*Session, error) {
	if portalURL == "" {
		portalURL = "https://administration-etrangers-en-france.interieur.gouv.fr/usagers/"
	}

	fmt.Printf("Opening ANEF portal in browser (%s)...\n", portalURL)
	_ = OpenBrowser(portalURL)

	sess, err := PromptCurlPaste(os.Stdin)
	if err != nil {
		return nil, err
	}

	sess.ImportSource = ImportBrowserAssisted
	return sess, nil
}
