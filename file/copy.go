package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

// Copy directory to target
func CopyDirectory(source, target string) error {
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		targetPath := filepath.Join(target, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(targetPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, targetPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, targetPath); err != nil {
				return err
			}
		default:
			if err := CopyFile(sourcePath, targetPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(targetPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		fInfo, err := entry.Info()
		if err != nil {
			return err
		}

		isSymlink := fInfo.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(targetPath, fInfo.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

// Copy source to target
func CopyFile(source, target string) error {
	data, err := os.ReadFile(source)
	if err != nil {
		return err
	}

	err = os.WriteFile(target, data, fs.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Copy a simlink
func CopySymLink(source, target string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, target)
}
