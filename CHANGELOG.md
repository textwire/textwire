# Release notes V6

## v6.0.1 (2026-11-19)

- 🐛 The file watcher no longer logs the same template parse error on every tick while the file stays unchanged. It logs the error once and logs it again only after the file is edited. Closes [#126](https://github.com/textwire/textwire/issues/126).
- 🐛 Directive argument lists now accept trailing commas, so multiline directives like `@insert("title", "Home page",)` work. Passing extra arguments now reports how many a directive accepts, like `@insert() accepts at most 2 arguments`, instead of the misleading `expected ')' but found ','` error. Closes [#129](https://github.com/textwire/textwire/issues/129).

## v6.0.0 (2026-11-19)

[Full migration guide from v5 to v6.](https://textwire.serhiicho.com/v6/upgrade)

- 🧑‍💻 BREAKING CHANGE! `defined()` now returns `false` for missing object properties and array indices instead of `true`. A property or index that exists with a `nil` value still counts as defined. Closes [#167](https://codeberg.org/textwire/textwire/issues/167).
- ✨ BREAKING CHANGE! `loop` object scope is now relative to the current loop instead of belonging only to a parent. If you want to get parent loop object use `loop.parent`. Closes [#138](https://codeberg.org/textwire/textwire/issues/138).
- 🧑‍💻 BREAKING CHANGE! `@break`, `@breakif`, `@continue` and `@continueif` now report a parser error when they target a loop outside of `@component` instead of being silently ignored. Closes [#160](https://codeberg.org/textwire/textwire/issues/160).
- ✨ `@for` loops now support the `loop` object with `index`, `iter`, `first` and `parent`. `last` is not available because a `for` loop may not have a known length. Closes [#130](https://codeberg.org/textwire/textwire/issues/130).
- 🐛 Removed single quotes around directives in error messages, e.g. `@elseif cannot come after @else` instead of `'@elseif' cannot come after '@else'`.
- 🧑‍💻 BREAKING CHANGE! Whitespace is handled by lines. A line is removed entirely, including its spaces and line ending, only when exactly one directive or comment is its sole non-whitespace content. That includes lines inside a multi-line directive, like the arguments of `@if(\n  true\n)`. A line with `{{ }}`,`{!! !!}` or `{{-- --}}` embedded code, several directives, or any other content keeps its whitespace exactly as written, and blank lines are unchanged. Block bodies are no longer trimmed at their edges, and component bodies are no longer trimmed either. Closes [#166](https://codeberg.org/textwire/textwire/issues/166).
- 🧑‍💻 BREAKING CHANGE! A component that includes itself, directly or through other components, now stops with an error instead of crashing the program with a stack overflow. The limit is how many components can render inside each other at the same time during one render, 500 by default. Components rendered one after another, like inside of a loop, don't count towards it. Change the limit with the new `MaxComponentDepth` config option.
- 🧑‍💻 BREAKING CHANGE! The short form `@insert('title', title)` and the `@reserve('name', fallback)` fallback are now escaped the same way as `{{ }}`. Previously they printed the value raw, which allowed XSS when user data was passed this way. Use the block form with `{!! !!}` when raw output is needed. Closes [#156](https://codeberg.org/textwire/textwire/issues/156).
- ✨ The file watcher now logs `refreshed <path>` when a template's content changes. Unchanged templates stay quiet.
- 🧑‍💻 BREAKING CHANGE! Writing to `global`, its properties, or a value aliased from it now returns an evaluation error instead of being silently ignored. Globals are frozen and read-only.
- 🧑‍💻 BREAKING CHANGE! `GlobalData` is now converted once when `NewTemplate` is called, so later changes to the original map are not picked up. Unsupported values now fail in `NewTemplate` instead of at render time.
- 🧑‍💻 BREAKING CHANGE! Public `pkg/` API changes: `evaluator.New` takes a third `global *value.Obj` argument, `value.NewScopeFromMap` takes a second `global *value.Obj` argument, and `Scope.AddGlobal` is removed. Pass `nil` to keep the previous behavior.
- 🐛 Errors for missing components, unused inserts, and unknown or duplicate `@pass`/`@passif` slots now point at the name string instead of the whole directive.

## [Release Notes V1](docs/CHANGELOG-V1.md)
## [Release Notes V2](docs/CHANGELOG-V2.md)
## [Release Notes V3](docs/CHANGELOG-V3.md)
## [Release Notes V4](docs/CHANGELOG-V4.md)
## [Release Notes V5](docs/CHANGELOG-V5.md)
## [Emojis Meaning](docs/EMOJIS.md)
