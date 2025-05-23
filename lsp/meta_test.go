package lsp

import (
	"strings"
	"testing"

	"github.com/textwire/textwire/v2/token"
)

func TestGetTokenMeta(t *testing.T) {
	t.Run("Invalid locale", func(t *testing.T) {
		_, err := GetTokenMeta(token.IF, "invalid")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	t.Run("Token missing docs", func(t *testing.T) {
		_, err := GetTokenMeta(token.EOF, "en")
		if err == nil {
			t.Error("expected error, got nil")
		}
	})

	testCases := []struct {
		name   string
		token  token.TokenType
		locale Locale
		expect string
	}{
		{"@if token", token.IF, "en", "@if(condition)"},
		{"@elseif token", token.ELSE_IF, "en", "@elseif(condition2)"},
		{"@each token", token.EACH, "en", "@each(item in items)"},
		{"@for token", token.FOR, "en", "@for(i = 0; i < items.len(); i++)"},
		{"@else token", token.ELSE, "en", "@else"},
		{"@dump token", token.DUMP, "en", "@dump(variable)"},
		{"@use token", token.USE, "en", "@use('layoutName')"},
		{"@insert token", token.INSERT, "en", "@insert('reservedName')"},
		{"@reserve token", token.RESERVE, "en", "@reserve('reservedName')"},
		{"@component token", token.COMPONENT, "en", "@component('path/to', { prop })"},
		{"@slot token", token.SLOT, "en", "@slot('name')"},
		{"@end token", token.END, "en", "@end"},
		{"@break token", token.BREAK, "en", "@break"},
		{"@continue token", token.CONTINUE, "en", "@continue"},
		{"@breakIf token", token.BREAK_IF, "en", "@breakIf(condition)"},
		{"@continueIf token", token.CONTINUE_IF, "en", "@continueIf(condition)"},
	}

	for _, tc := range testCases {
		t.Run("Valid "+tc.name, func(t *testing.T) {
			meta, err := GetTokenMeta(tc.token, tc.locale)
			if err != nil {
				t.Errorf("expected err to be nil, got error %v", err)
			}

			if !strings.Contains(meta, tc.expect) {
				t.Errorf("expected %s in meta, got %s", tc.expect, meta)
			}
		})
	}
}
