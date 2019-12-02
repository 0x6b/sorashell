# Contributing

## Setup Development Prerequisites

- Install [Go](https://golang.org/) 1.13 or later
- Install [rakyll/statik](https://github.com/rakyll/statik), asset embedder: `go get github.com/rakyll/statik`

## Fork on GitHub

Before you do anything else, login on [GitHub](https://github.com/) and [fork](https://help.github.com/articles/fork-a-repo/) this repository

## Clone Your Fork Locally

Install [Git](https://git-scm.com/) and clone your forked repository locally.

```sh
$ git clone https://github.com/<your-account>/sorashell.git
```

## Play with Your Fork

The project uses [Semantic Versioning 2.0.0](http://semver.org/) but you don't have to maintain version number.

1. Open your terminal, navigate to local repository directory
2. Build
   ```sh
   $ make sorashell
   ```
3. Create a new topic branch
   ```sh
   $ git checkout -b add-new-feature
   ```
4. Modify source code
5. Test
   ```sh
   $ make test
   ```

## Open a Pull Request

1. Commit your changes locally, [rebase onto upstream/master](https://github.com/blog/2243-rebase-and-merge-pull-requests), then push the changes to GitHub
   ```sh
   $ git push origin add-new-feature
   ```
2. Go to your fork on GitHub, switch to your topic branch, then click "Compare and pull request" button.

## References

- [c-bata/go-prompt](https://github.com/c-bata/go-prompt/)
- [soracom/soracom-cli](https://github.com/soracom/soracom-cli/)
- [SORACOM CLI 利用ガイド | SORACOM Developers](https://dev.soracom.io/jp/docs/cli_guide/)
