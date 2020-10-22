package main

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/caesurus/lz4"

	"github.com/itchio/lzma"
	"github.com/xi2/xz"
)

func extractGzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	r, err := gzip.NewReader(b)
	if err != nil {
		return
	}
	r.Multistream(false)

	var resB bytes.Buffer
	_, err = resB.ReadFrom(r)
	if err != nil {
		return
	}

	resData = resB.Bytes()

	return
}

func extractBzipData(data []byte) (resData []byte, err error) {
	ioreader := bytes.NewReader(data)

	r := bzip2.NewReader(ioreader)

	var buf bytes.Buffer
	n, err := io.Copy(&buf, r)
	// we could get an error due to extra data at the end of the file, just ignore and save what we have.
	if n > 0 {
		return buf.Bytes(), nil
	}

	return buf.Bytes(), err
}

// broken xz implemention
func extractXZData(data []byte) (resData []byte, err error) {
	ioreader := bytes.NewReader(data)

	r, err := xz.NewReader(ioreader, 0)
	if err != nil {
		err = fmt.Errorf("XZ NewReader error %s", err)
	}
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		err = fmt.Errorf("io.Copy error %s", err)
	}
	return buf.Bytes(), err
}

func extractLZMAData(data []byte) (resData []byte, err error) {
	ioreader := bytes.NewReader(data)

	r := lzma.NewReader(ioreader)

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		err = fmt.Errorf("io.Copy error %s", err)
	}
	return buf.Bytes(), err
}

func extractLZ4Data(data []byte) (resData []byte, err error) {
	ioreader := bytes.NewReader(data)

	r := lz4.NewReaderLegacy(ioreader)

	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		err = fmt.Errorf("io.Copy error %s", err)
	}
	return buf.Bytes(), err
}
