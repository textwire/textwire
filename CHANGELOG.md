# Release notes V5

## v5.1.6 (2026-10-24)

- 🐛 Fixed a parser crash when an invalid function call or oversized number was followed by an operator, like `defined() + 1`.
- 🐛 Fixed a crash when `@for` had an empty body.

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
