# Dip

(Markdown) Document instant preview.

Inspired by [grip].

## Features

* Offline-first: Dip renders Markdown with the help of [goldmark] - CSS/JS assets are bundled.
* Flexible: Dip can read from standard input, files, directories and even URLs! (WIP)
* Portable: Dip is cross-platform and runs on any OS supported by [fsnotify].

## Usage

`dip [<path>] [<address>]`, where path is the directory or file to serve as HTML, and address is the TCP address and/or port.

The default is equivalent to `dip . 8080`, which serves the current directory at port 8080.

## Installation

```
go install github.com/ongyx/dip/cmd/dip@latest
```

## Development

Dip uses Go for the backend server and Node.js for the frontend CSS/JS served to clients.

For ease of installation, the bundled assets are checked into this repository under `pkg/static/dist`.
You may want to update them periodically with:

```
npm update
npm run build
```

## License

Dip is licensed under the MIT License.

[grip]: https://github.com/joeyespo/grip
[goldmark]: https://github.com/yuin/goldmark
[fsnotify]: https://github.com/fsnotify/fsnotify
[github-markdown-css]: https://github.com/sindresorhus/github-markdown-css
