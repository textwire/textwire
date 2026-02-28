package evaluator

import (
	"testing"
)

func TestEvalArrayFunctions(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		// len
		{10, `{{ [].len() }}`, "0"},
		{20, `{{ [1, 2, 3].len() }}`, "3"},
		{30, `{{ [0, [2, [1, 2]]].len() }}`, "2"},
		// join
		{40, `{{ [1, 2, 3].join(", ") }}`, "1, 2, 3"},
		{50, `{{ ['Chiori', 'Venti', 'Raiden'].join(" ") }}`, "Chiori Venti Raiden"},
		{60, `{{ ['Chiori', 'Venti', 'Raiden'].join() }}`, "Chiori,Venti,Raiden"},
		{70, `{{ [].join() }}`, ""},
		{75, `{{ [1, 2.5, true, 'text'].join('-') }}`, "1-2.5-1-text"},
		{76, `{{ [1, nil, 3].join(', ') }}`, "1, , 3"},
		{77, `{{ [[1, 2], [3, 4]].join('|') }}`, "1, 2|3, 4"},
		// rand
		{80, `{{ [].rand() }}`, ""},
		{90, `{{ [1].rand() }}`, "1"},
		{100, `{{ ['Chiori'].rand() }}`, "Chiori"},
		{110, `{{ [[[4]]].rand().rand().rand() }}`, "4"},
		// reverse
		{120, `{{ [1, 2, 3].reverse() }}`, "3, 2, 1"},
		{130, `{{ ['Chiori'].reverse() }}`, "Chiori"},
		{140, `{{ [].reverse() }}`, ""},
		{150, `{{ ['Raiden', 'Venti', 'Chiori'].reverse() }}`, "Chiori, Venti, Raiden"},
		{160, `{{ [4, 3, [1, 2]].reverse() }}`, "1, 2, 3, 4"},
		// slice
		{170, `{{ [].slice(0) }}`, ""},
		{180, `{{ [1, 2, 3].slice(0) }}`, "1, 2, 3"},
		{190, `{{ [1, 2, 3].slice(-34) }}`, "1, 2, 3"}, // should change -35 to 0
		{200, `{{ [0, 1, 2, 3].slice(2) }}`, "2, 3"},
		{210, `{{ [0, 1, 2, 3].slice(5) }}`, ""},
		{220, `{{ [0, 1, 2, 3].slice(1, 3) }}`, "1, 2"},
		{230, `{{ ['Chiori', 'Venti', 'Raiden', "Nahida"].slice(1, 2) }}`, "Venti"},
		{
			240,
			`{{ ['Chiori', 'Venti', 'Raiden', "Nahida"].slice(0, -3) }}`,
			"Chiori, Venti, Raiden, Nahida",
		},
		{250, `{{ [1, 2, 3, 4].slice(-3, -1) }}`, "1, 2, 3, 4"},
		{251, `{{ [1, 2, 3].slice(2, 2) }}`, ""},
		{252, `{{ [1, 2, 3, 4].slice(-2, 2) }}`, "1, 2"},
		// shuffle
		{260, `{{ [].shuffle() }}`, ""},
		{270, `{{ [1].shuffle() }}`, "1"},
		{280, `{{ ['Chiori'].shuffle() }}`, "Chiori"},
		// contains
		{290, `{{ [].contains(1) }}`, "0"},
		{300, `{{ [1, 2, 3].contains(1) }}`, "1"},
		{310, `{{ [1, 2, 3].contains(2) }}`, "1"},
		{320, `{{ [1, 2, 3].contains(3) }}`, "1"},
		{330, `{{ [1, 2, 3].contains(4) }}`, "0"},
		{340, `{{ [{}, [1, 2], 3].contains(1) }}`, "0"},
		{350, `{{ [{}, [1, 2], 3].contains({}) }}`, "1"},
		{360, `{{ [{nice: 3, cool: 'anna'}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`, "1"},
		{370, `{{ [{cool: 'anna', nice: 3}, [1, 2], 3].contains({nice:3, cool: 'anna'}) }}`, "1"},
		{380, `{{ [[], 3].contains([]) }}`, "1"},
		{390, `{{ [[1, 2], 3].contains([1, 2]) }}`, "1"},
		{400, `{{ [[2, 1], 3].contains([1, 2]) }}`, "0"},
		{410, `{{ [1, 2].contains({}) }}`, "0"},
		{420, `{{ [{}, 21].contains({age: 21}) }}`, "0"},
		{430, `{{ [[], [1], [2]].contains([2]) }}`, "1"},
		{440, `{{ ![{}, 21].contains({age: 21}) }}`, "1"},
		{450, `{{ ![[], [1], [2]].contains([2]) }}`, "0"},
		{451, `{{ [1, nil, 3].contains(nil) }}`, "1"},
		{452, `{{ [1, 0].contains(1) }}`, "1"},
		{453, `{{ [1.5, 2.5].contains(1.5) }}`, "1"},
		{454, `{{ [1, 'Venti', 3].contains('Venti') }}`, "1"},
		{455, `{{ [[1, [2, 3]]].contains([1, [2, 3]]) }}`, "1"},
		// append
		{460, `{{ [].append(1) }}`, "1"},
		{470, `{{ [1, 2, 3].append(4) }}`, "1, 2, 3, 4"},
		{480, `{{ ['Chiori', 'Venti'].append('Raiden') }}`, "Chiori, Venti, Raiden"},
		{490, `{{ [1, 2].append(3, 4) }}`, "1, 2, 3, 4"},
		{500, `{{ [1, 2].append([3, 4]) }}`, "1, 2, 3, 4"},
		{510, `{{ [1, 2].append([3, 4]).len() }}`, "3"},
		// prepend
		{520, `{{ [].prepend(1) }}`, "1"},
		{530, `{{ [1, 2, 3].prepend(4) }}`, "4, 1, 2, 3"},
		{540, `{{ ['Chiori', 'Venti'].prepend('Raiden') }}`, "Raiden, Chiori, Venti"},
		{550, `{{ [1, 2].prepend(3, 4) }}`, "3, 4, 1, 2"},
		{560, `{{ [1, 2].prepend([3, 4]) }}`, "3, 4, 1, 2"},
		{570, `{{ [1, 2].prepend([3, 4]).len() }}`, "3"},
		// json
		{580, `{{ [].json() }}`, "[]"},
		{
			590,
			`{{ [1, 2.1, true, false, nil, "Chiori", []].json() }}`,
			`[1,2.1,true,false,null,"Chiori",[]]`,
		},
		{
			600,
			`{{ [{name: "Chiori", game: "Genshin Impact"}, -10].json() }}`,
			`[{"game":"Genshin Impact","name":"Chiori"},-10]`,
		},
		{610, `{{ [[[[[1,2]]]]].json() }}`, "[[[[[1,2]]]]]"},
		{
			611,
			`{{ [{name: 'Venti'}, {name: 'Chiori'}].json() }}`,
			`[{"name":"Venti"},{"name":"Chiori"}]`,
		},
		{612, `{{ [true, false, nil].json() }}`, "[true,false,null]"},
		{613, `{{ [{a: [1, {b: 2}]}].json() }}`, `[{"a":[1,{"b":2}]}]`},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
