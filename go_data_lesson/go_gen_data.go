package main

import (
	"flag"
	"fmt"
	"os"
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

func main() {
	flag.Usage = Usage
	key := flag.String("k", "", "Key Name")
	number := flag.Int64("n", -1, "Number Data")
	flag.Parse()

	if *key == "" {
		Usage()
	}

	if *number <= 0 {
		Usage()
	}

	for i := int64(1); i <= *number; i += 1 {
		fmt.Printf("%s=%d\n", *key, i)
	}
}
