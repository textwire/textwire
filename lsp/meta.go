package lsp

import (
	"embed"
	"path"
	"strings"
	"sync"

	"slices"

	"github.com/textwire/textwire/v3/lsp/utils"
	"github.com/textwire/textwire/v3/token"
)

// Locale represents a language locale for metadata.
// It's a 2 letter ISO 639-1 code (e.g. "en", "es", "fr").
// Codes: https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes
type Locale string

var (
	//go:embed metadata/*
	files embed.FS

	//go:embed inserts/*
	insertFiles embed.FS

	fileNamesOnce sync.Once
	fileNames     map[token.TokenType]string

	validLocales = []Locale{"en"}
)

// GetTokenMeta retrieves metadata for the given token type and locale.
func GetTokenMeta(tok token.TokenType, locale Locale) (string, error) {
	if !isValidLocale(locale) {
		return "", utils.ErrInvalidLocale(string(locale))
	}

	fileNamesOnce.Do(initFileNames)

	fileName, ok := fileNames[tok]
	if !ok {
		return "", utils.ErrNoMetadataFound(tok)
	}

	filePath := path.Join("metadata", string(locale), fileName)

	data, err := files.ReadFile(filePath)
	if err != nil {
		return "", utils.FailedToReadFile("meta", filePath, err)
	}

	return string(data), nil
}

// GetTokenInsert retrieves insert string for the given token type. This
// insert is used for autocompletion.
func GetTokenInsert(tok token.TokenType) (string, error) {
	fileNamesOnce.Do(initFileNames)

	fileName, ok := fileNames[tok]
	if !ok {
		return "", utils.ErrNoMetadataFound(tok)
	}

	filePath := path.Join("inserts", fileName)

	data, err := insertFiles.ReadFile(filePath)
	if err != nil {
		return "", utils.FailedToReadFile("insert", filePath, err)
	}

	return string(data), nil
}

func initFileNames() {
	fileNames = map[token.TokenType]string{}

	for dir, tok := range token.GetDirectives() {
		name := strings.ToLower(dir[1:])
		fileNames[tok] = name + ".md"
	}
}

func isValidLocale(locale Locale) bool {
	return slices.Contains(validLocales, locale)
}
