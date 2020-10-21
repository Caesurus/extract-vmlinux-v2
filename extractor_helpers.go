package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/caesurus/lz4"

	"github.com/itchio/lzma"
	"github.com/ulikunitz/xz"
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

func extractXZData(data []byte) (resData []byte, err error) {
	ioreader := bytes.NewReader(data)

	r, err := xz.NewReader(ioreader)
	if err != nil {
		err = fmt.Errorf("XZ NewReader error %s", err)
	}
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
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
