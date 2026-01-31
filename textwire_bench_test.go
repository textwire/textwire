package textwire

import (
	"testing"
)

func BenchmarkArrayJoinFunc(b *testing.B) {
	largeArr := make([]string, 10000)
	code := "{{ arr.join(' ') }}"

	b.ResetTimer()

	for b.Loop() {
		EvaluateString(code, map[string]any{
			"arr": largeArr,
		})
	}
}
