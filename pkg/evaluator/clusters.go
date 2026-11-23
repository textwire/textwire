package evaluator

import (
	"unicode"
	"unicode/utf8"

	"github.com/clipperhouse/uax29/v2/graphemes"
)

const (
	zeroWidthNonJoiner = '\u200C'
	zeroWidthJoiner    = '\u200D'
	thaiSaraAm         = '\u0E33'
	laoSaraAm          = '\u0EB3'

	emojiModifierFirst = 0x1F3FB
	emojiModifierLast  = 0x1F3FF
	tagCharacterFirst  = 0xE0020
	tagCharacterLast   = 0xE007F
)

func isLeadingClusterExtender(r rune) bool {
	switch r {
	case zeroWidthNonJoiner, zeroWidthJoiner, thaiSaraAm, laoSaraAm:
		return true
	}

	return unicode.IsMark(r) ||
		isInRange(r, emojiModifierFirst, emojiModifierLast) ||
		isInRange(r, tagCharacterFirst, tagCharacterLast)
}

func isInRange(r, first, last rune) bool {
	return r >= first && r <= last
}

func leadingClusterExtenderPrefixLen(s string) int {
	end := 0
	for end < len(s) {
		r, size := utf8.DecodeRuneInString(s[end:])
		if !isLeadingClusterExtender(r) {
			break
		}
		end += size
	}

	return end
}

func splitClusters(s string) []string {
	var clusters []string

	iterator := graphemes.FromString(s)
	for iterator.Next() {
		clusters = append(clusters, iterator.Value())
	}

	return clusters
}
