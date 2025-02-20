package lsp

import (
	"embed"
	"errors"
	"fmt"
	"sync"

	"github.com/textwire/textwire/v2/token"
)

// Locale represents a language locale for metadata.
// It's a 2 letter ISO 639-1 code (e.g. "en", "es", "fr").
// Codes: https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes
type Locale string

var (
	NoMetadataError = errors.New("no metadata found for token")
	FailToLoadMeta  = errors.New("failed to load metadata for a given file name")

	//go:embed metadata/*
	files embed.FS

	// cacheMutex ensures thread-safe access to tokenMetaCache.
	cacheMutex sync.RWMutex
)

var fileNames = map[token.TokenType]string{
	token.IF:      "if.md",
	token.ELSE_IF: "ifelse.md",
}

// GetTokenMeta returns a hover description for the given token type.
// If no description is found, an empty string is returned.
func GetTokenMeta(tok token.TokenType, locale Locale) ([]byte, error) {
	fileName, ok := fileNames[tok]
	if !ok {
		return []byte{}, NoMetadataError
	}

	meta, err := loadMeta(locale, fileName)
	if err != nil {
		return []byte{}, NoMetadataError
	}

	return meta, nil
}

// loadMeta loads metadata for a given file name and locale.
func loadMeta(locale Locale, fileName string) ([]byte, error) {
	filePath := fmt.Sprintf("meta/%s/%s", locale, fileName)

	data, err := files.ReadFile(filePath)
	if err != nil {
		return []byte{}, FailToLoadMeta
	}

	return data, nil
}
