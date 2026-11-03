package evaluator

import (
	"testing"
)

func TestObjJSON(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{10, `{{ {}.json() }}`, "{}"},
		{20, `{{ {one: {two: {}}}.json() }}`, "{&#34;one&#34;:{&#34;two&#34;:{}}}"},
		{
			30,
			`{!! {name: "Chiori", game: "Genshin Impact"}.json() !!}`,
			`{"game":"Genshin Impact","name":"Chiori"}`,
		},
		{
			40,
			`{!! user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() !!}`,
			`{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`,
		},
		{50, `{!! {a: {b: {c: {d: 1}}}}.json() !!}`, `{"a":{"b":{"c":{"d":1}}}}`},
		{
			60,
			`{!! {nums: [1, 2, 3], strs: ['a', 'b']}.json() !!}`,
			`{"nums":[1,2,3],"strs":["a","b"]}`,
		},
		{
			70,
			`{!! {quote: 'He said Hello', newline: 'A B'}.json() !!}`,
			`{"newline":"A B","quote":"He said Hello"}`,
		},
		{
			80,
			`{!! {active: true, count: nil, rate: 3.14}.json() !!}`,
			`{"active":true,"count":null,"rate":3.14}`,
		},
		{90, `{!! {z: 1, a: 2, m: 3}.json() !!}`, `{"a":2,"m":3,"z":1}`},
		{
			100,
			`{!! {user: {name: 'John', age: 30, hobbies: ['coding', 'gaming']}, active: true}.json() !!}`,
			`{"active":true,"user":{"age":30,"hobbies":["coding","gaming"],"name":"John"}}`,
		},
		{
			130,
			`{{ json = {desc: "Here"}.json(); }}<filter :prop='{{ "%s".format(json) }}' />`,
			`<filter :prop='{&#34;desc&#34;:&#34;Here&#34;}' />`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.id)
	}
}

