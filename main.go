package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/mpetavy/common"
	"os"
	"regexp"
	"strings"
)

var (
	file         = flag.String("f", "", "input file")
	search       common.MultiValueFlag
	breaker      common.MultiValueFlag
	breakerCount = flag.Int("bc", 1, "breaker count before")
	lineLength   = flag.Int("ll", 120, "max line length")
)

func init() {
	common.Init("1.0.0", "", "", "2018", "log file extractor", "mpetavy", fmt.Sprintf("https://github.com/mpetavy/%s", common.Title()), common.APACHE, nil, nil, nil, run, 0)

	flag.Var(&search, "s", "text to search for")
	flag.Var(&breaker, "b", "text to bread after")
}

func run() error {
	var breakerLines []string

	if !common.FileExists(*file) {
		return &common.ErrFileNotFound{FileName: *file}
	}

	ba, err := os.ReadFile(*file)
	if common.Error(err) {
		return err
	}

	searchRegexps := make([]*regexp.Regexp, len(search))
	for i, s := range search {
		searchRegexps[i], err = regexp.Compile(s)
		if common.Error(err) {
			return err
		}
	}

	breakerRegexps := make([]*regexp.Regexp, len(breaker))
	for i, s := range breaker {
		breakerRegexps[i], err = regexp.Compile(s)
		if common.Error(err) {
			return err
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(ba)))
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) > *lineLength {
			line = fmt.Sprintf("%s...", line[:*lineLength])
		}

		for _, searchRegexp := range searchRegexps {
			if searchRegexp.MatchString(line) {
				fmt.Printf("%s\n", line)
			}
		}

		for _, breakerRegexp := range breakerRegexps {
			if breakerRegexp.MatchString(line) {
				for _, bl := range breakerLines {
					fmt.Printf("%s\n", bl)
				}
				fmt.Printf("%s\n", line)
				fmt.Println("-----------------------------------------")
			}
		}

		if *breakerCount > 0 {
			if len(breakerLines) == *breakerCount {
				breakerLines = breakerLines[1:]
			}

			breakerLines = append(breakerLines, line)
		}
	}

	return nil
}

func main() {
	defer common.Done()

	common.Run([]string{"f", "s"})
}
