package main

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

func main() {
	var CliCmd = &cobra.Command{
		Use: "LinkSifter",
		// Short: "A program to check URLs against a list of patterns",
		// Long:  `LinkSifter is a command-line tool that allows you to check a list of URLs against patterns to find matches. It provides flexible matching options and customizable behavior to suit your needs.`,
		Run: CliCmdHelp,
	}

	var filein, fileout, filepattern string
	CliCmd.PersistentFlags().StringVarP(&filein, "input", "i", "", "File containing a list of URLs")
	CliCmd.PersistentFlags().StringVarP(&fileout, "output", "o", "", "File to save results to")
	CliCmd.PersistentFlags().StringVarP(&filepattern, "wordlist", "w", "", "File containing a list of patterns")

	var threads int
	var isverbose bool
	CliCmd.PersistentFlags().IntVarP(&threads, "threads", "t", 50, "Number of concurrent threads")
	CliCmd.PersistentFlags().BoolVarP(&isverbose, "verbose", "v", false, "Enable verbose mode")

	var isregex, isequal, islowercase, isallowercase bool
	CliCmd.PersistentFlags().BoolVarP(&isregex, "regex", "r", false, "Check using regex from the pattern file")
	CliCmd.PersistentFlags().BoolVarP(&isequal, "equal", "e", false, "Check if patterns from the file are equal to the compared value (useful when comparing against a part of the URL)")
	CliCmd.PersistentFlags().BoolVarP(&islowercase, "lowercase", "l", false, "Convert all URLs to lowercase")
	CliCmd.PersistentFlags().BoolVarP(&isallowercase, "all2lowercase", "L", false, "Convert URLs and patterns to lowercase")

	var scanpath, scanrawpath, scanfilename, scanrawquery bool
	CliCmd.PersistentFlags().BoolVarP(&scanpath, "path", "", false, "Check URL Path only (includes the filename)")
	CliCmd.PersistentFlags().BoolVarP(&scanrawpath, "rawpath", "", false, "Check URL Path only with out decoding (includes the filename)")
	CliCmd.PersistentFlags().BoolVarP(&scanfilename, "filename", "", false, "Check the Filename from the URL path only")
	CliCmd.PersistentFlags().BoolVarP(&scanrawquery, "rawquery", "", false, "Check the URL Querys only")

	if err := CliCmd.Execute(); err != nil {
		os.Exit(1)
	}
	homeDir, _ := os.UserHomeDir()
	// -----------------------------------------------------------------------------------------------------
	if strings.HasPrefix(fileout, "~/") {
		filepattern = filepath.Join(homeDir, strings.TrimPrefix(filepattern, "~/"))
	}
	patternLines, err := FileLines(filepattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open '%s': %s\n", filepattern, err.Error())
		return
	}
	patternLines = CleanSlice(patternLines)
	if isallowercase && !isregex {
		patternLines = CleanSlice(Slice2Lowercase(patternLines))
	} else if isregex && isverbose {
		newepatternLines := []string{}
		for _, r := range patternLines {
			_, err := regexp.Compile(r)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Malformed regex: '%s'\n", r)
				continue
			}
			newepatternLines = append(newepatternLines, r)
		}
		patternLines = newepatternLines
	}
	if len(patternLines) == 0 {
		fmt.Fprintf(os.Stderr, "Patterns file '%s' is empty\n", filepattern)
		return
	}
	// -----------------------------------------------------------------------------------------------------
	if strings.HasPrefix(filein, "~/") {
		filein = filepath.Join(homeDir, strings.TrimPrefix(filein, "~/"))
	}
	inLines, err := FileLines(filein)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open '%s': %s\n", filein, err.Error())
		return
	}
	if len(inLines) == 0 {
		fmt.Fprintf(os.Stderr, "URLs file '%s' is empty\n", filein)
		return
	}
	inLines = CleanSlice(inLines)
	// -----------------------------------------------------------------------------------------------------
	if strings.HasPrefix(fileout, "~/") {
		fileout = filepath.Join(homeDir, strings.TrimPrefix(fileout, "~/"))
	}
	if _, err := os.Stat(fileout); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(fileout), os.ModePerm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot create '%s': %s\n", filepath.Dir(fileout), err.Error())
			return
		}
	}
	logdatafile, err := os.OpenFile(fileout, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot create '%s': %s\n", fileout, err.Error())
		return
	}
	dataLogger := log.New(logdatafile, "", 0)
	// -----------------------------------------------------------------------------------------------------
	if isverbose {
		fmt.Fprintf(os.Stderr, "Total URLs/Patterns (%d/%d)\n\n", len(inLines), len(patternLines))
	}
	// -----------------------------------------------------------------------------------------------------
	concurrentWork := make(chan struct{}, threads)
	var wg sync.WaitGroup
	for _, line := range inLines {
		wg.Add(1)
		go func(s string) {
			defer func() {
				<-concurrentWork
				wg.Done()
			}()

			toScan := ""
			u, err := url.Parse(s)
			if err != nil {
				if isverbose {
					fmt.Fprintf(os.Stderr, "'%s', is malformed: %s\n", s, err.Error())
				}
				return
			}
			if scanpath {
				toScan = u.Path
			} else if scanrawpath {
				toScan = u.RawPath
			} else if scanfilename {
				toScan = filepath.Base(u.Path)
			} else if scanrawquery {
				toScan = u.RawQuery
			} else {
				toScan = s
			}
			if len(toScan) <= 1 {
				return
			}
			if (isallowercase || islowercase) && !isregex {
				toScan = strings.ToLower(toScan)
			}

			for _, p := range patternLines {
				if isregex {
					if okey, _ := regexp.MatchString(p, toScan); okey {
						if isverbose {
							fmt.Fprintln(os.Stderr, s)
						}
						dataLogger.Output(2, s+"\n")
						return
					}
				} else if isequal {
					if p == toScan {
						if isverbose {
							fmt.Fprintln(os.Stderr, s)
						}
						dataLogger.Output(2, s+"\n")
						return
					}
				} else {
					if strings.Contains(toScan, p) {
						if isverbose {
							fmt.Fprintln(os.Stderr, s)
						}
						dataLogger.Output(2, s+"\n")
						return
					}
				}
			}

		}(line)
	}
	wg.Wait()
}

func CliCmdHelp(cmd *cobra.Command, args []string) {
	threads, _ := cmd.Flags().GetInt("threads")
	filein, _ := cmd.Flags().GetString("input")
	fileout, _ := cmd.Flags().GetString("output")
	filepattern, _ := cmd.Flags().GetString("wordlist")
	if filein == "" || fileout == "" || filepattern == "" || threads < 1 {
		// fmt.Println("input, output and pattern files are required.")
		_ = cmd.Help()
		os.Exit(1)
	}
}
