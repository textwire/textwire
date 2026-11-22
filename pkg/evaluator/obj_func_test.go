package evaluator

import (
	"testing"
)

func TestObjJSON(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{"json of empty object", `{{ {}.json() }}`, "{}"},
		{
			"json escapes quotes in object",
			`{{ {one: {two: {}}}.json() }}`,
			"{&#34;one&#34;:{&#34;two&#34;:{}}}",
		},
		{
			"json of object with strings",
			`{!! {name: "Chiori", game: "Genshin Impact"}.json() !!}`,
			`{"game":"Genshin Impact","name":"Chiori"}`,
		},
		{
			"json of variable object with address",
			`{!! user = {address: {street: "Via Emilio Morosini", city: "Rome"}}; user.json() !!}`,
			`{"address":{"city":"Rome","street":"Via Emilio Morosini"}}`,
		},
		{
			"json of deeply nested object",
			`{!! {a: {b: {c: {d: 1}}}}.json() !!}`,
			`{"a":{"b":{"c":{"d":1}}}}`,
		},
		{
			"json of object with arrays",
			`{!! {nums: [1, 2, 3], strs: ['a', 'b']}.json() !!}`,
			`{"nums":[1,2,3],"strs":["a","b"]}`,
		},
		{
			"json of object with quote and newline",
			`{!! {quote: 'He said Hello', newline: 'A B'}.json() !!}`,
			`{"newline":"A B","quote":"He said Hello"}`,
		},
		{
			"json of mixed bool nil and float",
			`{!! {active: true, count: nil, rate: 3.14}.json() !!}`,
			`{"active":true,"count":null,"rate":3.14}`,
		},
		{
			"json sorts keys alphabetically",
			`{!! {z: 1, a: 2, m: 3}.json() !!}`,
			`{"a":2,"m":3,"z":1}`,
		},
		{
			"json of nested user object",
			`{!! {user: {name: 'John', age: 30, hobbies: ['coding', 'gaming']}, active: true}.json() !!}`,
			`{"active":true,"user":{"age":30,"hobbies":["coding","gaming"],"name":"John"}}`,
		},
		{
			"json embedded in html attribute",
			`{{ json = {desc: "Here"}.json(); }}<filter :prop='{{ "%s".format(json) }}' />`,
			`<filter :prop='{&#34;desc&#34;:&#34;Here&#34;}' />`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestObjCamel(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{"snake case and capitalized key", `{{ {first_name: 1, LastName: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			"uppercase keys",
			`{!! {HTTP: "https://", NAME: "Serhii"}.camel() !!}`,
			`{http: "https://", name: "Serhii"}`,
		},
		{
			"nested object keys",
			`{{ {First: 1, Second: {first_name: 1, LastName: 2}}.camel() }}`,
			`{first: 1, second: {firstName: 1, lastName: 2}}`,
		},
		{"empty object", `{{ {}.camel() }}`, `{}`},
		{"already camel case keys", `{{ {FirstName: 1, LastName: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			"digits in snake case keys",
			`{{ {user_1_name: 1, item2_count: 2}.camel() }}`,
			`{item2Count: 2, user1Name: 1}`,
		},
		{
			"repeated underscores",
			`{{ {first__name: 1, last___name: 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			"keys with spaces",
			`{{ {"first name": 1, "last name": 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{
			"keys with hyphens",
			`{{ {"first-name": 1, "last-name": 2}.camel() }}`,
			`{firstName: 1, lastName: 2}`,
		},
		{"leading underscores", `{{ {_name: 1, _private: 2}.camel() }}`,
			`{Name: 1, Private: 2}`,
		},
		{"trailing underscores", `{{ {name_: 1, value_: 2}.camel() }}`,
			`{name: 1, value: 2}`,
		},
		{
			"trailing underscores then json",
			`{!! {name_: 1, value_: 2}.camel().json() !!}`,
			`{"name":1,"value":2}`,
		},
		{
			"mixed acronym keys",
			`{{ {my_HTTP_request: 1, URL_path: 2}.camel() }}`,
			`{myHttpRequest: 1, urlPath: 2}`,
		},
		{
			"mixed separators",
			`{{ {"first_name-last": 1, "a_b_c": 2}.camel() }}`,
			`{aBC: 2, firstNameLast: 1}`,
		},
		{"uppercase snake case keys", `{{ {API_KEY: 1, USER_ID: 2}.camel() }}`,
			`{apiKey: 1, userId: 2}`,
		},
		{"snake case keys", `{{ {api_key: 1, user_id: 2}.camel() }}`,
			`{apiKey: 1, userId: 2}`,
		},
		{
			"nested snake case object",
			`{!! {user_data: {first_name: "John", last_name: "Doe"}}.camel() !!}`,
			`{userData: {firstName: "John", lastName: "Doe"}}`,
		},
		{
			"camel does not mutate original",
			`{!! obj = {first_name: "John"}; res = obj.camel(); obj.json() !!}`,
			`{"first_name":"John"}`,
		},
		{
			"camel of nested object does not mutate original",
			`{!! obj = {user_data: {first_name: "John"}}; res = obj.camel(); obj.json() !!}`,
			`{"user_data":{"first_name":"John"}}`,
		},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, nil, tc.name)
	}
}

func TestObjCamelKeyCollisions(t *testing.T) {
	cases := []struct {
		name   string
		inp    string
		expect string
	}{
		{
			name:   "key already in camel case wins over converted key",
			inp:    `{{ {a_b: 1, aB: 2}.camel() }}`,
			expect: `{aB: 2}`,
		},
		{
			name:   "key already in camel case wins over key first in byte order",
			inp:    `{{ {"a b": 1, a_b: 2, aB: 3}.camel() }}`,
			expect: `{aB: 3}`,
		},
		{
			name:   "key first in byte order wins when none is in camel case",
			inp:    `{{ {user_name: 1, "user-name": 2, "user name": 3}.camel() }}`,
			expect: `{userName: 3}`,
		},
		{
			name:   "hyphen comes before underscore in byte order",
			inp:    `{{ {_name: 1, "-name": 2}.camel() }}`,
			expect: `{Name: 2}`,
		},
		{
			name:   "keys in nested object collide by the same rule",
			inp:    `{{ {data: {a_b: 1, aB: 2, "c d": 3, c_d: 4}}.camel() }}`,
			expect: `{data: {aB: 2, cD: 3}}`,
		},
		{
			name:   "winning nested object is converted and losing one is dropped",
			inp:    `{{ {user_data: {first_name: 1}, userData: {last_name: 2}}.camel() }}`,
			expect: `{userData: {lastName: 2}}`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			evaluationExpected(t, tc.inp, tc.expect, nil, "")
		})
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
		name   string
		inp    string
		expect string
		data   map[string]any
	}{
		{"get simple key", `{{ {name: "Chiori"}.get('name') }}`, "Chiori", nil},
		{"get missing key from empty object", `{{ {}.get('name') }}`, "", nil},
		{"get empty key from empty object", `{{ {}.get('') }}`, "", nil},
		{
			"get nested dot path",
			`{{ {address: {street: "Hunan Road"}}.get('address.street') }}`,
			"Hunan Road",
			nil,
		},
		{
			"get path through nil value",
			`{{ {address: nil}.get('address.street') }}`,
			"",
			nil,
		},
		{
			"get path beyond nested value",
			`{{ {address: {street: "Pretty Road"}}.get('address.street.house') }}`,
			"",
			nil,
		},
		{"get numeric string key", `{{ {"1st": 1}.get('1st') }}`, "1", nil},
		{"get single dot", `{{ {name: "test"}.get('.') }}`, "", nil},
		{"get double dot", `{{ {name: "test"}.get('..') }}`, "", nil},
		{"get triple dot", `{{ {name: "test"}.get('...') }}`, "", nil},
		{"get leading dot", `{{ {name: "test"}.get('.name') }}`, "", nil},
		{"get trailing dot", `{{ {name: "test"}.get('name.') }}`, "", nil},
		{"get double dot in path", `{{ {a: {b: 1}}.get('a..b') }}`, "", nil},
		{"get deep path in empty object", `{{ {}.get('a.b.c.d.e') }}`, "", nil},
		{"get whitespace key", `{{ {a: 1}.get(' ') }}`, "", nil},
		{"get key with trailing space", `{{ {a: 1}.get('a ') }}`, "", nil},
		{"get key with leading space", `{{ {a: 1}.get(' a') }}`, "", nil},
		{"get key with space before dot", `{{ {a: 1}.get('a .b') }}`, "", nil},
		{"get key with dollar sign", `{{ {"na$me": 1}.get('na$me') }}`, "1", nil},
		{"get key with at sign", `{{ {"na@me": 1}.get('na@me') }}`, "1", nil},
		{"get empty string key", `{{ {"": "empty"}.get('') }}`, "empty", nil},
		// Unicode test cases
		{"get unicode key", `{{ {"naïve": 1}.get('naïve') }}`, "1", nil},
		{
			"get nested unicode path",
			`{{ {"naïve": {value: 2}}.get('naïve.value') }}`,
			"2",
			nil,
		},
		{
			"get japanese key",
			`{{ {"日本語": "Japanese"}.get('日本語') }}`,
			"Japanese",
			nil,
		},
		{
			"get nested chinese key",
			`{{ {"中文": {nested: "Chinese"}}.get('中文.nested') }}`,
			"Chinese",
			nil,
		},
		{
			"get emoji key",
			`{{ {"emoji🎉": "party"}.get('emoji🎉') }}`,
			"party",
			nil,
		},
		{
			"get nested accented path",
			`{{ {"café": {"résumé": {"naïve": "test"}}}.get('café.résumé.naïve') }}`,
			"test",
			nil,
		},
		{
			"get german umlaut key",
			`{{ {"München": "Munich"}.get('München') }}`,
			"Munich",
			nil,
		},
		{"get spanish tilde key", `{{ {"Ñoño": "child"}.get('Ñoño') }}`, "child", nil},
		// Long paths (10+ levels deep)
		{
			"get ten level path",
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: "deep"}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j') }}`,
			"deep",
			nil,
		},
		{
			"get thirteen level path",
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: {k: {l: {m: "very_deep"}}}}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j.k.l.m') }}`,
			"very_deep",
			nil,
		},
		{
			"get eleven level path missing last key",
			`{{ {a: {b: {c: {d: {e: {f: {g: {h: {i: {j: "target"}}}}}}}}}}.get('a.b.c.d.e.f.g.h.i.j.k') }}`,
			"",
			nil,
		},
		// Dot key precedence and fallback
		{
			"get dotted literal key takes precedence",
			`{{ {"a.b": "direct", a: {b: "nested"}}.get('a.b') }}`,
			"direct",
			nil,
		},
		{
			"get nested path",
			`{{ {a: {b: "path_value"}}.get('a.b') }}`,
			"path_value",
			nil,
		},
		{
			"get dotted literal key with three segments",
			`{{ {"a.b.c": "literal", a: {b: {c: "nested"}}}.get('a.b.c') }}`,
			"literal",
			nil,
		},
		{
			"get nested path with three segments",
			`{{ {a: {b: {c: "deep_path"}}}.get('a.b.c') }}`,
			"deep_path",
			nil,
		},
		{"get dotted literal key from object", `{{ {"a.b": "exists"}.get('a.b') }}`, "exists", nil},
		{"get dotted path from empty object", `{{ {}.get('a.b') }}`, "", nil},
		// Edge cases with dots in keys
		{
			"get leading dot key",
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('.a') }}`,
			"leading_dot",
			nil,
		},
		{
			"get trailing dot key",
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('a.') }}`,
			"trailing_dot",
			nil,
		},
		{
			"get just dot key",
			`{{ {".a": "leading_dot", "a.": "trailing_dot", ".": "just_dot"}.get('.') }}`,
			"just_dot",
			nil,
		},
		{
			"get double dot literal key",
			`{{ {"a..b": "double_dot", a: {b: "single"}}.get('a..b') }}`,
			"double_dot",
			nil,
		},
		{
			"get path starting with empty key",
			`{{ {"": {"b": "empty_key"}}.get('.b') }}`,
			"empty_key",
			nil,
		},
		// Mixed scenarios
		{
			"get dotted literal over nested path",
			`{{ {"x.y": "literal_xy", x: {y: "path_xy", z: "path_xz"}}.get('x.y') }}`,
			"literal_xy",
			nil,
		},
		{
			"get nested sibling path",
			`{{ {"x.y": "literal_xy", x: {y: "path_xy", z: "path_xz"}}.get('x.z') }}`,
			"path_xz",
			nil,
		},
		{
			"get longest dotted literal key",
			`{{ {"a.b": "ab", "a.b.c": "abc", a: {b: {c: "nested"}}}.get('a.b.c') }}`,
			"abc",
			nil,
		},
		{
			"get nested sibling when literal does not match",
			`{{ {"a.b": "ab", a: {b: "path_b", c: "path_c"}}.get('a.c') }}`,
			"path_c",
			nil,
		},
		// With structs
		{"get struct exported field", `{{ user.get('Name') }}`,
			"Chiori",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{"get struct field lowercase", `{{ user.get('name') }}`,
			"Chiori",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{"get struct int field", `{{ user.get('age') }}`,
			"30",
			map[string]any{"user": testUser{Age: 30}},
		},
		{"get struct nested struct field", `{{ user.get('address.street') }}`,
			"Hunan Road",
			map[string]any{"user": testUser{Address: testAddress{Street: "Hunan Road"}}},
		},
		{"get struct nested field with capitalized path", `{{ user.get('Address.City') }}`,
			"Sumeru",
			map[string]any{"user": testUser{Address: testAddress{City: "Sumeru"}}},
		},
		{"get struct field with wrong casing", `{{ user.get('aDDress.city') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{City: "Liyue"}}},
		},
		{"get missing struct field", `{{ user.get('address.zip') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{City: "Liyue"}}},
		},
		{"get missing nested struct field", `{{ user.get('address.street.house') }}`,
			"",
			map[string]any{"user": testUser{Address: testAddress{Street: "Hunan Road"}}},
		},
		{"get missing struct field name", `{{ user.get('missing') }}`,
			"",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{"get unexported struct field", `{{ user.get('secret') }}`,
			"",
			map[string]any{"user": testUser{secret: "hidden"}},
		},
		{"get struct slice field", `{{ user.get('tags') }}`,
			"a, b",
			map[string]any{"user": testUser{Tags: []string{"a", "b"}}},
		},
		{"get struct map field", `{{ user.get('meta.key') }}`,
			"value",
			map[string]any{"user": testUser{Meta: map[string]any{"key": "value"}}},
		},
		{"get nil struct pointer field", `{{ user.get('manager.name') }}`,
			"",
			map[string]any{"user": testUser{Name: "Chiori"}},
		},
		{"get nested struct pointer field", `{{ user.get('manager.name') }}`,
			"Paimon",
			map[string]any{"user": testUser{Manager: &testUser{Name: "Paimon"}}},
		},
		{"get missing struct int field lowercase", `{{ user.get('id') }}`,
			"",
			map[string]any{"user": struct{ ID int }{ID: 7}},
		},
		{"get struct int field uppercase", `{{ user.get('ID') }}`,
			"7",
			map[string]any{"user": struct{ ID int }{ID: 7}},
		},
		{"get struct pointer field", `{{ user.get('name') }}`,
			"Barbara",
			map[string]any{"user": &testUser{Name: "Barbara"}},
		},
		// Capitalized key fallback
		{"get capitalized key fallback", `{{ {"Name": "Ann"}.get('name') }}`, "Ann", nil},
		{
			"get exact lowercase key when both cases exist",
			`{{ {"name": "lower", "Name": "upper"}.get('name') }}`,
			"lower",
			nil,
		},
		{"get exact capitalized key", `{{ {"Name": "upper"}.get('Name') }}`, "upper", nil},
	}

	for _, tc := range cases {
		evaluationExpected(t, tc.inp, tc.expect, tc.data, tc.name)
	}
}
