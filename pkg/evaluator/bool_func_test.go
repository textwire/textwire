package evaluator

import "testing"

func TestEvalBinaryFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// binary
		{"true binary", `{{ true.binary() }}`, "1"},
		{"false binary", `{{ false.binary() }}`, "0"},
		{"true binary to float", `{{ true.binary().float() }}`, "1.0"},
		{"false binary to float", `{{ false.binary().float() }}`, "0.0"},
		// then
		{"true then string", `{{ true.then("nice") }}`, "nice"},
		{"false then string", `{{ false.then("nice") }}`, ""},
		{"true then int", `{{ true.then(1) }}`, "1"},
		{"true then array", `{{ true.then([1, 2]) }}`, "1, 2"},
		{"true then true", `{{ true.then(true) }}`, "1"},
		{"true then false", `{{ true.then(false) }}`, "0"},
		{"true then float", `{{ true.then(1.121) }}`, "1.121"},
		{"false then picks second string", `{{ false.then("first", "second") }}`, "second"},
		{"false then picks second int", `{{ false.then(1, 2) }}`, "2"},
		{"false then picks second array", `{{ false.then([1, 2], [3, 4]) }}`, "3, 4"},
		{"false then with true false picks false", `{{ false.then(true, false) }}`, "0"},
		{"false then with false true picks true", `{{ false.then(false, true) }}`, "1"},
		{"false then picks second float", `{{ false.then(1.121, 4.2141) }}`, "4.2141"},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}
