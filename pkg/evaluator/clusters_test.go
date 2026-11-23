package evaluator

import (
	"reflect"
	"strings"
	"testing"
)

func TestSplitClusters(t *testing.T) {
	tests := []struct {
		name   string
		inp    string
		expect []string
	}{
		{
			name:   "empty string has no clusters",
			inp:    "",
			expect: nil,
		},
		{
			name:   "single ASCII character forms one cluster",
			inp:    "a",
			expect: []string{"a"},
		},
		{
			name:   "plain ASCII characters form individual clusters",
			inp:    "abc",
			expect: []string{"a", "b", "c"},
		},
		{
			name:   "multibyte runes without marks form individual clusters",
			inp:    "日本語", // 日本語
			expect: []string{"日", "本", "語"},
		},
		{
			name:   "precomposed character forms one cluster",
			inp:    "\u00e9", // é (U+00E9, precomposed)
			expect: []string{"\u00e9"},
		},
		{
			name:   "combining acute accent attaches to its base character",
			inp:    "e\u0301", // e + U+0301 combining acute → é (decomposed)
			expect: []string{"e\u0301"},
		},
		{
			name:   "combining mark attaches to preceding character in text",
			inp:    "ae\u0301b", // a + é(decomposed) + b
			expect: []string{"a", "e\u0301", "b"},
		},
		{
			name:   "multiple combining marks attach to the same base character",
			inp:    "a\u0301\u0302\u0303", // a + acute + circumflex + tilde
			expect: []string{"a\u0301\u0302\u0303"},
		},
		{
			name:   "enclosing mark attaches to its base character",
			inp:    "a\u20dd", // a + U+20DD combining enclosing circle → a⃝
			expect: []string{"a\u20dd"},
		},
		{
			name:   "spacing combining mark attaches to its base character",
			inp:    "\u0915\u093e", // क + ा (Devanagari ka + aa vowel sign) → का
			expect: []string{"\u0915\u093e"},
		},
		{
			name:   "mixed combining mark categories attach to their base character",
			inp:    "a\u0301\u20dd\u093e", // a + acute + enclosing circle + Devanagari aa
			expect: []string{"a\u0301\u20dd\u093e"},
		},
		{
			name:   "separate marked characters form separate clusters",
			inp:    "e\u0301a\u0300o\u0308", // é à ö (all decomposed)
			expect: []string{"e\u0301", "a\u0300", "o\u0308"},
		},
		{
			name:   "variation selector attaches to its base character",
			inp:    "\u2764\ufe0f", // ❤️ (red heart + VS16)
			expect: []string{"\u2764\ufe0f"},
		},
		{
			name:   "skin tone modifier attaches to its emoji",
			inp:    "\U0001F44D\U0001F3FD", // 👍🏽 (thumbs up + medium skin tone)
			expect: []string{"\U0001F44D\U0001F3FD"},
		},
		{
			name:   "regional indicators pair to form flags",
			inp:    "\U0001F1E8\U0001F1F3\U0001F1F5\U0001F1F1", // 🇨🇳🇵🇱
			expect: []string{"\U0001F1E8\U0001F1F3", "\U0001F1F5\U0001F1F1"},
		},
		{
			name:   "unpaired regional indicator forms its own cluster",
			inp:    "\U0001F1E8\U0001F1F3\U0001F1F5", // 🇨🇳 + unpaired regional indicator P
			expect: []string{"\U0001F1E8\U0001F1F3", "\U0001F1F5"},
		},
		{
			name:   "emoji family joined with ZWJ forms one cluster",
			inp:    "\U0001F468\u200D\U0001F469\u200D\U0001F467", // 👨‍👩‍👧
			expect: []string{"\U0001F468\u200D\U0001F469\u200D\U0001F467"},
		},
		{
			name:   "rainbow flag joined with ZWJ forms one cluster",
			inp:    "\U0001F3F3\uFE0F\u200D\U0001F308", // 🏳️‍🌈 (white flag + VS16 + ZWJ + rainbow)
			expect: []string{"\U0001F3F3\uFE0F\u200D\U0001F308"},
		},
		{
			name:   "emoji ZWJ cluster stays separate from following text",
			inp:    "\U0001F468\u200D\U0001F469\u200D\U0001F467 hi", // 👨‍👩‍👧 hi
			expect: []string{"\U0001F468\u200D\U0001F469\u200D\U0001F467", " ", "h", "i"},
		},
		{
			name:   "CRLF forms one cluster between characters",
			inp:    "a\r\nb", // a + CRLF + b
			expect: []string{"a", "\r\n", "b"},
		},
		{
			// 🏴󠁧󠁢󠁥󠁮󠁧󠁿 (black flag + tag letters gbeng + cancel tag) = England flag
			name: "subdivision flag tags stay with the black flag",
			inp:  "\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F",
			expect: []string{
				"\U0001F3F4\U000E0067\U000E0062\U000E0065\U000E006E\U000E0067\U000E007F",
			},
		},
		{
			name:   "ZWJ does not join following non-emoji text",
			inp:    "a\u200Db", // a + ZWJ + b
			expect: []string{"a\u200D", "b"},
		},
		{
			name:   "leading ZWJ does not join following digits",
			inp:    "\u200D00", // ZWJ + "00"
			expect: []string{"\u200D", "0", "0"},
		},
		{
			name:   "Thai spacing mark attaches to its base character",
			inp:    "\u0E01\u0E33", // กำ (Thai ko kai + sara am)
			expect: []string{"\u0E01\u0E33"},
		},
		{
			name:   "decomposed Hangul syllable forms one cluster",
			inp:    "\u1112\u1161\u11AB", // 한 (Hangul jamo: h + a + n)
			expect: []string{"\u1112\u1161\u11AB"},
		},
		{
			name:   "leading combining mark forms its own cluster",
			inp:    "\u0301", // bare combining acute accent
			expect: []string{"\u0301"},
		},
		{
			name:   "leading combining mark stays separate from following character",
			inp:    "\u0301a", // combining acute + a
			expect: []string{"\u0301", "a"},
		},
		{
			name:   "consecutive leading combining marks form one cluster",
			inp:    "\u0301\u0302a", // acute + circumflex + a
			expect: []string{"\u0301\u0302", "a"},
		},
		{
			name:   "spaces and punctuation each form individual clusters",
			inp:    "a  b.",
			expect: []string{"a", " ", " ", "b", "."},
		},
		{
			name:   "invalid UTF-8 byte is preserved as its own cluster",
			inp:    "a\xffb", // a + invalid byte 0xFF + b
			expect: []string{"a", "\xff", "b"},
		},
		{
			name:   "combining mark attaches to preceding invalid UTF-8 byte",
			inp:    "\xff\u0301", // invalid byte 0xFF + combining acute
			expect: []string{"\xff\u0301"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := splitClusters(tt.inp)

			if !reflect.DeepEqual(actual, tt.expect) {
				t.Errorf("splitClusters(%q)\n actual: %q\nexpect: %q", tt.inp, actual, tt.expect)
			}

			if reassembled := strings.Join(actual, ""); reassembled != tt.inp {
				t.Errorf("clusters reassemble to %q, expect original input %q", reassembled, tt.inp)
			}
		})
	}
}
