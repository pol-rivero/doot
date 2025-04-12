# Dotfiles example

This directory contains a simple example of a dotfiles repository using `doot`.

This file (`README.md`) will **not** be installed to `$HOME` because it is listed in the `exclude_files` list of the configuration file (`doot/config.toml`).

## Important files

- [`doot/config.toml`](doot/config.toml): This is the configuration file for `doot`. It's optional if you don't need to customize the default behavior.
- [`doot/hooks`](doot/hooks): Directory with scripts that will be executed upon certain events.
- [`doot/commands`](doot/commands): Directory with custom commands that can be executed as `doot <command>`.
- [`laptop-dotfiles`](laptop-dotfiles): Host-specific directory with files that will be symlinked only if the machine name is `my-laptop`.
- [`root-dotfiles`](root-dotfiles): A nested doot directory that targets `/` instead of `$HOME`. They can be installed using the `doot root` custom command ([`doot/commands/root`](doot/commands/root)).
