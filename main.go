// main.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/architmishra-15/go-colors"
)

var licenseLocation string

const Url = "https://raw.githubusercontent.com/architmishra-15/license-cli/main/links.json"

var year_regex = regexp.MustCompile(`(?i)\[yyyy\]`)
var name_regex = regexp.MustCompile(`(?i)\[name of the author\]`)

func FindAndReplace(src string, year string, author string) string {
	result := year_regex.ReplaceAllString(src, year)
	result = name_regex.ReplaceAllString(result, author)
	return result
}

func LevenshteinDistance(s1, s2 string) int {
	str1 := []rune(s1)
	str2 := []rune(s2)

	m := len(str1)
	n := len(str2)

	prevRows := make([]int, n+1)
	currRows := make([]int, n+1)

	for j := 0; j <= n; j++ {
		prevRows[j] = j
	}

	for i := 1; i <= m; i++ {
		currRows[0] = i
		for j := 1; j <= n; j++ {
			if str1[i-1] == str2[j-1] {
				currRows[j] = prevRows[j-1]
			} else {
				currRows[j] = 1 + min(currRows[j-1], prevRows[j], prevRows[j-1])
			}
		}
		prevRows, currRows = currRows, prevRows
	}

	return prevRows[n]
}

func checkLicenceDir(path string) bool {
	file_in, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false
	}

	return file_in.IsDir()
}

func main() {
	// Setup license location based on OS
	home, _ := os.UserHomeDir()
	if runtime.GOOS == "windows" {
		licenseLocation = filepath.Join(home, "AppData", "Local", "license-cli", "licenses")
	} else {
		licenseLocation = filepath.Join(home, ".config", "license-cli", "licenses")
	}

	// Check if license directory exists
	if !checkLicenceDir(licenseLocation) {
		fmt.Printf("%sNo license directory found in%s %s%s%s%s\n",
			colors.FgRed, colors.Reset,
			colors.FgMagenta, colors.Underline, licenseLocation, colors.Reset)

		parent := filepath.Dir(licenseLocation)
		setupDir(parent)
	}

	// Parse command line arguments
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Printf("%sPlease enter a command%s\n", colors.Bold, colors.Reset)
		help()
		os.Exit(1)
	}

	command := strings.ToLower(args[0])

	switch command {
	case "help", "--help", "-h": help()
	case "version", "--version", "-v": version()
	
	case "add":
		if len(args) < 2 {
			fmt.Printf("%sError: 'add' requires a license name%s\n", colors.FgRed, colors.Reset)
			os.Exit(1)
		}

		licenseName := args[1]
		author := ""
		year := ""

		// Parse flags
		for i := 2; i < len(args); i++ {
			switch args[i] {
			case "--author", "-a":
				if i+1 < len(args) {
					author = args[i+1]
					i++
				}
			case "--year", "-y":
				if i+1 < len(args) {
					year = args[i+1]
					i++
				}
			}
		}

		addLicense(licenseName, author, year)

	case "list":
		listLicenses()

	case "get":
		getCurrentLicense()

	default:
		fmt.Printf("%sUnknown command: %s%s\n", colors.FgRed, command, colors.Reset)
		help()
		os.Exit(1)
	}
}
