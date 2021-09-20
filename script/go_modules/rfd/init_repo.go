package main

import (
	"os"
	"strings"
)

/**
Initialising a repo steps:

1. Create a branch 0001
2. Create a directory 0001
3. Copy template/0001/readme.md into <rfd root>/0001 directory
4. Copy template /0001/readme.md into <rfd root>.

*/

func initRepo() {
	create0001Rfd()
}

func create0001Rfd() {

	var fileExists = isFileExists(getRFDDirectory("0001") + sPathseparator + "readme.md")

	if fileExists {

		response := getUserInput("File exists. Overwrite (y/N)?")
		response = strings.ToUpper(response)

		switch response {
		case "N":
			printCancelled()

		case "NO":
			printCancelled()

		case "Y":
			initReadme()

		case "YES":
			initReadme()

		default:
			printCancelled()

		}


	} else {

		initReadme()

	}
}

func initReadme() {

	// Format the number to match nnnn
	formattedRFDNumber := "0001"
	title := "The " + appConfig.Organisation + " Request for Discussion Process"
	authors := "Gerry Kessell-Haak"
	state := "discussion"
	link := ""

	readmeFile := getRFDDirectory(formattedRFDNumber) + sPathseparator + "readme.md"

	if isFileExists(getRFDDirectory(formattedRFDNumber)) {

		if isFileExists(readmeFile) {
			err := os.Remove(readmeFile)
			CheckFatal(err)
		}

		err := os.Remove(getRFDDirectory(formattedRFDNumber))
		CheckFatal(err)

	}


	createReadme(&RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,

	}, newRepoTemplateFileLocation)
}

func printCancelled() {
	println("Operation cancelled.")
}

