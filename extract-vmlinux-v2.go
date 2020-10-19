package main

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
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(kernelFile)
	if err != nil {
		log.Fatal(err)
	}
	ke := NewKernelExtractor(&data, *ignoreValidation)
	ke.ExtractAll()
}
