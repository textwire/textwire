# Release notes V5

## v5.1.0 (2026-08-14)

- 🧑‍💻 Moved from GitHub to Codeberg

## v5.0.1 (2026-05-03)

- 🐛 Fixed bug with `{!!` output

## v5.0.0 (2026-05-03)

- 🐛 Change string escaping. It used to be that strings would be escaped using `raw` function, but now, strings are escaped when they are the result of `{{ <result> }}`. The final result between curly braces is scaped. Function [raw](https://textwire.codeberg.page/v4/functions/str#raw) is removed. Now, you need to use `{!!` and `!!}` to print raw unescaped string. Example: `{!! myData !!}`. Read [upgrade guide](https://textwire.codeberg.page/v5/upgrade).

## [Release Notes V1](.forgejo/CHANGELOG-V1.md)

## [Release Notes V2](.forgejo/CHANGELOG-V2.md)

## [Release Notes V3](.forgejo/CHANGELOG-V3.md)

## [Release Notes V4](.forgejo/CHANGELOG-V4.md)

## [Emojis Meaning](.forgejo/EMOJIS.md)
