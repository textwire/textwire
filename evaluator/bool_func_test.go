package evaluator

import "testing"

func TestEvalBinaryFunctions(t *testing.T) {
	cases := []struct {
		inp      string
		expected string
	}{
		// binary
		{`{{ true.binary() }}`, "1"},
		{`{{ false.binary() }}`, "0"},
		{`{{ !true.binary() }}`, "0"},
		{`{{ !false.binary() }}`, "1"},
		{`{{ true.binary().float() }}`, "1.0"},
		{`{{ false.binary().float() }}`, "0.0"},
		// then
		{`{{ true.then("nice") }}`, "nice"},
		{`{{ false.then("nice") }}`, ""},
		{`{{ true.then(1) }}`, "1"},
		{`{{ true.then([1, 2]) }}`, "1, 2"},
		{`{{ true.then(true) }}`, "1"},
		{`{{ true.then(false) }}`, "0"},
		{`{{ true.then(1.121) }}`, "1.121"},
		{`{{ false.then("first", "second") }}`, "second"},
		{`{{ false.then(1, 2) }}`, "2"},
		{`{{ false.then([1, 2], [3, 4]) }}`, "3, 4"},
		{`{{ false.then(true, false) }}`, "0"},
		{`{{ false.then(false, true) }}`, "1"},
		{`{{ false.then(1.121, 4.2141) }}`, "4.2141"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expected)
	}
}
