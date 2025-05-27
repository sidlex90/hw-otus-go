package main

import (
	"errors"
	"io"
	"os"

	//nolint:depguard
	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const iterationLimit int64 = 1024

func Copy(fromPath, toPath string, offset, limit int64) error {
	f, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer f.Close()

	offset, limit, err = normalizeOffsetLimit(f, offset, limit)
	if err != nil {
		return err
	}

	if offset > 0 {
		if _, err := f.Seek(offset, io.SeekStart); err != nil {
			return err
		}
	}

	outfile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outfile.Close()

	return copyWithProgress(outfile, f, limit)
}

func normalizeOffsetLimit(f *os.File, offset int64, limit int64) (int64, int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, 0, err
	}

	fSize := fi.Size()

	if fSize == 0 {
		return 0, 0, ErrUnsupportedFile
	}

	if offset > fSize {
		return 0, 0, ErrOffsetExceedsFileSize
	}

	if limit < 1 {
		limit = fSize
	}

	return offset, limit, nil
}

func copyWithProgress(to io.Writer, from io.Reader, limit int64) error {
	bar := pb.StartNew(int(limit))
	defer bar.Finish()

	readLimit := iterationLimit
	for l := limit; l > 0; {
		if l < iterationLimit {
			readLimit = l
		}

		n, err := io.CopyN(to, from, readLimit)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		bar.Add(int(n))

		l -= n
	}

	return nil
}
