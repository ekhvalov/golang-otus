package main

import (
	"bufio"
	"errors"
	"io"
	"os"
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
	err = seek(src, offset)
	if err != nil {
		return err
	}
	var bufSrc *bufio.Reader
	if limit > 0 {
		bufSrc = bufio.NewReader(io.LimitReader(src, limit))
	} else {
		bufSrc = bufio.NewReader(src)
	}

	bufDst := bufio.NewWriter(dst)
	_, err = io.Copy(bufDst, bufSrc)
	return err
}

func seek(f *os.File, offset int64) error {
	fileInfo, err := f.Stat()
	if err != nil {
		return err
	}
	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	s := fileInfo.Size()
	if offset > s {
		return ErrOffsetExceedsFileSize
	}
	_, err = f.Seek(offset, 0)
	if err != nil {
		return err
	}
	return nil
}
