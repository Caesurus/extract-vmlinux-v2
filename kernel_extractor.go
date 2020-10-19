package main

/*
	Original source: https://github.com/Caesurus/extract-vmlinux-v2
*/

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"sort"
)

type supportedAlgo struct {
	Name        string
	ExtractFunc func(data []byte) (resData []byte, err error)
	Suffix      string
	pattern     []byte
}

// KernelExtractor ...
type KernelExtractor struct {
	data             []byte
	algos            map[string]supportedAlgo
	ignoreValidation bool
}

// NewKernelExtractor ...
func NewKernelExtractor(data *[]byte, ignoreValidation bool) *KernelExtractor {
	k := KernelExtractor{data: *data}
	k.ignoreValidation = ignoreValidation

	k.algos = make(map[string]supportedAlgo)
	k.algos["GZIP"] = supportedAlgo{
		Name:        "GZIP",
		ExtractFunc: extractGzipData,
		Suffix:      "gz",
		pattern:     []byte("\037\213\010"),
	}

	k.algos["BZIP"] = supportedAlgo{
		Name:        "BZIP",
		ExtractFunc: nil,
		Suffix:      "bz",
		pattern:     []byte("BZh"),
	}

	k.algos["LZMA"] = supportedAlgo{
		Name:        "LZMA",
		ExtractFunc: extractLZMAData,
		Suffix:      "lzma",
		pattern:     []byte("\135\000\000"),
	}

	k.algos["LZOP"] = supportedAlgo{
		Name:        "LZOP",
		ExtractFunc: nil,
		Suffix:      "lzop",
		pattern:     []byte("\211\114\132"),
	}

	k.algos["LZ4"] = supportedAlgo{
		Name:        "LZ4",
		ExtractFunc: nil,
		Suffix:      "lz4",
		pattern:     []byte("\002!L\030"),
	}

	k.algos["XZ"] = supportedAlgo{
		Name:        "XZ",
		ExtractFunc: nil,
		Suffix:      "xz",
		pattern:     []byte("\3757zXZ\000"),
	}

	k.algos["ZSTD"] = supportedAlgo{
		Name:        "ZSTD",
		ExtractFunc: nil,
		Suffix:      "zstd",
		pattern:     []byte("(\265/\375"),
	}

	return &k
}
func (k KernelExtractor) isKernelImage(data []byte) bool {
	if k.ignoreValidation {
		return true
	}
	flagLinux := false
	flagParam := false

	flagLinux = bytes.Index(data, []byte("Linux")) > 0
	//flagSyscall := bytes.IndexAny(data, "syscall") > 0
	flagParam = bytes.Index(data, []byte("kernel/params.c")) > 0

	if flagLinux && flagParam {
		return true
	}
	return false
}

func (k KernelExtractor) callExtractor(algo supportedAlgo, offset int) (extractedData []byte, err error) {
	if algo.ExtractFunc != nil {
		extractedData, err := algo.ExtractFunc(k.ReturnBytes(offset))
		if err == nil {
			if k.isKernelImage(extractedData) {
				return extractedData, err
			}
			err = fmt.Errorf("Doesn't look like that was a valid Kernel Image, use -i to dump extracted content anyway")
		}
	} else {
		err = fmt.Errorf("Currently don't support %s extraction", algo.Name)
	}
	return nil, err
}

// ExtractAll will attempt to extract all recognized compressed files
func (k KernelExtractor) ExtractAll() map[string][]byte {
	var files = make(map[string][]byte)

	for desc, algo := range k.algos {
		found, offsets := k.searchPattern(algo.pattern)
		if found {
			fmt.Printf("%s header found at %v\n", desc, offsets)
			for _, offset := range offsets {
				fmt.Printf("Attempting extraction with %s offset:%d \n", desc, offset)
				data, err := k.callExtractor(algo, offset)

				if nil != err {
					fmt.Println(err)
				} else if len(data) > 0 {
					fmt.Printf("%d bytes extracted\n", len(data))
					filename := fmt.Sprintf("vmlinux_%s_%d.bin", algo.Name, offset)
					files[filename] = data
				}
			}
		} else {
			fmt.Printf(" No %s found\n", desc)
		}
	}
	return files
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
func (k KernelExtractor) searchGZIPPattern() (bool, []int) {
	return k.searchPattern(k.algos["GZIP"].pattern)
}
func (k KernelExtractor) searchBZIPPattern() (bool, []int) {
	return k.searchPattern(k.algos["BZIP"].pattern)
}
func (k KernelExtractor) searchLZMAPattern() (bool, []int) {
	return k.searchPattern(k.algos["LZMA"].pattern)
}
func (k KernelExtractor) searchLZOPPattern() (bool, []int) {
	return k.searchPattern(k.algos["LZOP"].pattern)
}
func (k KernelExtractor) searchLZ4Pattern() (bool, []int) {
	return k.searchPattern(k.algos["LZ4"].pattern)
}
func (k KernelExtractor) searchXZPattern() (bool, []int) {
	return k.searchPattern(k.algos["XZ"].pattern)
}
func (k KernelExtractor) searchZSTDPattern() (bool, []int) {
	return k.searchPattern(k.algos["ZSTD"].pattern)
}

// ListAllHeadersFound ...
func (k KernelExtractor) ListAllHeadersFound() error {

	for desc, algo := range k.algos {
		found, offsets := k.searchPattern(algo.pattern)
		if found {
			fmt.Printf("%s header found at %v\n", desc, offsets)
		} else {
			fmt.Printf(" No %s found\n", desc)
		}

	}

	return nil
}
