# soracom-shell

Interactive shell for [`soracom-cli`](https://github.com/soracom/soracom-cli/). Features include:

- commands and flags automatic and interactive completion
- just works by copying the cross-compiled binary file into the target environment. There is no need to build an environment or solve dependencies.

## Install

...

## Usage

1. Install and configure `soracom-cli` according to repo's README
2. Run `soracom-shell`
3. Enjoy
    - the shell automatically provides suggestions as you type
    - hit <kbd>Tab</kdb> to move cursor
    - <kbd>Space</kbd> to select

## Build

```console
$ make
```

## Test

```console
$ make test
```
