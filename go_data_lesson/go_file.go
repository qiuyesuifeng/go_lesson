package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
	os.Exit(1)
}

func DoProcess(input []*os.File, output *os.File, arg interface{}) {
	// example
	// 处理一个文件数据
	br := bufio.NewReader(input[0])
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		} else {
			realLine := strings.TrimRight(line, "\n")

			// v=XXXX
			fields := strings.Split(realLine, "=")
			if len(fields) != 2 {
				continue
			}

			value := fields[1]
			data := value + "\n"
			_, err := output.Write([]byte(data))
			if err != nil {
				// do something
				fmt.Printf("[DoProcess Fail][Line]%s\n", realLine)
			}
		}
	}
}

func main() {
	flag.Usage = Usage
	input := flag.String("i", "", "Input File Name")
	output := flag.String("o", "", "Output File Name")
	example := flag.String("e", "", "Example Name")
	flag.Parse()

	if *input == "" {
		Usage()
	}

	if *output == "" {
		Usage()
	}

	if *example != "single" && *example != "multi" {
		Usage()
	}

	var err error
	outputFile, err := os.OpenFile(*output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	CheckError(err)
	defer outputFile.Close()

	var inputFiles []*os.File
	inputs := strings.Split(*input, ",")
	for _, v := range inputs {
		inputFile, err := os.OpenFile(v, os.O_RDONLY, 0)
		CheckError(err)
		defer inputFile.Close()

		inputFiles = append(inputFiles, inputFile)
	}

	DoProcess(inputFiles, outputFile, *example)
}
