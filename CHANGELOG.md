# Release notes V6

## v6.0.0 (2026-11-06)

[Full migration guide from v5 to v6.](https://textwire.serhiicho.com/v6/upgrade)

- ✨ BREAKING CHANGE! `loop` object scope is now relative to the current loop instead of belonging only to a parent. If you want to get parent loop object use `loop.parent`. Closes [#138](https://codeberg.org/textwire/textwire/issues/138).
- 🧑‍💻 BREAKING CHANGE! `@break`, `@breakif`, `@continue` and `@continueif` now report a parser error when they target a loop outside of `@component` instead of being silently ignored. Closes [#160](https://codeberg.org/textwire/textwire/issues/160).
- ✨ `@for` loops now support the `loop` object with `index`, `iter`, `first` and `parent`. `last` is not available because a `for` loop may not have a known length. Closes [#130](https://codeberg.org/textwire/textwire/issues/130).
- 🐛 Removed single quotes around directives in error messages, e.g. `@elseif cannot come after @else` instead of `'@elseif' cannot come after '@else'`.
- 🧑‍💻 BREAKING CHANGE! Whitespace is handled by lines. A line that holds only a directive, like `@if`, `@each`, `@end`, `@component` or `@break`, or only `{{ }}` or `{!! !!}` embedded code, is removed with its indentation and line ending. Any other line is kept exactly as written. Block bodies are no longer trimmed at their edges, and component bodies are no longer trimmed either. Closes [#166](https://codeberg.org/textwire/textwire/issues/166).

## [Release Notes V1](docs/CHANGELOG-V1.md)
## [Release Notes V2](docs/CHANGELOG-V2.md)
## [Release Notes V3](docs/CHANGELOG-V3.md)
## [Release Notes V4](docs/CHANGELOG-V4.md)
## [Release Notes V5](docs/CHANGELOG-V5.md)
## [Emojis Meaning](docs/EMOJIS.md)
