package evaluator

import (
	"testing"
)

func TestEvalObjectFunctions(t *testing.T) {
	cases := []struct {
		id     int
		inp    string
		expect string
	}{
		// json
		{10, `{{ {}.json() }}`, "{}"},
		{20, `{{ {one: {two: {}}}.json() }}`, `{"one":{"two":{}}}`},
		{
			30,
			`{{ {name: "Chiori", game: "Genshin Impact"}.json() }}`,
			`{"game":"Genshin Impact","name":"Chiori"}`,
		},
		{
			40,
			`{{ user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() }}`,
			`{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`,
		},
		{50, `{{ {a: {b: {c: {d: 1}}}}.json() }}`, `{"a":{"b":{"c":{"d":1}}}}`},
		{
			60,
			`{{ {nums: [1, 2, 3], strs: ['a', 'b']}.json() }}`,
			`{"nums":[1,2,3],"strs":["a","b"]}`,
		},
		{
			70,
			`{{ {quote: 'He said Hello', newline: 'A B'}.json() }}`,
			`{"newline":"A B","quote":"He said Hello"}`,
		},
		{
			80,
			`{{ {active: true, count: nil, rate: 3.14}.json() }}`,
			`{"active":true,"count":null,"rate":3.14}`,
		},
		{90, `{{ {z: 1, a: 2, m: 3}.json() }}`, `{"a":2,"m":3,"z":1}`},
		{
			100,
			`{{ {user: {name: 'John', age: 30, hobbies: ['coding', 'gaming']}, active: true}.json() }}`,
			`{"active":true,"user":{"age":30,"hobbies":["coding","gaming"],"name":"John"}}`,
		},
		{618, `{{ {value: (1.0/0.0)}.json() }}`, `{"value":null}`},
		{
			619,
			`{{ {nan: (0.0/0.0), inf: (1.0/0.0), ninf: (-1.0/0.0)}.json() }}`,
			`{"inf":null,"nan":null,"ninf":null}`,
		},
		// camel
		{
			630,
			`{{ {first_name: 1, LastName: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			640,
			`{{ {HTTP: "https://", NAME: "Serhii"}.camel() }}`,
			`{http: "https://", name: "Serhii"}`,
		},
		{
			650,
			`{{ {First: 1, Second: {first_name: 1, LastName: 2}}.camel() }}`,
			`{first: 1, second: {firstName: 1, lastName: 2}}`,
		},
		{660, `{{ {}.camel() }}`, `{}`},
		{
			670,
			`{{ {FirstName: 1, LastName: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			680,
			`{{ {user_1_name: 1, item2_count: 2}.camel() }}`,
			`{item2Count: 2, user1Name: 1}`,
		},
		{
			690,
			`{{ {first__name: 1, last___name: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			700,
			`{{ {"first name": 1, "last name": 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			710,
			`{{ {"first-name": 1, "last-name": 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			730,
			`{{ {_name: 1, _private: 2}.camel() }}`,
			`{Name: 1, Private: 2}`,
		},
		{
			740,
			`{{ {name_: 1, value_: 2}.camel() }}`,
			`{name: 1, value: 2}`,
		},
		{
			750,
			`{{ {name_: 1, value_: 2}.camel().json() }}`,
			`{"name":1,"value":2}`,
		},
		{
			760,
			`{{ {my_HTTP_request: 1, URL_path: 2}.camel() }}`,
			`{myHttpRequest: 1, urlPath: 2}`,
		},
		{
			770,
			`{{ {"first_name-last": 1, "a_b_c": 2}.camel() }}`,
			`{aBC: 2, firstNameLast: 1}`,
		},
		{
			780,
			`{{ {API_KEY: 1, USER_ID: 2}.camel() }}`,
			`{apiKey: 1, userId: 2}`,
		},
		{
			790,
			`{{ {api_key: 1, user_id: 2}.camel() }}`,
			`{apiKey: 1, userId: 2}`,
		},
		{
			800,
			`{{ {user_data: {first_name: "John", last_name: "Doe"}}.camel() }}`,
			`{userData: {firstName: "John", lastName: "Doe"}}`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.id)
	}
}
