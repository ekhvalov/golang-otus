package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	dst, err := os.OpenFile(toPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
		_ = dst.Close()
	}()
	size, err := seek(src, offset)
	if err != nil {
		return err
	}
	var bufSrc *bufio.Reader
	if limit > 0 {
		bufSrc = bufio.NewReader(io.LimitReader(src, limit))
		size = limit
	} else {
		bufSrc = bufio.NewReader(src)
	}
	bufDst := bufio.NewWriter(dst)

	bar := pb.Full.Start64(size)
	barReader := bar.NewProxyReader(bufSrc)

	_, err = io.Copy(bufDst, barReader)
	return err
}

func seek(f *os.File, offset int64) (int64, error) {
	fileInfo, err := f.Stat()
	if err != nil {
		return -1, err
	}
	if !fileInfo.Mode().IsRegular() {
		return -1, ErrUnsupportedFile
	}
	size := fileInfo.Size()
	if offset > size {
		return -1, ErrOffsetExceedsFileSize
	}
	_, err = f.Seek(offset, 0)
	if err != nil {
		return -1, err
	}
	return size - offset, nil
}
