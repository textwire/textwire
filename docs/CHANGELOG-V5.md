# Release notes V5

## v5.2.3 (2026-11-06)

- 🐛 Fixed a panic on a `@for` post statement that isn't followed by `)`, like `@for(i = 0; i < 1; x@end`. It now reports error. Closes [#162](https://codeberg.org/textwire/textwire/issues/162).
- 🚀 Improved printing of deeply nested arrays and objects. `String()`, `JSON()`, `Dump()` and `Native()` now check for cycles once per value instead of at every nesting level, so printing scales linearly with depth. Closes [#163](https://codeberg.org/textwire/textwire/issues/163).
    | Case | Depth | Before | After | Speedup |
    |---|---|---|---|---|
    | Nested arrays | 1,000 | 127 ms | 1.21 ms | 105x |
    | Nested arrays | 2,000 | 551 ms | 2.59 ms | 213x |
    | Nested arrays | 4,000 | 2.44 s | 5.80 ms | 420x |
    | Nested objects | 1,000 | 184 ms | 3.60 ms | 51x |
    | Nested objects | 2,000 | 812 ms | 11.4 ms | 71x |
    | Nested objects | 4,000 | 3.64 s | 40.5 ms | 90x |
    | Go `list.json()` | 1,000 | 215 ms | 8.95 ms | 24x |
    | Go `list.json()` | 2,000 | 963 ms | 30.3 ms | 32x |
    | Go `list.json()` | 4,000 | 3.64 s | 84.9 ms | 43x |
    | Go `{{ list }}` | 1,000 | 211 ms | 7.57 ms | 28x |
    | Go `{{ list }}` | 2,000 | 936 ms | 26.6 ms | 35x |
    | Go `{{ list }}` | 4,000 | 3.62 s | 81.2 ms | 45x |

## v5.2.2 (2026-11-05)

- 🐛 Fixed `expected expression after '+'` pointing at the left operand instead of the operator, so multi-line expressions like `{{ 1\n\n\n+ }}` now report the line with `+`.
- 🐛 Fixed `capitalize()` using upper case instead of title case for the first letter. `"ǆemal".capitalize()` now gives `ǅemal` instead of `Ǆemal`.
- 🐛 Fixed `Program.String()` panicking on empty `@if`, `@elseif`, `@for`, `@each` and `@insert` bodies. Empty `@else` blocks now print as `@else` instead of disappearing.
- 🐛 Fixed a parser panic when `@elseif` follows `@else`, like `@if(true)@else@elseif(true)@end`. It now reports `@elseif cannot come after @else`.

## v5.2.1 (2026-11-05)

- 🐛 Fixed parser hanging forever on a broken `@for` init or post statement without `@end`, like `@for(x)`. It now stops at EOF again. Closes [#151](https://codeberg.org/textwire/textwire/issues/151).
- 🐛 Fixed `reverse()` and `truncate()` splitting flags and ZWJ emoji sequences, like `🇨🇳🇵🇱`, `👨‍👩‍👧` and `🏳️‍🌈`, into separate characters. Closes [#152](https://codeberg.org/textwire/textwire/issues/152).
- 🐛 Fixed `len()`, `at()`, `first()` and `last()` splitting accented letters, emoji and flags into runes. They now use the same character boundaries as `reverse()` and `truncate()`, so `len()` counts characters. Closes [#153](https://codeberg.org/textwire/textwire/issues/153).
- 🐛 Added missing `@end` to hover message on `@component`.
- 🐛 Fixed token columns counting UTF-8 bytes instead of characters, which shifted error columns, editor links and loop detection on lines with non-ASCII text. Closes [#154](https://codeberg.org/textwire/textwire/issues/154).
- 🐛 Fixed `@reserve()`, `@insert()` and default `@pass` printing extra quotes, and `@insert()` with a block not printing its block. Closes [#155](https://codeberg.org/textwire/textwire/issues/155).

## v5.2.0 (2026-11-04)

- ✨ Added support for object property assignment through index notation. Participant [@peterlodri-sec](https://codeberg.org/peterlodri-sec). Closes [#147](https://codeberg.org/textwire/textwire/issues/147) and [#132](https://codeberg.org/textwire/textwire/issues/132).

## v5.1.20 (2026-11-03)

- 🐛 Fixed `illegal token` errors showing invisible characters as a plain space. Non-breaking space and other invisible characters, like NUL, are now shown by code point, e.g. `U+00A0`. Invalid UTF-8 bytes are shown as `0xFF` instead of `�`. Closes [#145](https://codeberg.org/textwire/textwire/issues/145).
- 🐛 Fixed the invalid assignment error printing `got 'error'` for non-literal left sides. It now names the expression, e.g. `got 'prefix expression'` for `-x = 2` or `got 'infix expression'` for `x + 1 = 2`. Closes [#143](https://codeberg.org/textwire/textwire/issues/143).
- 🐛 Fixed `decimal()` crashing on huge decimal counts, like `1.decimal(".", 9223372036854775807)` or `"12".decimal(".", 2000000000)`. It now fails with `result is too large`, the same limit as `repeat()`. Closes [#140](https://codeberg.org/textwire/textwire/issues/140) and [#142](https://codeberg.org/textwire/textwire/issues/142).
- 🐛 Fixed `abs()`, `int()`, `ceil()`, `floor()` and `round()` returning wrong numbers outside the int64 range. `(-9223372036854775808).abs()` and `9223372036854775807.float().int()` now fail with `result is out of int64 range` instead of silently wrapping. Closes [#144](https://codeberg.org/textwire/textwire/issues/144).
- 🚀  Improved nested `@if`, `@for` and `@each` evaluation speed. Deeply nested templates no longer take quadratic time. Evaluation of nested `@if` blocks now scales linearly. Closes [#141](https://codeberg.org/textwire/textwire/issues/141).
    | Nested `@if` levels | Before | After | Speedup |
    |---|---|---|---|
    | 1,000 | 30.9 ms | 5.9 ms | 5.2x |
    | 5,000 | 344 ms | 26.8 ms | 12.8x |
    | 10,000 | 1,182 ms | 57.4 ms | 20.6x |
    | 20,000 | 4,466 ms | 99.9 ms | 44.7x |
    | 40,000 | 17,231 ms | 212 ms | 81.2x |

## v5.1.19 (2026-11-02)

- 🐛 Fixed `@for()` reporting a bad init as a bad post statement. Init failures now report `expects statement as init conditional`.
- 🐛 Fixed `reverse()` and `truncate()` splitting characters made of several code points. Accents written as separate code points and emoji with skin-tone modifiers now stay together, so `"e\u0301x".truncate(1)` gives `é...`. Invalid UTF-8 is preserved.
- 🐛 Fixed `capitalize()` replacing an invalid first byte with U+FFFD. It now returns the string unchanged.
- 🐛 Fixed template errors from file paths or messages containing `%` being mangled, like `/tmp/100%tw` becoming `/tmp/100%!t(MISSING)w`.
- 🐛 Fixed global data with unsupported Go types, like channels, being silently stored. `AddGlobal` now returns the error.

## v5.1.18 (2026-11-01)

- 🧑‍💻 Improved nested ternary in the true branch logic. Before, it needed parentheses, now `a ? b ? "x" : "y" : "z"` parses without parentheses.
- 🐛 Fixed `decimal()` on strings: it now applies the requested separator to values that already have enough decimals, keeps large integers exact, and leaves non-numeric strings like `nan` unchanged.
- 🐛 Fixed unclosed `@if`, `@each`, `@for` and `@component` directives. The first three reported `illegal token ''` on line 1, while `@component` silently parsed. They now report `missing @end for @if started on line N`.
- 🐛 Fixed `slice()` returning an array that shared its backing array with the original, so changing the slice silently changed the source array.
- 🐛 Fixed self-referencing arrays and objects crashing the process with a stack overflow or out of memory. Assignments that would create a circular reference now fail with `Cannot assign a to a[0]: assignment would create a circular reference.` String, JSON, dump and native conversion also detect cycles, and cyclic Go maps, slices and pointers are rejected during conversion. Example is `{{ a = [1]; a[0] = a; a }}` or `{{ b = {x: 1}; b.y = b; b }}`.

## v5.1.17 (2026-10-29)

- 🐛 Fixed file watcher serving an edited template before its layout and components were linked.
- 🐛 Fixed single `&` and `|` being silently treated as `+`. They now produce a syntax error instead.
- 🐛 Fixed `reverse()` mutating the original array. It now returns a reversed copy.
- 🐛 Fixed `get()` on objects not finding struct fields that dot access finds.
- 🐛 Fixed `@each` and `@for` loop variables leaking out of the loop and clashing with or overwriting outer variables. Loop variables now live in their own scope.

## v5.1.15 (2026-10-28)

- ♻️ Refactored file watcher. Added unit tests to it.

## v5.1.14 (2026-10-27)

- 🐛 Fixed `Template.Response` sending error pages with HTTP 200. It now writes 500 for render errors and 404 for missing templates.
- 🐛 Fixed `locateFiles` treating any path containing `.tw` as a template. It now only loads files ending in the template extension, and no longer mangles names in directories like `my.tweaks/`.
- 🐛 Fixed file watcher mutating programs in place while renders use them. It now parses and links fresh programs and swaps the whole slice under the write lock.

## v5.1.13 (2026-10-27)

- 🐛 Fixed native maps with non-string keys collapsing into one object entry. Keys now use their formatted values, so `map[int]string{1: "one", 2: "two"}` retains both entries.
- 🐛 Fixed `@each` failing when array elements have different types. The loop variable is now redefined for each element instead of being type-checked as an assignment.
- 🐛 Fixed `.capitalize()` corrupting strings that start with a multi-byte character. It now uppercases the first Unicode rune without splitting its UTF-8 encoding.
- 🐛 Fixed array `.rand()` always returning the first element. It now selects a random element.
- 🐛 Fixed float decrement failing for negative values, underflowing values between 0 and 1, and mutating aliased values. Decrement now assigns a new float value.
- 🐛 Fixed absolute `TemplateDir` paths losing their leading slash. Template directories are now cleaned without turning absolute paths into relative ones.
- 🐛 Fixed `!` rejecting booleans from data, logical expressions and functions. It now negates boolean values regardless of their source.
- 🐛 Fixed a panic when a `@for` init statement is missing its operator, like `@for(i 0; ...)`. It now reports a syntax error.
- 🐛 Fixed `uint64` values above `math.MaxInt64` rendering as negative numbers. They now report `unsigned integer '...' is out of int64 range` instead of silently wrapping to `-1`.
- 🐛 Fixed `@reserve` and `@slot` rendering fallback or empty inside conditionals and loops. Child evaluation contexts now share directive data, which remains available across loop iterations. For example, with `@insert("title", "Home")`, this now renders `Home` instead of `Untitled`:
  ```textwire
  @if(showTitle)
      @reserve("title", "Untitled")
  @end
  ```

## v5.1.12 (2026-10-26)

- 🐛 Fixed `@dump` printing strings and object keys without HTML escaping, which let dumped data inject markup into the page. They're now escaped with `html.EscapeString`.
- 🐛 Fixed a panic when accessing an object with an empty string key, like `{{ {a: 1}[""] }}`. It now returns `nil` like any other missing key.
- 🐛 Fixed named Go types and fixed-size arrays in template data. `type Role string`, `time.Duration` and `[3]int` now convert instead of failing. Unsupported kinds like `func`, `chan` or `complex` no longer leave a nil value that panicked on access; they report `unsupported value type 'func()' for field 'Fn'` instead.
- 🧑‍💻 Fixed README.md file examples.
- 🧑‍💻 Removed a useless random generator creation from `shuffle()`. It was created and thrown away before the shuffle ran, and Go seeds the global RNG automatically.

## v5.1.11 (2026-10-26)

- 🐛 Fixed saving a template with a syntax error while `FileWatcher` is on crashing the whole app. The watcher now logs the parse error and keeps the previously working program instead of dereferencing a nil one.
- 🐛 Fixed the file watcher calling `os.Exit` when it failed to walk the template folder. It now logs the error and keeps watching.
- 🐛 Fixed a panic when the right operand is missing inside parentheses, like `{{ (1 + }}`. It now reports `expected expression after '+'` instead of crashing the parser, REPL and file watcher.
- 🐛 Fixed `.json()` producing invalid JSON. String values and object keys are now encoded with `encoding/json`, so control characters, invalid UTF-8 and quotes in keys are escaped properly. HTML characters are left as written inside `{!! !!}`, since it's the explicit no-escaping operator.

## v5.1.10 (2026-10-25)

- 🐛 Fixed `@break`, `@continue`, `@breakif` and `@continueif` being silently ignored outside of a loop. They now report `'@break' can only be used inside a loop` (or the matching directive name).
- 🐛 Fixed a NUL byte in a template silently cutting the output. The lexer used `0` to mark the end of input, so everything after the byte disappeared; NUL is now treated as a regular character and the rest of the template renders.
- 🐛 Fixed illegal-token errors showing mojibake for non-ASCII characters. `{{ café }}` now reports `illegal token 'é'` instead of `illegal token 'Ã'`.
- 🐛 Fixed `-9223372036854775808` failing with `cannot parse '9223372036854775808' as 'INT'`. The magnitude doesn't fit into `int64`, but the negative value does, so the literal is now parsed as a whole. Integer arithmetic still wraps around like Go.
- 🐛 Fixed non-finite floats rendering as `+Inf.0`, `-Inf.0` and `NaN.0`. They now print as `+Inf`, `-Inf` and `NaN`, like Go, without a misleading decimal part.

## v5.1.9 (2026-10-25)

- 🐛 Fixed string escapes with backslashes. `{{ "a\\" }}` now gives `a\` instead of swallowing the closing quote, and `{{ "C:\\new" }}` no longer turns the escaped backslash into a newline.
- 🐛 Fixed unclosed comments silently deleting the rest of the template. A `{{-- comment` without the closing `--}}` used to make everything after it disappear from the output, leaving you with no clue why. Now it reports `missing closing '--}}'` pointing at the end of the file.
- 🐛 Fixed `truncate()` cutting bytes instead of characters. `'привіт'.truncate(2)` now gives `пр...` and `'привіт'.truncate(3, '')` gives `при` instead of leaving invalid UTF-8 in the output.
- 🐛 Fixed float division by zero rendering `+Inf.0` or `NaN.0` in the page. `{{ 1.0 / 0.0 }}` now reports the same `division by zero` error as integer division.
- 🐛 Fixed expressions without `;` between them being accepted and glued together. `{{ 1 2 }}` used to render `12` and `{{ name name }}` rendered `AnnaAnna`; now they report `expected ';' but found '2'` (or the token that follows).
- 🐛 Fixed object keys accepting almost any token. `{#: 1}` and `{1.5: 2}` now report `object key must be an identifier or a string, got '...'`, and expression keys like `{'a' + 'b': 1}` are rejected instead of quietly building a different object.

## v5.1.8 (2026-10-24)

- 🐛 Fixed `@for` running the post statement after `@break`, incrementing the loop variable one extra time.
- 🐛 Fixed missing closing `}}` or `!!}` being silently accepted.
- 🐛 Fixed unterminated string literals being accepted instead of raising a syntax error. This `{{ "unterminated }}` would evaluate to `unterminated }}` before.

## v5.1.7 (2026-10-24)

- 🐛 Fixed `str.at()` panicking when a negative index points before the start of the string.
- 🐛 Fixed `slice()` panicking when the start index is after the end index.
- 🐛 Fixed `decimal()` panicking when the number of decimals is negative.
- 🐛 Fixed `repeat()` panicking or exhausting memory when the count is too large.
- 🧑‍💻 Change make commands to work through container instead of relying on local Go installation.

## v5.1.6 (2026-10-24)

- 🐛 Fixed a parser crash when an invalid function call or oversized number was followed by an operator, like `defined() + 1`.
- 🐛 Fixed a crash when `@for` had an empty body.
- 🐛 Fixed `@for` looping forever when its condition evaluated to `nil`.
- 🐛 Fixed `@break` and `@continue` in a nested loop also stopping the outer loop.
- 🐛 Fixed modulo by zero panicking instead of returning a division-by-zero error.

## v5.1.5 (2026-10-22)

- 🐛 Fixed a lexer panic when a template ended with an unclosed `{{`, `{!`, or `{{-`. The same crash could happen in the language server while typing.
- 🐛 Fixed an infinite loop when parsing a malformed object literal, such as a missing `}` or `:`. The language server could freeze on the same input.

## v5.1.4 (2026-10-22)

- 🐛 Remove mutation from `<object>.camel()` function. Now it matches behavior that is written in docs.

## v5.1.3 (2026-09-25)

- 🐛 Fixed wrong URLs to documentation.

## v5.1.2 (2026-09-23)

- 🐛 Fixed wrong scope with `@each` loop and `@if`. It wasn't extending the parent scope. Therefore, this code would result in 0:

```textwire
{{ sum = 0 }}
@each(num in [1, 2, 3])
    {{ sum = sum + num }}
@end
{{ sum }}
```

## v5.1.1 (2026-09-13)

- 🐛 Fixed issue where you couldn't start the server with no Textwire config

## v5.1.0 (2026-08-14)

- 🧑‍💻 Moved from GitHub to Codeberg

## v5.0.1 (2026-05-03)

- 🐛 Fixed bug with `{!!` output

## v5.0.0 (2026-05-03)

- 🐛 Change string escaping. It used to be that strings would be escaped using `raw` function, but now, strings are escaped when they are the result of `{{ <result> }}`. The final result between curly braces is scaped. Function [raw](https://textwire.serhiicho.com/v4/functions/str.html#raw) is removed. Now, you need to use `{!!` and `!!}` to print raw unescaped string. Example: `{!! myData !!}`. Read [upgrade guide](https://textwire.serhiicho.com/v5/upgrade.html).

## Previous Versions

- [Release Notes V1](CHANGELOG-V1.md)
- [Release Notes V2](CHANGELOG-V2.md)
- [Release Notes V3](CHANGELOG-V3.md)
- [Release Notes V4](CHANGELOG-V4.md)
- [Emojis Meaning](EMOJIS.md)
