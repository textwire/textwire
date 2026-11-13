package evaluator

import (
	"slices"
	"testing"

	"codeberg.org/textwire/textwire/v5/pkg/value"
)

func TestEvalArrFunctions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		// len
		{"len of empty array", `{{ [].len() }}`, "0"},
		{"len of three elements", `{{ [1, 2, 3].len() }}`, "3"},
		{"len counts only top-level elements", `{{ [0, [2, [1, 2]]].len() }}`, "2"},
		// join
		{"join ints with comma", `{{ [1, 2, 3].join(", ") }}`, "1, 2, 3"},
		{
			"join strings with space",
			`{{ ['Chiori', 'Venti', 'Raiden'].join(" ") }}`,
			"Chiori Venti Raiden",
		},
		{
			"join strings with default comma",
			`{{ ['Chiori', 'Venti', 'Raiden'].join() }}`,
			"Chiori,Venti,Raiden",
		},
		{"join empty array", `{{ [].join() }}`, ""},
		{
			"join mixed types with dash",
			`{{ [1, 2.5, true, 'text'].join('-') }}`,
			"1-2.5-1-text",
		},
		{"join with nil element", `{{ [1, nil, 3].join(', ') }}`, "1, , 3"},
		{"join nested arrays", `{{ [[1, 2], [3, 4]].join('|') }}`, "1, 2|3, 4"},
		// rand
		{"rand on empty array", `{{ [].rand() }}`, ""},
		{"rand on single int", `{{ [1].rand() }}`, "1"},
		{"rand on single string", `{{ ['Chiori'].rand() }}`, "Chiori"},
		{"rand unwraps nested arrays", `{{ [[[4]]].rand().rand().rand() }}`, "4"},
		// reverse
		{"reverse three ints", `{{ [1, 2, 3].reverse() }}`, "3, 2, 1"},
		{"reverse single string", `{{ ['Mona'].reverse() }}`, "Mona"},
		{"reverse empty array", `{{ [].reverse() }}`, ""},
		{
			"reverse three strings",
			`{{ ['Raiden', 'Venti', 'Barbara'].reverse() }}`,
			"Barbara, Venti, Raiden",
		},
		{"reverse nested array", `{{ [4, 3, [1, 2]].reverse() }}`, "1, 2, 3, 4"},
		{"reverse does not mutate original", `{{ a = [1, 2, 3]; b = a.reverse(); a }}`, "1, 2, 3"},
		// slice
		{"slice empty array from zero", `{{ [].slice(0) }}`, ""},
		{"slice full array from zero", `{{ [1, 2, 3].slice(0) }}`, "1, 2, 3"},
		{
			"slice with large negative start clamps to zero",
			`{{ [1, 2, 3].slice(-34) }}`,
			"1, 2, 3",
		}, // should change -35 to 0
		{"slice from index two", `{{ [0, 1, 2, 3].slice(2) }}`, "2, 3"},
		{"slice start beyond length", `{{ [0, 1, 2, 3].slice(5) }}`, ""},
		{"slice from one to three", `{{ [0, 1, 2, 3].slice(1, 3) }}`, "1, 2"},
		{
			"slice strings one to two",
			`{{ ['Chiori', 'Venti', 'Raiden', "Nahida"].slice(1, 2) }}`,
			"Venti",
		},
		{
			"slice with negative end clamps to length",
			`{{ ['Kokomi', 'Venti', 'Raiden', "Nahida"].slice(0, -3) }}`,
			"Kokomi, Venti, Raiden, Nahida",
		},
		{
			"slice with both negative bounds clamps",
			`{{ [1, 2, 3, 4].slice(-3, -1) }}`,
			"1, 2, 3, 4",
		},
		{"slice with equal start and end", `{{ [1, 2, 3].slice(2, 2) }}`, ""},
		{"slice negative start to positive end", `{{ [1, 2, 3, 4].slice(-2, 2) }}`, "1, 2"},
		{"slice start greater than end", `{{ [1, 2, 3].slice(2, 1) }}`, ""},
		{"slice start beyond length with end one", `{{ [1, 2, 3].slice(3, 1) }}`, ""},
		{
			"slice returns independent copy",
			`{{ a = [1, 2, 3]; b = a.slice(0); b[0] = 9; a }}`,
			"1, 2, 3",
		},
		{
			"slice returns independent subarray copy",
			`{{ a = [1, 2, 3]; b = a.slice(1, 2); b[0] = 9; a }}`,
			"1, 2, 3",
		},
		// shuffle
		{"shuffle empty array", `{{ [].shuffle() }}`, ""},
		{"shuffle single int", `{{ [1].shuffle() }}`, "1"},
		{"shuffle single string", `{{ ['Cryo'].shuffle() }}`, "Cryo"},
		// contains
		{"contains in empty array", `{{ [].contains(1) }}`, "0"},
		{"contains first element", `{{ [1, 2, 3].contains(1) }}`, "1"},
		{"contains middle element", `{{ [1, 2, 3].contains(2) }}`, "1"},
		{"contains last element", `{{ [1, 2, 3].contains(3) }}`, "1"},
		{"contains missing element", `{{ [1, 2, 3].contains(4) }}`, "0"},
		{"contains int in mixed array", `{{ [{}, [1, 2], 3].contains(1) }}`, "0"},
		{"contains empty object", `{{ [{}, [1, 2], 3].contains({}) }}`, "1"},
		{
			"contains object with all pairs matching",
			`{{ [{nice: 3, cool: 'anna'}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`,
			"1",
		},
		{
			"contains object regardless of pair order",
			`{{ [{cool: 'anna', nice: 3}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`,
			"1",
		},
		{"contains empty array", `{{ [[], 3].contains([]) }}`, "1"},
		{"contains array", `{{ [[1, 2], 3].contains([1, 2]) }}`, "1"},
		{"contains array with different order", `{{ [[2, 1], 3].contains([1, 2]) }}`, "0"},
		{"contains object in int array", `{{ [1, 2].contains({}) }}`, "0"},
		{"contains object without matching pairs", `{{ [{}, 21].contains({age: 21}) }}`, "0"},
		{"contains nested array", `{{ [[], [1], [2]].contains([2]) }}`, "1"},
		{"negated contains missing object", `{{ ![{}, 21].contains({age: 21}) }}`, "1"},
		{"negated contains present array", `{{ ![[], [1], [2]].contains([2]) }}`, "0"},
		{"contains nil element", `{{ [1, nil, 3].contains(nil) }}`, "1"},
		{"contains nil from missing property", `{{ [{}.name].contains(nil) }}`, "1"},
		{"contains duplicate int", `{{ [1, 0].contains(1) }}`, "1"},
		{"contains float", `{{ [1.5, 2.5].contains(1.5) }}`, "1"},
		{"contains string in mixed array", `{{ [1, 'Venti', 3].contains('Venti') }}`, "1"},
		{"contains nested array structure", `{{ [[1, [2, 3]]].contains([1, [2, 3]]) }}`, "1"},
		{
			"contains matches raw array element to plain array",
			`{{ html = '<br>'; [[raw(html)]].contains([html]) }}`,
			"1",
		},
		{
			"contains matches raw object element to plain object",
			`{{ obj = { html: '<br>' }; [raw(obj)].contains(obj) }}`,
			"1",
		},
		{
			"contains in raw array finds plain string",
			`{{ raw(['<a>', '<b>', '<c>']).contains('<b>') }}`,
			"1",
		},
		// append
		{"append to empty array", `{{ [].append(1) }}`, "1"},
		{"append int", `{{ [1, 2, 3].append(4) }}`, "1, 2, 3, 4"},
		{
			"append string",
			`{{ ['Chiori', 'Venti'].append('Raiden') }}`,
			"Chiori, Venti, Raiden",
		},
		{"append multiple arguments", `{{ [1, 2].append(3, 4) }}`, "1, 2, 3, 4"},
		{"append array flattens", `{{ [1, 2].append([3, 4]) }}`, "1, 2, 3, 4"},
		{"append array then len", `{{ [1, 2].append([3, 4]).len() }}`, "3"},
		{
			"append twice with reassignment",
			`{{ a = []; a = a.append(1); a = a.append(2); a }}`,
			"1, 2",
		},
		{
			"append arrays with reassignment",
			`{{ a = []; a = a.append([1, 2]); a = a.append([3, 4]); a }}`,
			"1, 2, 3, 4",
		},
		// prepend
		{"prepend to empty array", `{{ [].prepend(1) }}`, "1"},
		{"prepend int", `{{ [1, 2, 3].prepend(4) }}`, "4, 1, 2, 3"},
		{
			"prepend string",
			`{{ ['Chiori', 'Venti'].prepend('Raiden') }}`,
			"Raiden, Chiori, Venti",
		},
		{"prepend multiple arguments", `{{ [1, 2].prepend(3, 4) }}`, "3, 4, 1, 2"},
		{"prepend array flattens", `{{ [1, 2].prepend([3, 4]) }}`, "3, 4, 1, 2"},
		{"prepend array then len", `{{ [1, 2].prepend([3, 4]).len() }}`, "3"},
		// json
		{"json of empty array", `{{ [].json() }}`, "[]"},
		{
			"json of mixed values",
			`{!! [1, 2.1, true, false, nil, "Chiori", []].json() !!}`,
			`[1,2.1,true,false,null,"Chiori",[]]`,
		},
		{
			"json of object and negative int",
			`{!! [{name: "Skirk", game: "Genshin Impact"}, -10].json() !!}`,
			`[{"game":"Genshin Impact","name":"Skirk"},-10]`,
		},
		{"json of deeply nested array", `{{ [[[[[1,2]]]]].json() }}`, "[[[[[1,2]]]]]"},
		{
			"json of array of objects",
			`{!! [{name: 'Venti'}, {name: 'Escoffier'}].json() !!}`,
			`[{"name":"Venti"},{"name":"Escoffier"}]`,
		},
		{"json of booleans and nil", `{{ [true, false, nil].json() }}`, "[true,false,null]"},
		{
			"json of nested object in array",
			`{!! [{a: [1, {b: 2}]}].json() !!}`,
			`[{"a":[1,{"b":2}]}]`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestArrRandFunc_CanSelectAnyElement(t *testing.T) {
	const calls = 100

	elems := []value.Literal{
		&value.Int{Val: 1},
		&value.Int{Val: 2},
		&value.Int{Val: 3},
	}
	arr := &value.Arr{Elements: elems}

	seen := make(map[value.Literal]bool)
	for range calls {
		got, err := arrRandFunc(arr)
		if err != nil {
			t.Fatal(err)
		}
		if !slices.Contains(elems, got) {
			t.Fatalf("rand returned an element outside the array: %v", got)
		}
		seen[got] = true
	}

	for i, elem := range elems {
		if !seen[elem] {
			t.Errorf("element %d was never returned in %d calls", i, calls)
		}
	}
}
