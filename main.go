package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"io/ioutil"
	"strings"
)

var (
	inputFile    = flag.String("i", "", "input file")
	search       common.MultiValueFlag
	breaker      common.MultiValueFlag
	breakerCount = flag.Int("bc", 1, "breaker count before")
	lineLength   = flag.Int("ll", 120, "max line length")
)

func init() {
	common.Init(false, "1.0.0", "", "", "2018", "log file extractor", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)

	flag.Var(&search, "s", "text to search for")
	flag.Var(&breaker, "b", "text to bread after")
}

func run() error {
	var breakerLines []string

	if !common.FileExists(*inputFile) {
		return &common.ErrFileNotFound{FileName: *inputFile}
	}

	ba, err := ioutil.ReadFile(*inputFile)
	if common.Error(err) {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(ba)))
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > *lineLength {
			line = fmt.Sprintf("%s...", line[:*lineLength])
		}

		for _, s := range search {
			if strings.Contains(line, s) {
				fmt.Printf("%s\n", line)
			}
		}

		for _, s := range breaker {
			if strings.Contains(line, s) {
				for _, bl := range breakerLines {
					fmt.Printf("%s\n", bl)
				}
				fmt.Printf("%s\n", line)
				fmt.Println("-----------------------------------------")
			}
		}

		if len(breakerLines) == *breakerCount {
			breakerLines = breakerLines[1:]
		}

		breakerLines = append(breakerLines, line)
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run(nil)
}
