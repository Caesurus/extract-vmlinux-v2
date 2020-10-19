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
	parser := argparse.NewParser(appName, "A more robust vmlinux extractor")
	var inputFile *string = parser.String("f", "file", &argparse.Options{Help: "Input file to be processed"})
	var ignoreValidation *bool = parser.Flag("i", "ignore", &argparse.Options{Help: "Ignore kernel verification, extract whatever is found"})
	var listOnly *bool = parser.Flag("l", "list", &argparse.Options{Help: "List headers and offset found"})
	var printVersion *bool = parser.Flag("V", "version", &argparse.Options{Help: "Print version info"})

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(1)
	}

	if *printVersion {
		fmt.Printf("%s: %d.%d.%d\n", appName, versionMajor, versionMinor, versionPatch)
		os.Exit(0)
	}


	if len(*inputFile) > 0 {
		kernelFile, err := os.Open(*inputFile)
		if err != nil {
			fmt.Println("couldn't open file")
			os.Exit(2)
		}

		data, err := ioutil.ReadAll(kernelFile)
		if err != nil {
			log.Fatal(err)
		}
		ke := NewKernelExtractor(&data, *ignoreValidation)

		if *listOnly{
			ke.ListAllHeadersFound()
			os.Exit(0)
		}

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

}
