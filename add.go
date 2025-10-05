// add.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/architmishra-15/go-colors"
)

func addLicense(licenseName, author, year string) {
	// Normalize license name (case-insensitive)
	licenseName = strings.ToLower(licenseName)

	// Check if license exists locally
	licensePath := filepath.Join(licenseLocation, licenseName)
	content, err := os.ReadFile(licensePath)

	if err != nil {
		if os.IsNotExist(err) {
			// Try to find similar license names
			fmt.Printf("%sLicense '%s' not found locally%s\n", colors.FgRed, licenseName, colors.Reset)
			findSimilarLicense(licenseName)
			os.Exit(1)
		}
		fmt.Printf("%sError reading license: %v%s\n", colors.FgRed, err, colors.Reset)
		os.Exit(1)
	}

	// Get user info if not provided
	userInfo := getUsrInfo()
	if author == "" {
		author = strings.TrimSpace(userInfo.name)
	}
	if year == "" {
		year = strconv.Itoa(userInfo.year)
	}

	// Replace placeholders
	licenseText := FindAndReplace(string(content), year, author)

	// Write to current directory
	outputPath := filepath.Join(".", "LICENSE")
	err = os.WriteFile(outputPath, []byte(licenseText), 0644)
	if err != nil {
		fmt.Printf("%sError writing LICENSE file: %v%s\n", colors.FgRed, err, colors.Reset)
		os.Exit(1)
	}

	fmt.Printf("%s%sSuccessfully created LICENSE file with %s license%s\n",
		colors.FgGreen, colors.Bold, licenseName, colors.Reset)
	fmt.Printf("Author: %s%s%s\n", colors.FgCyan, author, colors.Reset)
	fmt.Printf("Year: %s%s%s\n", colors.FgCyan, year, colors.Reset)
}

func findSimilarLicense(targetName string) {
	// Read all files in license directory
	files, err := os.ReadDir(licenseLocation)
	if err != nil {
		return
	}

	var suggestions []string
	minDistance := 3 // Only suggest if distance is <= 3

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		distance := LevenshteinDistance(targetName, strings.ToLower(file.Name()))
		if distance <= minDistance {
			suggestions = append(suggestions, file.Name())
		}
	}

	if len(suggestions) > 0 {
		fmt.Printf("\n%sDid you mean:%s\n", colors.FgYellow, colors.Reset)
		for _, suggestion := range suggestions {
			fmt.Printf("  - %s\n", suggestion)
		}
	}
}

func listLicenses() {
	files, err := os.ReadDir(licenseLocation)
	if err != nil {
		fmt.Printf("%sError reading license directory: %v%s\n", colors.FgRed, err, colors.Reset)
		os.Exit(1)
	}

	fmt.Printf("%s%sAvailable licenses:%s\n", colors.FgGreen, colors.Bold, colors.Reset)
	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("  - %s\n", file.Name())
		}
	}
}

func getCurrentLicense() {
	licensePath := filepath.Join(".", "LICENSE")
	content, err := os.ReadFile(licensePath)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%sNo LICENSE file found in current directory%s\n", colors.FgYellow, colors.Reset)
			os.Exit(1)
		}
		fmt.Printf("%sError reading LICENSE file: %v%s\n", colors.FgRed, err, colors.Reset)
		os.Exit(1)
	}

	fmt.Println(string(content))
}

