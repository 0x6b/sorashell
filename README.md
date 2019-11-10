# soracom-shell

Interactive shell for [`soracom-cli`](https://github.com/soracom/soracom-cli/). Features include:

- commands and flags automatic and interactive completion
- just works by copying the cross-compiled binary file into the target environment. There is no need to build an environment or solve dependencies.

## Install

...

## Usage

1. Install and configure `soracom-cli` according to the repo's README
2. Run `soracom-shell`
3. The shell automatically provides suggestions as you type. Enjoy!
    - <kbd>Tab</kbd> to move cursor
    - <kbd>Space</kbd> to select

## Development

### Setup

1. Install [Go](https://golang.org/) 1.13 or later
2. Install [rakyll/statik](https://github.com/rakyll/statik), asset embedder: `go get github.com/rakyll/statik`

### Build

```console
$ make soracom-shell
```

### Test

```console
$ make test
```

## Acknowledgements

[c-bata/go-prompt](https://github.com/c-bata/go-prompt/) and [soracom/soracom-cli](https://github.com/soracom/soracom-cli/), which does all the heavy lifting.

## License

MIT. See [LICENSE](LICENSE) for details.
