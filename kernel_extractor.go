package main

import (
	"bytes"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
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
	data []byte
	algos              map[string]supportedAlgo
	ignoreValidation bool
}

// NewKernelExtractor ...
func NewKernelExtractor(data *[]byte, ignoreValidation bool) *KernelExtractor {
	k := KernelExtractor{data: *data}
	k.ignoreValidation = ignoreValidation

	k.algos = make(map[string]supportedAlgo)
	k.algos["GZIP"] = supportedAlgo{
		Name:        "GZIP",
		ExtractFunc: gUnzipData,
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
		ExtractFunc: nil,
		Suffix:      "lzma",
		pattern:     []byte("\135\000\000\000"),
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
	if k.ignoreValidation{
		return true
	}

	flagLinux := bytes.IndexAny(data, "Linux") > 0
	flagSyscall := bytes.IndexAny(data, "syscall") > 0
	//flagParam := bytes.IndexAny(data, "kernel/params.c") > 0

	if flagLinux && flagSyscall {
		return true
	}
	return false
}

func (k KernelExtractor) callExtractor(algo supportedAlgo, offset int) (err error) {
	if algo.ExtractFunc != nil {
		file, err := ioutil.TempFile("", algo.Name)
		if err != nil {
			log.Fatal(err)
		}

		extractedData, err := algo.ExtractFunc(k.ReturnBytes(offset))
		if err == nil {
			if k.isKernelImage(extractedData) {
				bytesWritten, err := file.Write(extractedData)
				if err == nil {
					fmt.Printf("Wrote %d to file: %s\n", bytesWritten, file.Name())
				} else {
					fmt.Println("Extraction failed", err)
				}
			} else {
				fmt.Println("Doesn't look like that was a valid Kernel Image, use -i to dump extracted content anyway")
			}
		}
	} else {
		fmt.Printf("Currently don't support %s extraction\n", algo.Name)
	}
	return err
}

// ExtractAll will attempt to extract all recognized compressed files
func (k KernelExtractor) ExtractAll() error {
	for desc, algo := range k.algos {
		found, offsets := k.searchPattern(algo.pattern)
		if found {
			fmt.Printf("%s header found at %v\n", desc, offsets)
			for _, offset := range offsets {
				fmt.Printf("Attempting extraction with %s offset:%d \n", desc, offset)
				_ = k.callExtractor(algo, offset)
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
