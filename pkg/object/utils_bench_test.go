package object

import (
	"testing"
)

func BenchmarkNativeSliceToArrayObject(b *testing.B) {
	cases := []struct {
		name string
		size int
	}{
		{"small", 100},
		{"medium", 1000},
		{"large", 10_000},
		{"huge", 100_000},
	}

	for _, tc := range cases {
		b.Run(tc.name, func(b *testing.B) {
			slice := make([]any, tc.size)
			for i := 0; i < tc.size; i++ {
				switch i % 4 {
				case 0:
					slice[i] = "string-" + string(rune('a'+i%26))
				case 1:
					slice[i] = i * 42
				case 2:
					slice[i] = float64(i) * 1.5
				case 3:
					slice[i] = i%2 == 0
				}
			}

			b.ResetTimer()

			for b.Loop() {
				_ = nativeSliceToArrayObject(slice)
			}
		})
	}
}
