package main

import (
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"sort"
)

type supportedAlgo struct {
	Name        string
	SearchFunc  func() (bool, []int)
	ExtractFunc func(data []byte) (resData []byte, err error)
	Suffix      string
}

// KernelExtractor ...
type KernelExtractor struct {
	data []byte
	//algos map[string]func() (bool, []int)
	algos map[string]supportedAlgo
}

// NewKernelExtractor ...
func NewKernelExtractor(data *[]byte) *KernelExtractor {
	k := KernelExtractor{data: *data}
	k.algos = make(map[string]supportedAlgo)
	k.algos["GZIP"] = supportedAlgo{
		Name:        "GZIP",
		SearchFunc:  k.searchGZIPPattern,
		ExtractFunc: gUnzipData,
		Suffix:      "gz",
	}
	/*
		k.algos = make(map[string]func() (bool, []int))
		k.algos["BZIP"] = k.searchBZIPPattern
		k.algos["GZIP"] = k.searchGZIPPattern
		k.algos["LZ4"] = k.searchLZ4Pattern
		k.algos["LZMA"] = k.searchLZMAPattern
		k.algos["LZO"] = k.searchLZOPPattern
		k.algos["XZ"] = k.searchXZPattern
		k.algos["ZSTD"] = k.searchZSTDPattern
	*/

	return &k
}

// ExtractAll will attempt to extract all recognized compressed files
func (k KernelExtractor) ExtractAll() error {
	for desc, algo := range k.algos {
		found, offsets := algo.SearchFunc()
		if found {
			fmt.Printf("%s header found at %v\n", desc, offsets)
			for _, offset := range offsets {
				fmt.Printf("Attempting extraction with %s offset:%d \n", desc, offset)

				file, err := ioutil.TempFile("", desc)
				if err != nil {
					log.Fatal(err)
				}

				extractedData, err := algo.ExtractFunc(k.ReturnBytes(offset))
				if err == nil{
					bytesWritten, err := file.Write(extractedData)
					if err == nil {
						fmt.Printf("Wrote %d to file: %s\n", bytesWritten, file.Name())
					} else {
						fmt.Println("Extraction failed", err)
					}
				}
			}
		} else {
			fmt.Printf(" No %s found\n", desc)
		}

	}
	return nil
}

// ReturnBytes Return []bytes from a given offset into the buffer
func (k KernelExtractor) ReturnBytes(offset int) []byte {
	return k.data[offset:]
}

func (k KernelExtractor) searchPattern(pattern []byte) (bool, []int) {
	index := suffixarray.New(k.data)
	offsets := index.Lookup(pattern, -1)
	if offsets == nil {
		return false, nil
	}
	sort.Ints(offsets)
	return true, offsets
}

func (k KernelExtractor) searchXZPattern() (bool, []int) {
	pattern := []byte("\3757zXZ\000")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchGZIPPattern() (bool, []int) {
	pattern := []byte("\037\213\010")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchBZIPPattern() (bool, []int) {
	pattern := []byte("BZh")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchLZMAPattern() (bool, []int) {
	pattern := []byte("\135\000\000\000")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchLZOPPattern() (bool, []int) {
	pattern := []byte("\211\114\132")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchLZ4Pattern() (bool, []int) {
	pattern := []byte("\002!L\030")
	return k.searchPattern(pattern)
}
func (k KernelExtractor) searchZSTDPattern() (bool, []int) {
	pattern := []byte("(\265/\375")
	return k.searchPattern(pattern)
}

// ListAllHeadersFound ...
func (k KernelExtractor) ListAllHeadersFound() error {

	for desc, algo := range k.algos {

		found, offsets := algo.SearchFunc()
		if found {
			fmt.Printf("%s header found at %v\n", desc, offsets)
		} else {
			fmt.Printf(" No %s found\n", desc)
		}

	}

	return nil
}
