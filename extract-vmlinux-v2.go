package main

/*
	Original source: https://github.com/Caesurus/extract-vmlinux-v2
*/

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

func main() {
	parser := argparse.NewParser("extract-vmlinux-v2", "A more robust vmlinux extractor")
	opts := argparse.Options{}
	opts.Required = true

	var kernelFile *os.File = parser.File("f", "file", os.O_RDWR, os.FileMode(0600), &opts)
	var ignoreValidation *bool = parser.Flag("i", "ignore", &argparse.Options{})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(kernelFile)
	if err != nil {
		log.Fatal(err)
	}
	ke := NewKernelExtractor(&data, *ignoreValidation)
	files := ke.ExtractAll()

	for filename, extractedData := range files {
		file, err := ioutil.TempFile("", filename)
		if err != nil {
			log.Fatal(err)
		}
		bytesWritten, err := file.Write(extractedData)
		if err == nil {
			fmt.Printf("Wrote %d bytes to file: %s\n", bytesWritten, file.Name())
		} else {
			fmt.Println("Extraction failed", err)
		}
	}

}
