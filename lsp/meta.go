package lsp

import (
	"embed"
	"fmt"
	"path"
	"sync"

	"slices"

	"github.com/textwire/textwire/v2/token"
)

// Locale represents a language locale for metadata.
// It's a 2 letter ISO 639-1 code (e.g. "en", "es", "fr").
// Codes: https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes
type Locale string

var (
	//go:embed metadata/*
	files embed.FS

	fileNamesOnce sync.Once
	fileNames     map[token.TokenType]string

	validLocales = []Locale{"en"}
)

// GetTokenMeta retrieves metadata for the given token type and locale.
func GetTokenMeta(tok token.TokenType, locale Locale) (string, error) {
	if !isValidLocale(locale) {
		return "", fmt.Errorf("invalid locale: %s", locale)
	}

	fileNamesOnce.Do(initFileNames)

	fileName, ok := fileNames[tok]
	if !ok {
		return "", fmt.Errorf("no metadata found for token: %v", tok)
	}

	return loadMeta(locale, fileName)
}

func GetFiles() embed.FS {
	return files
}

func initFileNames() {
	for dir, tok := range token.GetDirectives() {
		fileNames[tok] = dir[1:] + ".md"
	}
}

// loadMeta loads metadata from the embedded files for the given locale and file name.
func loadMeta(locale Locale, fileName string) (string, error) {
	filePath := path.Join("metadata", string(locale), fileName)

	data, err := files.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	return string(data), nil
}

func isValidLocale(locale Locale) bool {
	return slices.Contains(validLocales, locale)
}