func TestObjCamel(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
	}{
		{
			630,
			`{{ {first_name: 1, LastName: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			640,
			`{!! {HTTP: "https://", NAME: "Serhii"}.camel() !!}`,
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
			`{!! {name_: 1, value_: 2}.camel().json() !!}`,
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
			`{!! {user_data: {first_name: "John", last_name: "Doe"}}.camel() !!}`,
			`{userData: {firstName: "John", lastName: "Doe"}}`,
		},
		{
			801,
			`{!! obj = {first_name: "John"}; res = obj.camel(); obj.json() !!}`,
			`{"first_name":"John"}`,
		},
		{
			802,
			`{!! obj = {user_data: {first_name: "John"}}; res = obj.camel(); obj.json() !!}`,
			`{"user_data":{"first_name":"John"}}`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.id)
	}
}

type testAddress struct {
	Street string
	City   string
}

type testUser struct {
	Name    string
	Age     int
	Address testAddress
	Tags    []string
	Meta    map[string]any
	Manager *testUser
	secret  string
}

func TestObjGet(t *testing.T) {
	cases := []struct {
		id     uint
		inp    string
		expect string
		data   map[string]any
	}{
		{810, `{{ {name: "Chiori"}.get('name') }}`, "Chiori", nil},
		{820, `{{ {}.get('name') }}`, "", nil},
		{821, `{{ {}.get('') }}`, "", nil},
		{830, `{{ {address: {street: "Hunan Road"}}.get('address.street') }}`, "Hunan Road", nil},
		{840, `{{ {address: nil}.get('address.street') }}`, "", nil},
		{850, `{{ {address: {street: "Pretty Road"}}.get('address.street.house') }}`, "", nil},
		{860, `{{ {"1st": 1}.get('1st') }}`, "1", nil},
		{890, `{{ {name: "test"}.get('.') }}`, "", nil},
		{900, `{{ {name: "test"}.get('..') }}`, "", nil},
		{910, `{{ {name: "test"}.get('...') }}`, "", nil},
		{920, `{{ {name: "test"}.get('.name') }}`, "", nil},
		{930, `{{ {name: "test"}.get('name.') }}`, "", nil},
		{940, `{{ {a: {b: 1}}.get('a..b') }}`, "", nil},
		{950, `{{ {}.get('a.b.c.d.e') }}`, "", nil},
		{960, `{{ {a: 1}.get(' ') }}`, "", nil},
		{970, `{{ {a: 1}.get('a ') }}`, "", nil},
		{980, `{{ {a: 1}.get(' a') }}`, "", nil},
		{990, `{{ {a: 1}.get('a .b') }}`, "", nil},
		{1000, `{{ {"na$me": 1}.get('na$me') }}`, "1", nil},
		{1010, `{{ {"na@me": 1}.get('na@me') }}`, "1", nil},
		{1020, `{{ {"": "empty"}.get('') }}`, "empty", nil},
		// Unicode test cases
		{1030, `{{ {"naïve": 1}.get('naïve') }}`, "1", nil},
		{1040, `{{ {"naïve": {value: 2}}.get('naïve.value') }}`, "2", nil},
		{1050, `{{ {"日本語": "Japanese"}.get('日本語') }}`, "Japanese", nil},
		{1060, `{{ {"中文": {nested: "Chinese"}}.get('中文.nested') }}`, "Chinese", nil},
		{1070, `{{ {"emoji🎉": "party"}.get('emoji🎉') }}`, "party", nil},
		{
			1080,
			`{{ {"café": {"résumé": {"naïve": "test"}}}.get('café.résumé.naïve') }}`,
			"test",
			nil,
		},
		{1090, `{{ {"München": "Munich"}.get('München') }}`, "Munich", nil},
		{1100, `{{ {"Ñoño": "child"}.get('Ñoño') }}`, "child", nil},
		// Long paths (10+ levels deep)
		{
			1110,
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: "deep"}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j') }}`,
			"deep",
			nil,
		},
		{
			1120,
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: {k: {l: {m: "very_deep"}}}}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j.k.l.m') }}`,
			"very_deep",
			nil,
		},
		{
			1130,
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: "target"}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j.k') }}`,
			"",
			nil,
		},
		// Dot key precedence and fallback
		{1140, `{{ {"a.b": "direct", a: {b: "nested"}}.get('a.b') }}`, "direct", nil},
		{1150, `{{ {a: {b: "path_value"}}.get('a.b') }}`, "path_value", nil},
		{1160, `{{ {"a.b.c": "literal", a: {b: {c: "nested"}}}.get('a.b.c') }}`, "literal", nil},
		{1170, `{{ {a: {b: {c: "deep_path"}}}.get('a.b.c') }}`, "deep_path", nil},
		{1180, `{{ {"a.b": "exists"}.get('a.b') }}`, "exists", nil},
		{1190, `{{ {}.get('a.b') }}`, "", nil},
		// Edge cases with dots in keys
		{
			1200,
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('.a') }}`,
			"leading_dot",
			nil,
		},
		{
			1210,
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('a.') }}`,
			"trailing_dot",
			nil,
		},
		{
			1220,
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('.') }}`,
			"just_dot",
			nil,
		},
		{1230, `{{ {"a..b": "double_dot", a: {b: "single"}}.get('a..b') }}`, "double_dot", nil},
		{1240, `{{ {"": {"b": "empty_key"}}.get('.b') }}`, "empty_key", nil},
		// Mixed scenarios
		{
			1250,
			`{{ {"x.y": "literal_xy", x: {y: "path_xy", z: "path_xz"}}.get('x.y') }}`,
			"literal_xy",
			nil,
		},
		{
			1260,
			`{{ {"x.y": "literal_xy", x: {y: "path_xy", z: "path_xz"}}.get('x.z') }}`,
			"path_xz",
			nil,
		},
		{
			1270,
			`{{ {"a.b": "ab", "a.b.c": "abc", a: {b: {c: "nested"}}}.get('a.b.c') }}`,
			"abc",
			nil,
		},
		{1280, `{{ {"a.b": "ab", a: {b: "path_b", c: "path_c"}}.get('a.c') }}`, "path_c", nil},
		// With structs
		{
			1290,
			`{{ user.get('Name') }}`,
			"Chiori",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{
			1300,
			`{{ user.get('name') }}`,
			"Chiori",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{
			1310,
			`{{ user.get('age') }}`,
			"30",
			map[string]any{"user": testUser{Age: 30}},
		},
		{
			1320,
			`{{ user.get('address.street') }}`,
			"Hunan Road",
			map[string]any{"user": testUser{Address: testAddress{Street: "Hunan Road"}}},
		},
		{
			1330,
			`{{ user.get('Address.City') }}`,
			"Sumeru",
			map[string]any{"user": testUser{Address: testAddress{City: "Sumeru"}}},
		},
		{
			1340,
			`{{ user.get('aDDress.city') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{City: "Liyue"}}},
		},
		{
			1350,
			`{{ user.get('address.zip') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{City: "Liyue"}}},
		},
		{
			1360,
			`{{ user.get('address.street.house') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{Street: "Hunan Road"}}},
		},
		{
			1370,
			`{{ user.get('missing') }}`,
			"",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{
			1380,
			`{{ user.get('secret') }}`,
			"",
			map[string]any{"user": testUser{secret: "hidden"}},
		},
		{
			1390,
			`{{ user.get('tags') }}`,
			"a, b",
			map[string]any{"user": testUser{Tags: []string{"a", "b"}}},
		},
		{
			1400,
			`{{ user.get('meta.key') }}`,
			"value",
			map[string]any{"user": testUser{Meta: map[string]any{"key": "value"}}},
		},
		{
			1410,
			`{{ user.get('manager.name') }}`,
			"",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{
			1420,
			`{{ user.get('manager.name') }}`,
			"Paimon",
			map[string]any{"user": testUser{Manager: &testUser{Name: "Paimon"}}},
		},
		{
			1430,
			`{{ user.get('id') }}`,
			"",
			map[string]any{"user": struct{ ID int }{ID: 7}},
		},
		{
			1440,
			`{{ user.get('ID') }}`,
			"7",
			map[string]any{"user": struct{ ID int }{ID: 7}},
		},
		{
			1450,
			`{{ user.get('name') }}`,
			"Barbara",
			map[string]any{"user": &testUser{Name: "Barbara"}},
		},
		// Capitalized key fallback
		{1460, `{{ {"Name": "Ann"}.get('name') }}`, "Ann", nil},
		{1470, `{{ {"name": "lower", "Name": "upper"}.get('name') }}`, "lower", nil},
		{1480, `{{ {"Name": "upper"}.get('Name') }}`, "upper", nil},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.data, tc.id)
	}
}
