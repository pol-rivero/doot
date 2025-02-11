# doot

A fast and simple dotfiles manager that just gets the job done.

<table>
  <tr>
    <th>Install from the AUR <br> <img src="https://img.shields.io/aur/version/doot-bin"/></th>
    <td><strong>Pre-compiled binary</strong> <br> <code>yay -S doot-bin</code></td>
    <td><strong>Build from source</strong> <br> <code>yay -S doot</code></td>
    <td><strong>Latest Git commit</strong> <br> <code>yay -S doot-git</code></td>
  </tr>
  <tr>
    <th>Download from GitHub <br> <img src="https://img.shields.io/github/v/release/pol-rivero/doot"/></td></th>
    <td><strong><a href="https://github.com/pol-rivero/doot/releases/latest">Compiled binaries</a></strong></td>
    <td></td>
    <td></td>
  </tr>
</table>


## Usage

### Create or update symlinks

Simply run `doot` (or `doot install`) from anywhere in your system. It will symlink all files in your dotfiles directory to your home directory, creating directories as needed.  
The subsequent runs will incrementally update the symlinks, adding the new files and directories, and removing references to files that are no longer in the dotfiles directory.

```sh
git clone https://your-dotfiles.git ~/.dotfiles # or any other directory

doot
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

Pass `--crypt` to add a file as a private (encrypted) file. See [the documentation](docs/encryption.md) for more information.  
If you have more than one machine and this file is only applicable to the current one, pass `--host` to add it as a host-specific file. See `hosts` in the configuration file below.


### Advanced usage

- [`doot crypt`: Manage private (encrypted) files](https://github.com/pol-rivero/doot/wiki/Private-(encrypted)-files)

- [`doot bootstrap`: Automatically download and apply your dotfiles](docs/bootstrap.md)

- [Hooks: Run custom scripts before and after the installation process](docs/hooks.md)


## Dotfiles directory

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

`doot` reads an optional configuration file `<dotfiles dir>/doot/config.toml` in the root of the dotfiles directory. This file won't be symlinked by default. These are the available options and their default values:

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
# IMPORTANT: See warning below!
include_files = []

# If set to true, files and directories in the root of the dotfiles directory will be prefixed with a dot. For example, `<dotfiles dir>/config/foo` will be symlinked to `~/.config/foo`.
# This is useful if you don't want to have hidden files in the root of the dotfiles directory.
implicit_dot = true

# Top-level files and directories that won't be prefixed with a dot if `implicit_dot` is set to true. Each entry is the name of a file or directory in the root of the dotfiles directory.
implicit_dot_ignore = [
  "bin"
]

# Key-value pairs of "host name" -> "host-specific directory".
# In the example below, <dotfiles dir>/laptop-dots/.zshrc will be symlinked to ~/.zshrc, taking precedence over <dotfiles dir>/.zshrc, if the hostname is "my-laptop".
# If `implicit_dot` is set to true, the host-specific directories also count as top-level. For example, <dotfiles dir>/laptop-dots/config/foo will be symlinked to ~/.config/foo.
[hosts]
# my-laptop = "laptop-dots"
```

> [!WARNING]
> The glob matching has some quirks that you should be aware of:
> 1. Unlike standard glob patterns, `**/file.txt` will **NOT** match `file.txt`. [Issue link](https://github.com/gobwas/glob/issues/58).
> 2. When a directory matches a glob in `exclude_files`, it will **NOT** be explored recursively (so its contents will *never* be symlinked, even if they would have matched a glob in `include_files`).  
>   This is done to improve performace and is usually the desired behavior. If you want to exclude all the files in `some-dir` except for `some-dir/images/important.png`, do the following:
>   ```toml
>   # Exclude all children of some-dir, but not some-dir itself, so that it can be explored
>   exclude_files = [ "some-dir/**" ]
>   include_files = [
>       # Include images so that it can be explored. Children are NOT included (no trailing `/**`)
>       "some-dir/images",
>       "some-dir/images/important.png"  # Include the file you want
>   ]
