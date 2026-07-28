package session

import (
	"github.com/arihant-dev/anef-tracker/pkg/importer"
)

func Tokenize(input string) []string {
	return importer.Tokenize(input)
}
