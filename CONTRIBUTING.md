# Contributing to pvu

We want to make contributing to this project as easy and transparent as possible.

## Project structure

- `cmd/pvu` - CLI entrypoint.
- `internal/app` - Main application flow and command handling.
- `internal/k8s` - Kubernetes client setup and cluster interactions.
- `internal/models` - Shared data models used across the project.
- `internal/tui` - Interactive terminal UI.
- `internal/ui` - Rendering helpers and layout code used by the TUI.
- `internal/version` - Version, commit, and build date information.

## Development

Common tasks are available through the Makefile.

Tidy modules:

```shell
make tidy
```

Build the binary locally:

```shell
make build
```

Run the project:

```shell
make run
```

Run the TUI against the default namespace:

```shell
make tui
```

Build local release artifacts with GoReleaser snapshot mode:

```shell
make release
```

Remove local build artifacts:

```shell
make clean
```

## Commits

Commit messages should be well formatted. This project uses Conventional Commits.

```shell
<type>[<scope>]: <short summary>
   │      │             │
   │      │             └─> Summary in present tense. Not capitalized. No period at the end.
   │      │
   │      └─> Scope (optional): for example app, k8s, tui, ui, docs, build
   │
   └─> Type: chore, docs, feat, fix, refactor, style, or test
```

You can follow the specification on their website:
https://www.conventionalcommits.org/en/v1.0.0/

## Pull requests

Pull requests are welcome.

Before opening a pull request:

1. Fork the repository and create your branch from `main`.
2. Keep changes focused and avoid mixing unrelated work in one pull request.
3. Add or update tests when behavior changes.
4. Update documentation when flags, behavior, or install steps change.
5. Make sure the project builds successfully before submitting.

Pull requests are the standard way to propose and review changes before merging them into the main codebase.

## Issues

GitHub issues are used to track bugs, regressions, and feature requests.

When opening an issue, please include:

- A clear summary of the problem.
- Steps to reproduce it.
- Expected behavior.
- Actual behavior.
- Environment details when relevant, such as OS, Kubernetes version, and how the plugin was installed.

If you are reporting a TUI issue, include terminal size or a screenshot when possible.

## Plugin notes

This project is distributed as a kubectl plugin. Kubectl discovers plugins from executables whose names start with `kubectl-` and are available on `PATH`.

If you are making changes related to packaging or install behavior, please verify:

- Manual installation still works with the `kubectl-pvu` binary name.
- Krew manifest changes match the generated archive names.
- Plugin naming stays consistent with Krew naming guidance.

## License

By contributing to pvu, you agree that your contributions will be licensed under the LICENSE file in the root directory of this source tree.