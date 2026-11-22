package value

import "testing"

func TestObjJSON(t *testing.T) {
	t.Run("escapes keys", func(t *testing.T) {
		obj := NewObj(map[string]Literal{
			`a"b<c>`: &Str{Val: `<script>`},
		})

		got, err := obj.JSON()
		if err != nil {
			t.Fatal(err)
		}

		expect := `{"a\"b<c>":"<script>"}`
		if got != expect {
			t.Fatalf("expected %q, got %q", expect, got)
		}
	})
}

func TestIsPreferredKey(t *testing.T) {
	cases := []struct {
		name     string
		key      string
		otherKey string
		camelKey string
		expected bool
	}{
		{
			name:     "camel case key vs snake case key",
			key:      "aB",
			otherKey: "a_b",
			camelKey: "aB",
			expected: true,
		},
		{
			name:     "snake case key vs camel case key",
			key:      "a_b",
			otherKey: "aB",
			camelKey: "aB",
			expected: false,
		},
		{
			name:     "camel case key vs key with space",
			key:      "aB",
			otherKey: "a b",
			camelKey: "aB",
			expected: true,
		},
		{
			name:     "key with space vs camel case key",
			key:      "a b",
			otherKey: "aB",
			camelKey: "aB",
			expected: false,
		},
		{
			name:     "none in camel case, first in byte order",
			key:      "user name",
			otherKey: "user_name",
			camelKey: "userName",
			expected: true,
		},
		{
			name:     "none in camel case, last in byte order",
			key:      "user_name",
			otherKey: "user name",
			camelKey: "userName",
			expected: false,
		},
		{
			name:     "none in camel case, numbers, first in byte order",
			key:      "1skirk",
			otherKey: "2skirk",
			camelKey: "skirk",
			expected: true,
		},
		{
			name:     "none in camel case, numbers, last in byte order",
			key:      "2skirk",
			otherKey: "1skirk",
			camelKey: "skirk",
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isPreferredKey(tc.key, tc.otherKey, tc.camelKey)

			if actual != tc.expected {
				t.Fatalf(
					"isPreferredKey(%q, %q, %q) = %v, expected %v",
					tc.key, tc.otherKey, tc.camelKey, actual, tc.expected,
				)
			}
		})
	}
}
