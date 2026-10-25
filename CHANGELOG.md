# Release notes V5

## v5.1.9 (2026-10-25)

- 🐛 Fixed string escapes with backslashes. `{{ "a\\" }}` now gives `a\` instead of swallowing the closing quote, and `{{ "C:\\new" }}` no longer turns the escaped backslash into a newline.
- 🐛 Fixed unclosed comments silently deleting the rest of the template. A `{{-- comment` without the closing `--}}` used to make everything after it disappear from the output, leaving you with no clue why. Now it reports `missing closing '--}}'` pointing at the end of the file.
- 🐛 Fixed `truncate()` cutting bytes instead of characters. `'привіт'.truncate(2)` now gives `пр...` and `'привіт'.truncate(3, '')` gives `при` instead of leaving invalid UTF-8 in the output.
- 🐛 Fixed float division by zero rendering `+Inf.0` or `NaN.0` in the page. `{{ 1.0 / 0.0 }}` now reports the same `division by zero` error as integer division.
- 🐛 Fixed expressions without `;` between them being accepted and glued together. `{{ 1 2 }}` used to render `12` and `{{ name name }}` rendered `AnnaAnna`; now they report `expected ';' but found '2'` (or the token that follows).

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

## [Release Notes V1](.forgejo/CHANGELOG-V1.md)
## [Release Notes V2](.forgejo/CHANGELOG-V2.md)
## [Release Notes V3](.forgejo/CHANGELOG-V3.md)
## [Release Notes V4](.forgejo/CHANGELOG-V4.md)
## [Emojis Meaning](.forgejo/EMOJIS.md)
