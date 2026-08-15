# Textwire. A templating language for Go

<p align="center">
<a href="https://github.com/textwire/textwire/actions/workflows/go.yml"><img src="https://github.com/textwire/textwire/actions/workflows/go.yml/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/textwire/textwire"><img src="https://goreportcard.com/badge/github.com/textwire/textwire"></a>
<a href="https://github.com/textwire/textwire/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
</p>

<p align="center"><img src="https://avatars.githubusercontent.com/u/152325030" width="150" height="150" /></p>

Textwire is a simple yet powerful templating language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string.


Moved to [codeberg.org/textwire/textwire](https://codeberg.org/textwire/textwire)

To upgrade to the new location and receive updates, change all of your references in Go files from GitHub to Codeberg:

```bash
go get codeberg.org/textwire/textwire/v5
```

Change imports in source:

```diff
- github.com/textwire/textwire/v5
+ codeberg.org/textwire/textwire/v5
```

After all references are changed, run tidy command:

```bash
go mod tidy
```
