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

## License

Dip is licensed under the MIT License.

The following resources are vendored in `static/`:

* **[github-markdown-css]** CSS for rendering HTML similarly to Github.

[grip]: https://github.com/joeyespo/grip
[goldmark]: https://github.com/yuin/goldmark
[fsnotify]: https://github.com/fsnotify/fsnotify
[github-markdown-css]: https://github.com/sindresorhus/github-markdown-css
