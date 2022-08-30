package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/nobishino/gohackc/tokenizer"
)

var out string

func main() {
	os.Exit(exec())
}

func exec() int {
	flag.StringVar(&out, "o", "Tokens.xml", "output file name")
	flag.Parse()
	if len(flag.Args()) < 1 {
		log.Println("should specify source file")
		return 1
	}
	srcPath := flag.Args()[0]
	src, err := os.Open(srcPath)
	if err != nil {
		log.Println(err)
		return 1
	}
	defer src.Close()
	tknzr := tokenizer.NewTokenizer(src)

	out, err := os.Create(dstPath(srcPath))
	if err != nil {
		log.Println(err)
		return 1
	}
	defer out.Close()

	if err := tokenizer.ToXML(out, tknzr); err != nil {
		log.Println(err)
		return 1
	}
	return 0
}

func dstPath(srcPath string) string {
	dir := filepath.Dir(srcPath)
	return filepath.Join(dir, out)
}
