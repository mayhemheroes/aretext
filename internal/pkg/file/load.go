package file

import (
	"io"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/wedaly/aretext/internal/pkg/text"
)

// Load reads a file from disk and starts a watcher to detect changes.
// This will remove the POSIX end-of-file indicator (line feed at end of file).
func Load(path string, watcherPollInterval time.Duration) (*text.Tree, Watcher, error) {
	f, err := os.Open(path)
	if err != nil {
		// Return the error directly so callers can use os.IsNotExist(err) to check if the file exists.
		return nil, nil, err
	}
	defer f.Close()

	lastModifiedTime, size, err := lastModifiedTimeAndSize(f)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "lastModifiedTime()")
	}

	tree, checksum, err := readContentsAndChecksum(f)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "readContentsAndChecksum()")
	}

	// POSIX files end with a single line feed to indicate the end of the file.
	// We remove it from the tree to simplify editor operations; we'll add it back when saving the file.
	removePosixEof(tree)

	watcher := newFileWatcher(watcherPollInterval, path, lastModifiedTime, size, checksum)

	return tree, watcher, nil
}

func readContentsAndChecksum(f *os.File) (*text.Tree, string, error) {
	checksummer := NewChecksummer()
	r := io.TeeReader(f, checksummer)
	tree, err := text.NewTreeFromReader(r)
	if err != nil {
		return nil, "", errors.Wrapf(err, "text.NewTreeFromReader()")
	}
	return tree, checksummer.Checksum(), nil
}

func lastModifiedTimeAndSize(f *os.File) (time.Time, int64, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return time.Time{}, 0, errors.Wrapf(err, "f.Stat()")
	}

	return fileInfo.ModTime(), fileInfo.Size(), nil
}

func removePosixEof(tree *text.Tree) {
	if endsWithLineFeed(tree) {
		lastPos := tree.NumChars() - 1
		tree.DeleteAtPosition(lastPos)
	}
}

func endsWithLineFeed(tree *text.Tree) bool {
	reader := tree.ReaderAtPosition(tree.NumChars(), text.ReadDirectionBackward)
	var buf [1]byte
	if n, err := reader.Read(buf[:]); err != nil || n == 0 {
		return false
	}
	return buf[0] == '\n'
}