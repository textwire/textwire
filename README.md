# Textwire. A templating language for Go

<p align="center">
<a href="https://github.com/textwire/textwire/actions/workflows/go.yml"><img src="https://github.com/textwire/textwire/actions/workflows/go.yml/badge.svg"></a>
<a href="https://goreportcard.com/report/github.com/textwire/textwire"><img src="https://goreportcard.com/badge/github.com/textwire/textwire"></a>
<a href="https://github.com/textwire/textwire/blob/master/LICENSE"><img src="https://img.shields.io/badge/License-MIT-yellow.svg"></a>
</p>

<p align="center"><img src="https://textwire.github.io/img/logo.png" width="150" height="150" /></p>

Textwire is a simple yet powerful templating language for Go. It is designed to easily inject variables from Go code into a template file or just a regular string. It is inspired by Go's syntax and has a similar syntax to make it easier for Go developers to learn and use it.

### [Read Official Documentation](https://textwire.github.io)

## Installation

```bash
go get github.com/textwire/textwire/v2
```

## Neovim and VSCode Support
If you use [Neovim](https://neovim.io/) or [VSCode](https://code.visualstudio.com/) as your primary editor, you can install the [Neovim Plugin](https://github.com/textwire/textwire.nvim) or [VSCode Extension](https://marketplace.visualstudio.com/items?itemName=SerhiiCho.textwire) to get syntax highlighting and other features for Textwire.

## License
The Textwire project is licensed under the [MIT License](https://github.com/textwire/textwire/blob/master/LICENSE)

## Contribute
### Without Container Engine
You'll need to have Go installed on your machine.

### With Container Engine
If you use a container engine like [🦦 Podman](https://podman.io/) or [🐳 Docker](https://app.docker.com/), here are the steps that you can make:

#### Build an image
To build an image, navigate to the root of the project and run this command:
```bash
docker compose build
# for Podman, use:
podman-compose build
```

#### Run the container
To run a container, navigate to the root of the project and run this command:
```bash
docker compose run --rm app
# for Podman, use:
podman-compose run --rm app
```

Optionally, if you want to be able to run `make run` in your container to check Textwire page in your browser (for manual testing purposes), then you need to run this:
```bash
docker compose run --rm -p 8080:8080 app
# for Podman, use:
podman-compose run --rm -p 8080:8080 app
```