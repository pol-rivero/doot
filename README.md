# doot

A fast, simple and intuitive dotfiles manager that just gets the job done. **[Should you try it?](https://github.com/pol-rivero/doot/wiki/Should-I-use-doot%3F)**

## Install

<details>
<summary>From the <b>AUR</b> (Arch Linux, Manjaro, and other Arch-based distributions)</summary>

&nbsp;  
I recommend installing the `doot-bin` package, which is a pre-compiled binary.

|Pre-compiled binary|Build from source|Latest Git commit|
|---|---|---|
|`yay -S doot-bin`|`yay -S doot`|`yay -S doot-git`|

&nbsp;  
</details>



<details>
<summary>Using <b>Homebrew</b> (Linux and macOS)</summary>

&nbsp;  
Install `doot` from Homebrew:

```sh
brew install pol-rivero/tap/doot
```

Make sure to run `brew update && brew upgrade` periodically to keep `doot` up to date.

&nbsp;  
</details>



<details>
<summary>Using the <b>installer script</b> (Linux and macOS)</summary>

&nbsp;  
Run the following command:

```sh
curl -sSL get-doot.polrivero.com | sh
```

- You can inspect the script before running it: `curl -sSL get-doot.polrivero.com | cat`

- Make sure to run this command periodically or set up a cron job in order to keep `doot` up to date.

- To uninstall, run the following command: `sudo rm $(which doot)`

&nbsp;  
</details>



<details>
<summary>Linux manual installation</summary>

&nbsp;  
Go to the [latest GitHub release](https://github.com/pol-rivero/doot/releases/latest) and download either `doot-linux-x86_64` or `doot-linux-arm64` depending on your architecture, rename it to `doot`.  
Make it executable and move it to any directory in your `PATH`:

```sh
chmod +x doot
sudo mv doot /usr/local/bin
```

**Want to contribute?**  
If your distribution doesn't have a package for `doot`, consider helping out by creating and submitting it to your distribution's package manager. Please [open an issue](https://github.com/pol-rivero/doot/issues) in order to discuss it and coordinate the effort.

&nbsp;  
</details>



<details>
<summary>macOS manual installation</summary>

&nbsp;  
Go to the [latest GitHub release](https://github.com/pol-rivero/doot/releases/latest) and download either `doot-darwin-x86_64` or `doot-darwin-arm64` depending on your architecture, rename it to `doot`.  
Make it executable and move it to any directory in your `PATH`:

```sh
chmod +x doot
sudo mv doot /usr/local/bin
```

&nbsp;  
</details>



<details>
<summary>Windows</summary>

&nbsp;  
**Windows is not officially supported.** I'm not sure how Windows handles symlinks, so I can't guarantee that `doot` will work as expected.  
If you want to give it a try, you can download the latest release from the [GitHub releases page](https://github.com/pol-rivero/doot/releases/latest).

&nbsp;  
</details>

## Usage

If this is your first time setting up a dotfiles repository, read the [Getting Started](https://github.com/pol-rivero/doot/wiki/Getting-Started) guide.

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

After that, if you have set `DOOT_DIR` in your shell configuration file (`~/.bashrc` or equivalent), you can just run `doot` as usual.

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

# Files and directories that are always symlinked, overriding `exclude_files`. Each entry is a glob pattern relative to the dotfiles directory.
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
