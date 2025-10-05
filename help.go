// help.go
package main

import (
	"fmt"

	"github.com/architmishra-15/go-colors"
)

func help() {
	fmt.Printf("%s%sLicense CLI - Manage project licenses easily%s\n\n", colors.FgCyan, colors.Bold, colors.Reset)
	
	fmt.Printf("%s%sUSAGE:%s\n", colors.Bold, colors.FgMagenta, colors.Reset)
	fmt.Println("  license <command> [arguments] [flags]")
	
	fmt.Printf("\n%s%sCOMMANDS:%s\n", colors.Bold, colors.FgMagenta, colors.Reset)
	fmt.Printf("  %sadd%s <license-name>     Add a license to current directory\n", colors.FgGreen, colors.Reset)
	fmt.Printf("  %slist%s                   List all available licenses\n", colors.FgGreen, colors.Reset)
	fmt.Printf("  %sget%s                    Display current LICENSE file\n", colors.FgGreen, colors.Reset)
	fmt.Printf("  %shelp%s                   Show this help message\n", colors.FgGreen, colors.Reset)
	
	fmt.Printf("\n%s%sFLAGS:%s\n", colors.Bold, colors.FgMagenta, colors.Reset)
	fmt.Printf("  %s--author, -a%s <name>   Specify author name (default: git user.name)\n", colors.FgCyan, colors.Reset)
	fmt.Printf("  %s--year, -y%s <year>     Specify year (default: current year)\n", colors.FgCyan, colors.Reset)
	fmt.Printf("  %s--help, -h%s            Show help message\n", colors.FgCyan, colors.Reset)
	
	fmt.Printf("\n%s%sEXAMPLES:%s\n", colors.Bold, colors.FgMagenta, colors.Reset)
	fmt.Printf("%s", colors.FgYellow)
	fmt.Println("  license add mit")
	fmt.Println("  license add apache-2.0 --author \"John Doe\" --year 2024")
	fmt.Println("  license list")
	fmt.Println("  license get")
	fmt.Printf("%s", colors.Reset)
	
	fmt.Printf("\n%s%s%sNOTE:%s\n", colors.Underline, colors.Bold, colors.FgRed, colors.Reset)
	fmt.Println("  ● License names are case-insensitive.")
	fmt.Println("  ● On first run, you'll be prompted to download licenses from remote.")
}

func version() {
	fmt.Printf(" %slicense-cli:%s", colors.Bold, colors.Reset)
	fmt.Printf("        %s1.2.1%s\n", colors.FgGreen, colors.Reset)
	fmt.Printf(" %sAuthor:%s", colors.Bold, colors.Reset)
	fmt.Printf("             %sArchit%s\n", colors.FgGreen, colors.Reset)
	fmt.Printf(" %sLicense:%s", colors.Bold, colors.Reset)
	fmt.Printf("            %sGPL-3.0%s\n", colors.FgGreen, colors.Reset)
	fmt.Printf(" %sSource Code:%s", colors.Bold, colors.Reset)
	fmt.Printf("        %s%shttps://github.com/architmishra-15/license-cli%s\n", colors.FgGreen, colors.Underline, colors.Reset)
}
