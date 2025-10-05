// get_license.go
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/architmishra-15/go-colors"
)

type UserInfo struct {
	name string
	year int
}

// Structure of the JSON the application will get after requesting from the repo on GitHub
type JSON struct {
	License string `json:"license"`
	URL     string `json:"url"`
}

func getLicensesJson() map[string]string {
	resp, err := http.Get(Url)
	if err != nil {
		log.Fatalf("%s[ERROR] Could not get to URL: %v%s\n", colors.FgRed, err, colors.Reset)
		return nil
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("%s[ERROR] Could not read the response body: %v%s\n", colors.FgRed, err, colors.Reset)
		return nil
	}

	var licenses map[string]string
	err = json.Unmarshal(body, &licenses)
	if err != nil {
		log.Fatalf("%s[ERROR] Cannot unmarshal JSON: %v%s\n", colors.FgRed, err, colors.Reset)
		return nil
	}

	return licenses
}

// getUsrInfo function returns the final data struct containing the information of the user's name and the current year
func getUsrInfo() *UserInfo {
	y, _, _ := time.Now().Date()

	reply := exec.Command("git", "config", "--global", "user.name")
	output, err := reply.Output()
	
	name := ""
	if err == nil {
		name = strings.TrimSpace(string(output))  // ADDED TrimSpace to remove newlines
	}

	return &UserInfo{
		name: name,
		year: y,
	}
}

func downloadLicenses() {
	licenses := getLicensesJson()

	if len(licenses) == 0 {  // ADDED check for empty map
		fmt.Printf("%sNo licenses found to download%s\n", colors.FgYellow, colors.Reset)
		return
	}

	successCount := 0  // ADDED success counter
	for key, val := range licenses {
		resp, err := http.Get(val)
		if err != nil {
			fmt.Printf("%sCould not get %s license: %v%s\n", colors.FgRed, key, err, colors.Reset)
			continue
		}
		
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()  // MOVED close here to ensure it happens
		
		if err != nil {
			fmt.Printf("%sCould not read content for %s: %v%s\n", colors.FgRed, key, err, colors.Reset)
			continue
		}
		
		createFile(key, body)
		fmt.Printf("%s✓ Downloaded: %s%s\n", colors.FgGreen, key, colors.Reset)
		successCount++
	}
	
	fmt.Printf("\n%s%sSuccessfully downloaded %d/%d licenses%s\n",
		colors.TrueColor(63, 224, 141), colors.Bold, successCount, len(licenses), colors.Reset)
}

