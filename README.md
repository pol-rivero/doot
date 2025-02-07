# doot

A fast and simple dotfiles manager that just gets the job done.

<!-- | Install from the AUR | Install it with `go get` | Download the binary |
| --- | --- | --- |
| `yay -S doot` | `go get github.com/pol-rivero/doot` | [GitHub Releases](https://github.com/pol-rivero/doot/releases/tag/latest) |
| ![AUR](https://img.shields.io/aur/version/doot) | ![Go](https://img.shields.io/github/go-mod/go-version/pol-rivero/doot) | ![GitHub Releases](https://img.shields.io/github/v/release/pol-rivero/doot) | -->

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

- [`doot crypt`: Manage private (encrypted) files](docs/encryption.md)

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

## Why make another dotfiles manager?

I've tried many dotfiles managers, but none of them really satisfied me. Here are some managers I've tried:

- [GNU Stow](https://www.gnu.org/software/stow/): Works great for what it does, and there are tons of resources online for learning how to use it, but:
  - It lacks support for host-specific files, which makes it unsuitable when you have more than one machine.
  - It symlinks the whole directory instead of individual files. This means having to add a lot of `.gitignore` files, and also losing **all** those ignored files (which can sometimes be very important) when resetting the branch or switching to another dotfiles manager. I've been there and it's not fun.
  - No native support for encrypted files. This is not a deal-breaker, but it's nice to have.

- [RCM](https://thoughtbot.github.io/rcm/): It's almost perfect, except for a [couple of issues](https://github.com/thoughtbot/rcm/issues/306) (see below) that I worked around by adding post-install hooks. However, it's also very slow, and adding the hooks made it even slower. Linking my personal dotfiles repository with RCM takes around 10 seconds, while `doot` takes [TO BE MEASURED]. 10 seconds doesn't seem like much time to wait, but I add and remove files from my dotfiles constantly, and waiting every time I do that gets old very quickly.
  - RCM doesn't allow linking files/directories that start with a dot and are inside nested directories.
  - It doesn't clean up stale/dead symlinks when a filed is moved/renamed/deleted in the dotfiles repository. 
  - The last release was in 2022 and doesn't seem to be actively maintained.
  - No native support for encrypted files. This is not a deal-breaker, but it's nice to have.

- [yadm](https://yadm.io/) and [chezmoi](https://www.chezmoi.io/): Both are very mature and have tons of features. However, having that many features (especially *templates*) forces them to create dotfiles as *files* instead of *symlinks*, meaning that:
  - Automatic changes to the dotfiles (made by other programs) won't be *automagically* reflected in the dotfiles repository. You need to remember to check for changes periodically.
  - It forces you to use their CLI to manage the dotfiles, and I prefer to use a Git GUI for that.

- [dotbot](https://github.com/anishathalye/dotbot): The configuration file is mandatory, and it requires you to manually list each file to be symlinked, instead of deducing it from the directory structure. This is a deal-breaker for me as it requires double the work for no reason.
