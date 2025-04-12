# doot

A fast and simple dotfiles manager that just gets the job done.

<table>
  <tr>
    <th>Install from the AUR <br> <a href="https://aur.archlinux.org/packages/doot-bin"><img src="https://img.shields.io/aur/version/doot-bin"/></a></th>
    <td><strong>Pre-compiled binary</strong> <br> <code>yay -S doot-bin</code></td>
    <td><strong>Build from source</strong> <br> <code>yay -S doot</code></td>
    <td><strong>Latest Git commit</strong> <br> <code>yay -S doot-git</code></td>
  </tr>
  <tr>
    <th>Download from GitHub <br> <a href="https://github.com/pol-rivero/doot/releases/latest"><img src="https://img.shields.io/github/v/release/pol-rivero/doot"/></a></th>
    <td><strong><a href="https://github.com/pol-rivero/doot/releases/latest">Download latest</a></strong></td>
    <td><a href="https://github.com/pol-rivero/doot/releases">Older releases</a></td>
    <td></td>
  </tr>
</table>


## Usage

### Create or update symlinks

Simply run `doot` (or `doot install`) from anywhere in your system. It will symlink all files in your dotfiles directory to your home directory, creating directories as needed.  
The subsequent runs will incrementally update the symlinks, adding the new files and directories, and removing references to files that are no longer in the dotfiles directory.

```sh
git clone https://your-dotfiles.git ~/.dotfiles # (or any other directory)

doot  # Installs or updates the symlinks
```

To remove the symlinks, run:

```sh
doot clean
```

Pass `--full-clean` to the `install` or `clean` commands to search for all symlinks that point to the dotfiles directory, even if they were created by another program. This is useful if you created symlinks manually or your dotfiles installation has somehow become corrupted. 


### Add a new file to the dotfiles directory

You could manually move the file to the dotfiles directory and run `doot` to symlink it, but there's a command to do it in one step:

```sh
doot add ./some/file [/other/file ...]
```

Pass `--crypt` to add a file as a private (encrypted) file. See [the documentation](https://github.com/pol-rivero/doot/wiki/Private-(encrypted)-files) for more information.  
If you have more than one machine and this file is only applicable to the current one, pass `--host` to add it as a host-specific file. See `hosts` in the configuration file below.


### Advanced usage

- [`doot crypt`: Manage private (encrypted) files](https://github.com/pol-rivero/doot/wiki/Private-(encrypted)-files)

- [`doot bootstrap`: Automatically download and apply your dotfiles](https://github.com/pol-rivero/doot/wiki/Bootstrap)

- [Hooks: Run custom scripts before and after the installation process](https://github.com/pol-rivero/doot/wiki/Hooks)

- [Need more control? Create your own custom commands](https://github.com/pol-rivero/doot/wiki/Custom-Commands)


## Example

See the [`example`](example) directory for a complete example of a simple dotfiles repository.

## Dotfiles directory location

By default, `doot` searches for your dotfiles in commonly used directories. In order of priority, it looks for the first directory that exists:

1. `$DOOT_DIR`

2. `$XDG_DATA_HOME/dotfiles` (or `$HOME/.local/share/dotfiles` if `XDG_DATA_HOME` is not set)

3. `$HOME/.dotfiles`

Notice how you can set the `DOOT_DIR` environment variable to use any custom directory. The first time you run `doot`, if that variable is not yet defined globally, you can set it inline:

```sh
DOOT_DIR=/path/to/your/dotfiles doot
```

After that, if you have `DOOT_DIR` set in your shell configuration file (`~/.bashrc` or equivalent), you can just run `doot` as usual.

## Configuration file

`doot` reads an optional configuration file: `<dotfiles dir>/doot/config.toml`. This file won't be symlinked when installing. These are the available options and their default values:

```toml
# The target directory for the symlinks. Can contain environment variables.
target_dir = "$HOME"

# Files and directories to ignore. Each entry is a glob pattern relative to the dotfiles directory.
# IMPORTANT: Hidden files/directories are ignored by default. If you set `implicit_dot` to false, you should remove the `**/.*` pattern from this list.
exclude_files = [
  "**/.*",
  "LICENSE",
  "README.md",
]

# Files and directories that are always symlinked, even if they start with a dot or match a pattern in `exclude_files`. Each entry is a glob pattern relative to the dotfiles directory.
include_files = []

# You can get a large performance boost by setting this to `false`, but read this first:
# https://github.com/pol-rivero/doot/wiki/Tip:-set-explore_excluded_dirs-to-false
explore_excluded_dirs = true

# If set to true, files and directories in the root of the dotfiles directory will be prefixed with a dot. For example, `<dotfiles dir>/config/foo` will be symlinked to `~/.config/foo`.
# This is useful if you don't want to have hidden files in the root of the dotfiles directory.
implicit_dot = true

# Top-level files and directories that won't be prefixed with a dot if `implicit_dot` is set to true. Each entry is the name of a file or directory in the root of the dotfiles directory.
implicit_dot_ignore = [
  "bin"
]

# Command and flags to use for displaying diffs. Use any tool and format you like, but it must accept 2 positional arguments for the files to compare.
diff_command = "diff --unified --color=always"

# Key-value pairs of "host name" -> "host-specific directory".
# In the example below, <dotfiles dir>/laptop-dots/.zshrc will be symlinked to ~/.zshrc, taking precedence over <dotfiles dir>/.zshrc, if the hostname is "my-laptop".
# If `implicit_dot` is set to true, the host-specific directories also count as top-level. For example, <dotfiles dir>/laptop-dots/config/foo will be symlinked to ~/.config/foo.
[hosts]
# my-laptop = "laptop-dots"
```
