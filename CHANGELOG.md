# Release notes

## v4.0.0 (2026-03-20)

- ⚠️ BREAKING CHANGE! Changed increment and decrement expressions into statements to match Go's style. `{{ x++ }}` and `{{ x-- }}`  are not returning any value now, they equivalent to `{{ x = x + 1 }}` and `{{ x = x - 1 }}`. Closes [#88](https://github.com/textwire/textwire/issues/88). [Read more about this change](https://textwire.github.io/blog/2026-03-20-textwire-v4).
- ⚠️ BREAKING CHANGE! Remove support for `@breakIf` and `@continueIf`, use lovercased versions `@breakif` and `@continueif`.
- 🧑‍💻 `NewTemplate()` function now returns `*fail.Error` instead of Go's `error` type.
- 🧑‍💻 Improved internal logic for evaluator. Separated literal values from nonliteral. Closes [#90](https://github.com/textwire/textwire/issues/90).

## [Release Notes V1](.github/CHANGELOG-V1.md)

## [Release Notes V2](.github/CHANGELOG-V2.md)

## [Release Notes V3](.github/CHANGELOG-V3.md)

## [Emojis Meaning](.github/EMOJIS.md)
