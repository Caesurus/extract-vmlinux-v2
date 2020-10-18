package main

import (
	"bytes"
	"compress/gzip"
)

func gUnzipData(data []byte) (resData []byte, err error) {
	b := bytes.NewBuffer(data)

	//var r io.Reader
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
