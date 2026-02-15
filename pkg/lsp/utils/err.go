package utils

import (
	"fmt"

	"github.com/textwire/textwire/v3/pkg/token"
)

func ErrInvalidLocale(locale string) error {
	return fmt.Errorf("invalid locale: %s", locale)
}

func ErrNoMetadataFound(tok token.TokenType) error {
	return fmt.Errorf("no metadata found for token: %v", tok)
}

func FailedToReadFile(readTarget, filePath string, err error) error {
	return fmt.Errorf("failed to read %s file %s: %w", readTarget, filePath, err)
}
