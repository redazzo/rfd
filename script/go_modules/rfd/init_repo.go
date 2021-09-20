package main

import (
	"github.com/go-git/go-git/v5"
	"os"
)

/**
Initialising a repo steps:

1. Create a branch 0001
2. Create a directory 0001, and checkout that branch
3. Copy template/0001/readme.md into <rfd root>/0001 directory
4. Copy template /0001/readme.md into <rfd root>.
5. Stage, commit, push to remote, and update upstream tracking
6.
*/

func initRepo() {
	create0001Rfd()
}

func create0001Rfd() error {

	// Format the number to match nnnn
	formattedRFDNumber := "0001"

	// Branch, write the readme file, stage, commit, push, and set upstream

	// Create a branch named as per "nnnn"
	err, r, w, _ := createBranch(formattedRFDNumber)
	CheckFatal(err)

	err, _ = copyReadme()
	CheckFatal(err)

	// Stage and commit
	_, err = w.Add(appConfig.RFDRelativeDirectory + sPathseparator + formattedRFDNumber + sPathseparator + "readme.md")
	CheckFatal(err)

	logger.traceLog("Committing ...")
	_, err = w.Commit("Earmark branch", &git.CommitOptions{
		All: true,
	})
	CheckFatal(err)

	// Push to origin and set upstream
	err = pushToOrigin(r)
	CheckFatal(err)

	err = setUpstream(r, formattedRFDNumber)

	return err
}

func copyReadme() (error, *os.File) {

	return nil, nil
}
