package lsp

import (
	"embed"
	"fmt"
	"log"
	"sync"

	"github.com/textwire/textwire/v2/token"
)

// Locale represents a language locale for metadata.
// It's a 2 letter ISO 639-1 code (e.g. "en", "es", "fr").
// Codes: https://en.wikipedia.org/wiki/List_of_ISO_639_language_codes
type Locale string

//go:embed meta/*
var files embed.FS

// tokenMetaCache caches metadata for tokens by locale.
var tokenMetaCache = make(map[Locale]map[token.TokenType]string)

// cacheMutex ensures thread-safe access to tokenMetaCache.
var cacheMutex sync.RWMutex

// GetTokenMeta returns a hover description for the given token type.
// If no description is found, an empty string is returned.
func GetTokenMeta(tok token.TokenType, locale Locale) string {
	metaMap, ok := getMetaFromCache(locale)

	if !ok {
		metaMap = loadMetaFromFile(locale)
		saveMeta(locale, metaMap)
	}

	// Return the metadata for the token.
	meta, ok := metaMap[tok]
	if !ok {
		log.Printf("no metadata found for token: %v in locale: %s", tok, locale)
		return ""
	}

	return meta
}

func getMetaFromCache(locale Locale) (map[token.TokenType]string, bool) {
	cacheMutex.RLock()
	metaMap, ok := tokenMetaCache[locale]
	cacheMutex.RUnlock()

	return metaMap, ok
}

// loadMetaFromFile loads metadata for all tokens in the given locale.
func loadMetaFromFile(locale Locale) map[token.TokenType]string {
	metaMap := make(map[token.TokenType]string)

	// Define the list of tokens and their corresponding file names.
	tokens := []struct {
		tok  token.TokenType
		file string
	}{
		{token.IF, "if.md"},
		{token.ELSE_IF, "ifelse.md"},
		// Add more tokens here.
	}

	// Load metadata for each token.
	for _, t := range tokens {
		filePath := fmt.Sprintf("meta/%s/%s", locale, t.file)
		data, err := files.ReadFile(filePath)
		if err != nil {
			log.Printf("failed to load metadata for token %v in locale %s: %v", t.tok, locale, err)
			continue
		}
		metaMap[t.tok] = string(data)
	}

	return metaMap
}

func saveMeta(locale Locale, metaMap map[token.TokenType]string) {
	cacheMutex.Lock()
	tokenMetaCache[locale] = metaMap
	cacheMutex.Unlock()
}
