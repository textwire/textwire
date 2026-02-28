package evaluator

import "testing"

func TestEvalBinaryFunctions(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		// binary
		{10, `{{ true.binary() }}`, "1"},
		{20, `{{ false.binary() }}`, "0"},
		{30, `{{ true.binary().float() }}`, "1.0"},
		{40, `{{ false.binary().float() }}`, "0.0"},
		// then
		{50, `{{ true.then("nice") }}`, "nice"},
		{60, `{{ false.then("nice") }}`, ""},
		{70, `{{ true.then(1) }}`, "1"},
		{80, `{{ true.then([1, 2]) }}`, "1, 2"},
		{90, `{{ true.then(true) }}`, "1"},
		{100, `{{ true.then(false) }}`, "0"},
		{110, `{{ true.then(1.121) }}`, "1.121"},
		{120, `{{ false.then("first", "second") }}`, "second"},
		{130, `{{ false.then(1, 2) }}`, "2"},
		{140, `{{ false.then([1, 2], [3, 4]) }}`, "3, 4"},
		{150, `{{ false.then(true, false) }}`, "0"},
		{160, `{{ false.then(false, true) }}`, "1"},
		{170, `{{ false.then(1.121, 4.2141) }}`, "4.2141"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
