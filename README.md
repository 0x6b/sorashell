# sorashell

Interactive shell for [`soracom-cli`](https://github.com/soracom/soracom-cli/). Features include:

- commands and flags automatic, interactive, fuzzy-match completion
- just works by copying the cross-compiled binary file into the target environment. There is no need to build an environment or solve dependencies.

## Install

...

## Usage

1. Install and configure `soracom-cli` according to the repo's README or [SORACOM CLI 利用ガイド | SORACOM Developers](https://dev.soracom.io/jp/docs/cli_guide/) (Japanese)
2. Run `sorashell`
3. The shell automatically provides suggestions as you type. Enjoy!
    - <kbd>Tab</kbd> to move cursor or display suggestions
    - <kbd>Space</kbd> to select

### Tips and Tricks

You can:

- fuzzy match commands, options, IMSIs, device IDs, etc. Try with `subscribers get --imsi
` and type SIM status (online or offline), or a part of it's name. 
- use pipe (`|`) or redirect (`>`) with your risk e.g.
    ```console
    SORASHELL> subscribers list --speed-class-filter s1.slow | jq '.[] | .imsi' > slow-subscribers.txt 
    ```
- run local shell command by preceding the command with `!` e.g.
    ```console
    SORASHELL> !cd /Users/soracom/Desktop
    SORASHELL> !ls
    ```

## Limitation

- `sorashell` has rough edges, and is not yet suitable for non-technical users.
- macOS only at this moment. PRs welcome!

## Acknowledgements

[c-bata/go-prompt](https://github.com/c-bata/go-prompt/) and [soracom/soracom-cli](https://github.com/soracom/soracom-cli/), which does all the heavy lifting.

## Contributing

Please read [CONTRIBUTING](CONTRIBUTING.md) for details.

## License

MIT. See [LICENSE](LICENSE) for details.
