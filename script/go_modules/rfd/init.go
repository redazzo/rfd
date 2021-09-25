package main

import (
	"github.com/go-git/go-git/v5"
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

	// Stage and commit
	logger.traceLog("Staging ...")
	r, err := git.PlainOpen(".")
	CheckFatal(err)
	w, err := r.Worktree()
	CheckFatal(err)

	tmpPrefix := getPathPrefix()
	_, err = w.Add(tmpPrefix + "0001" + sPathseparator)
	CheckFatal(err)
	_, err = w.Add(tmpPrefix + "0001" + sPathseparator + "readme.md")
	CheckFatal(err)
	_, err = w.Add("readme.md")
	CheckFatal(err)

	logger.traceLog("Committing ...")
	_, err = w.Commit("Initialising repository", &git.CommitOptions{
		All: true,
	})
	CheckFatal(err)

	// Push to origin
	logger.traceLog("Pushing to origin ...")
	err = pushToOrigin(r)
	CheckFatal(err)
	logger.traceLog("Pushed to origin")

}

func create0001Rfd() {

	var fileExists = exists(getRFDDirectory("0001") + sPathseparator + "readme.md")

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

	formattedRFDNumber := "0001"
	title := "The " + appConfig.Organisation + " Request for Discussion Process"
	authors := "Gerry Kessell-Haak"
	state := "discussion"
	link := ""

	readmeFile := getRFDDirectory(formattedRFDNumber) + sPathseparator + "readme.md"

	if exists(getRFDDirectory(formattedRFDNumber)) {

		if exists(readmeFile) {
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
		appConfig.RFDStates,

	}, newRepoTemplateFileLocation)

	copyToRoot(readmeFile, "readme.md", true)

}

func printCancelled() {
	println("Operation cancelled.")
}

