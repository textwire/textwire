package textwire

import (
	"testing"
)

func BenchmarkArrayJoinFunc(b *testing.B) {
	arr := make([]string, 10000)
	code := "{{ arr.join(' ') }}"

	b.ResetTimer()

	for b.Loop() {
		EvaluateString(code, map[string]any{"arr": arr})
	}
}

func BenchmarkArrayAppendFunc(b *testing.B) {
	arr := make([]struct{}, 10000)
	o1 := struct{}{}
	o2 := struct{}{}
	o3 := struct{}{}
	o4 := struct{}{}
	code := "{{ arr.append(new) }}"

	b.ResetTimer()

	for b.Loop() {
		EvaluateString(code, map[string]any{
			"arr": arr,
			"o1":  o1,
			"o2":  o2,
			"o3":  o3,
			"o4":  o4,
		})
	}
}
