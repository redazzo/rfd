package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"os"
	"strconv"
)

/*
	An RFD can be in one of two branch states:
		1. Newly created, or still being updated and not yet ready for mainlining into the trunk. In this case
		   there won't yet be an separate RFD directory in the trunk. Instead there will be a branch called nnnn,
		   where nnnn is a 4-digit number. On this branch there will be a directory called /nnnn, and a readme.md file
		   located at /nnnn/readme.md
		2. Trunk - the nnnn branch will have been merged into the trunk.

		Algorithm to create a new RFD:
		1. Fetch all branches that match nnnn naming format, and keep a record of the greatest nnnn
		2. Fetch all directories on the trunk that match nnnn naming format, and compare with prior greatest nnnn and keep a record of the highest nnnn.
		3. Create a branch nnnn+1 --> mmmm
		4. Create a readme.md file --> mmmm\readme.md

*/

func NewRFD() {
	logger.traceLog("Creating new RFD")

	rfdNumber := getMaxRFDNumber()

	title := getUserInput("Enter title of RFD: ")
	authors := getUserInput("Enter authors, comma delimited: ")

	err := createRFD(rfdNumber)
	CheckFatal(err)

	fmt.Println("Title: " + title)
	fmt.Println("Authors: " + authors)
	fmt.Println("RFD ID: " + strconv.Itoa(rfdNumber))

}

func getMaxRFDNumber() int {

	// Fetch branches
	//logger.traceLog("git branch")

	r, err := git.PlainOpen(".")
	CheckFatal(err)

	// Length of the HEAD history
	logger.traceLog("git rev-list HEAD --count")

	// ... retrieving the HEAD reference
	ref, err := r.Head()
	CheckFatal(err)

	// ... retrieves the commit history
	cIter, err := r.Log(&git.LogOptions{From: ref.Hash()})
	CheckFatal(err)

	// ... just iterates over the commits
	var cCount int
	err = cIter.ForEach(func(c *object.Commit) error {
		cCount++

		return nil
	})
	CheckFatal(err)

	fmt.Println(cCount)

	//_, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
	//	URL:      "https://github.com/go-git/go-git",
	//	Progress: os.Stdout,
	//})

	CheckFatal(err)

	return 5
}

func getUserInput(txt string) string {
	fmt.Println(txt)
	reader := bufio.NewReader(os.Stdin)
	responseTxt, err := reader.ReadString('\n')
	CheckFatal(err)
	return responseTxt
}

func createRFD(rfdNumber int) error {

	return errors.New("test error")
}
