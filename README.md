# Dip

(Markdown) Document instant preview.

Inspired by [grip].

## Features

- Offline-first: Dip renders Markdown like Github does with the help of [goldmark] and [github-markdown-css] - CSS/JS assets are bundled.
- Flexible: Dip can read from standard input, files, directories and even URLs!
- Fast (reasonably): Dip's server is written in Go with HTMX for live reloading.

Disclaimer: Dip is intended primarily for development, and is not meant to be hosted as an online service.
Live reloading is implemented by Server-Sent Events over HTTP/2 Cleartext so connections are _not_ secure.

## Usage

`dip [<path>] [<address>]`, where path is the directory or file to serve as HTML, and address is the TCP address and/or port.

The default is equivalent to `dip . localhost:8080`, which serves the current directory at port 8080 on localhost.

## Installation

```
go install github.com/ongyx/dip/cmd/dip@latest
```

## Development

Dip is composed of a server (`cmd/`, `pkg/`) that serves minified Typescript (`src/`) to clients.

For ease of installation, the bundled assets are checked into this repository under `pkg/asset/dist`.
You may want to update them periodically with:

```
npm update && npm run build
```

After which you can start the server with `npm run serve`.

## License

Dip is licensed under the MIT License.

[grip]: https://github.com/joeyespo/grip
[goldmark]: https://github.com/yuin/goldmark
[fsnotify]: https://github.com/fsnotify/fsnotify
[github-markdown-css]: https://github.com/sindresorhus/github-markdown-css
