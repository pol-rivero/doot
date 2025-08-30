package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pol-rivero/doot/lib/common"
	"github.com/pol-rivero/doot/lib/common/log"
)

func HardlinkOrCopyFile(sourcePath, destinationPath string) error {
	err := os.Link(sourcePath, destinationPath)
	if err == nil {
		return nil
	}
	log.Info("Could not hardlink %s to %s: %v. Falling back to copy.", sourcePath, destinationPath, err)

	return CopyFile(sourcePath, destinationPath)
}

func MoveOrCopyFile(sourcePath, destinationPath string) error {
	err := os.Rename(sourcePath, destinationPath)
	if err == nil {
		return nil
	}
	log.Info("Could not move %s to %s: %v. Falling back to copy + delete.", sourcePath, destinationPath, err)

	if err = CopyFile(sourcePath, destinationPath); err != nil {
		return err
	}
	if err = os.Remove(sourcePath); err != nil {
		return fmt.Errorf("failed to remove %q: %w", sourcePath, err)
	}
	return nil
}

func CopyFile(sourcePath, destinationPath string) error {
	info, err := os.Lstat(sourcePath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(destinationPath), 0o755); err != nil {
		return fmt.Errorf("failed to create parent directory for %q: %w", destinationPath, err)
	}

	if common.IsSymlink(info) {
		return copySymlink(sourcePath, destinationPath)
	}
	if info.Mode().IsRegular() {
		return copyRegularFile(sourcePath, destinationPath, info.Mode())
	}
	return fmt.Errorf("unsupported file type for %q", sourcePath)
}

func copySymlink(sourcePath, destinationPath string) error {
	target, err := os.Readlink(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to read symlink %q: %w", sourcePath, err)
	}

	err = os.Remove(destinationPath)
	if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("failed to remove existing symlink %q: %w", destinationPath, err)
	}

	return os.Symlink(target, destinationPath)
}

func copyRegularFile(sourcePath, destinationPath string, fileMode os.FileMode) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source file %q: %w", sourcePath, err)
	}
	defer in.Close()

	if err := removeIfSymlink(destinationPath); err != nil {
		return fmt.Errorf("failed to remove existing symlink %q: %w", destinationPath, err)
	}

	// Create with a safe default, permissions will be fixed after copy
	out, err := os.OpenFile(destinationPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("failed to open or create destination file %q: %w", destinationPath, err)
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("failed to copy file contents from %q to %q: %w", sourcePath, destinationPath, err)
	}
	if err := out.Sync(); err != nil {
		return fmt.Errorf("failed to flush file %q: %w", destinationPath, err)
	}

	if err := os.Chmod(destinationPath, fileMode.Perm()); err != nil {
		return fmt.Errorf("failed to change file mode for %q: %w", destinationPath, err)
	}
	return nil
}

func removeIfSymlink(path string) error {
	// Delete symlink so that it will be recreated on copy.
	// Otherwise, OpenFile will open and overwrite the symlink target instead of the symlink itself.
	info, err := os.Lstat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if common.IsSymlink(info) {
		return os.Remove(path)
	}
	return nil
}
