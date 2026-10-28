<p align="center"><img src="https://codeberg.org/textwire/textwire/raw/branch/master/.forgejo/assets/textwire-logo.png" width="200" height="200" /></p>

# Textwire. A templating language for Go

Textwire is a simple yet powerful templating language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string.


> [!IMPORTANT]
> The project moved from `github.com/textwire/textwire` to `codeberg.org/textwire/textwire`. See [Migrating from GitHub](#migrating-from-github) below.

### [Read Official Documentation](https://textwire.serhiicho.com)

## Installation

```bash
go get codeberg.com/textwire/textwire/v5
```

## Migrating from GitHub

1. Replace the import path in your code:

   ```bash
   grep -rl 'github.com/textwire/textwire' --include='*.go' . | xargs sed -i 's#github.com/textwire/textwire#codeberg.org/textwire/textwire#g'
   ```

2. Fetch the new module and drop the old one:

   ```bash
   go get codeberg.org/textwire/textwire/v5
   go mod tidy
   ```

3. Check that everything still builds and passes:

   ```bash
   go build ./... && go test ./...
   ```

## Neovim and VSCode Support

If you use [Neovim](https://neovim.io/) or [VSCode](https://code.visualstudio.com/) as your primary editor, you can install the [Neovim Plugin](https://codeberg.org/textwire/textwire.nvim) or [VSCode Extension](https://marketplace.visualstudio.com/items?itemName=SerhiiCho.textwire) to get syntax highlighting and other features for Textwire.

## License

The Textwire project is licensed under the [MIT License](https://codeberg.org/textwire/textwire/src/branch/master/LICENSE)

## Thanks To

Thanks to Codeberg and Forgejo for hosting Textwire.

<p align="center">
<img src="https://codeberg.org/textwire/textwire/raw/branch/master/.forgejo/assets/codeberg-logo.png" width="50" height="50" />
<img src="https://codeberg.org/textwire/textwire/raw/branch/master/.forgejo/assets/textwire-logo.png" width="50" height="50" />
<img src="https://codeberg.org/textwire/textwire/raw/branch/master/.forgejo/assets/forgejo-logo.png" width="50" height="50" />
</p>
