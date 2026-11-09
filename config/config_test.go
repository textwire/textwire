package config

import "testing"

func TestConfigure_MaxComponentDepth(t *testing.T) {
	cases := []struct {
		name   string
		value  int
		expect int
	}{
		{name: "zero keeps the default", value: 0, expect: DefaultMaxComponentDepth},
		{name: "negative keeps the default", value: -1, expect: DefaultMaxComponentDepth},
		{name: "positive overrides the default", value: 50, expect: 50},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			conf := New()
			conf.Configure(&Config{TemplateDir: t.TempDir(), MaxComponentDepth: tc.value})

			if conf.MaxComponentDepth != tc.expect {
				t.Fatalf("expected %d, got %d", tc.expect, conf.MaxComponentDepth)
			}
		})
	}
}
