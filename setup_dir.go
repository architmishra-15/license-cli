// setup_dir.go
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/architmishra-15/go-colors"
)

func setupDir(path string) {
	var usrconfirm string
	fmt.Printf("Do you wish to download the licenses from GitHub? [Y/n]: ")
	fmt.Scanf("%5s", &usrconfirm)

	if strings.ToLower(usrconfirm) == "y" || strings.ToLower(usrconfirm) == "yes" || usrconfirm == "" {
		fmt.Printf("%sStarting the download...%s\n", colors.FgGreen, colors.Reset)
		beginDownload(path)
	} else {
		fmt.Printf("%sPlease download the licenses in order to use them.%s\n", colors.FgYellow, colors.Reset)
		fmt.Println("You can download them later by deleting the config directory and running the tool again.")
		os.Exit(0)
	}
}

func beginDownload(path string) {
	// Create the license directory
	if err := os.MkdirAll(licenseLocation, 0755); err != nil {
		log.Fatalf("%sError creating directory: %v%s\n", colors.FgRed, err, colors.Reset)
	}
	fmt.Printf("%sCreated directory: %s%s\n", colors.FgGreen, licenseLocation, colors.Reset)

	// Download all licenses
	downloadLicenses()
}

func createFile(title string, content []byte) {
	path := filepath.Join(licenseLocation, strings.ToLower(title))  // CHANGED to use filepath.Join and make lowercase

	err := os.WriteFile(path, content, 0644)
	if err != nil {
		log.Printf("%sCould not write to file: %s%s\n", colors.FgRed, path, colors.Reset)
		return  // CHANGED from panic to return
	}
}

