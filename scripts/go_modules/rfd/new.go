package main

import (
	"bufio"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/go-git/go-git/v5/storage/memory"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
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

	newRFDNumber := getMaxRFDNumber() + 1
	logger.traceLog("New RFD Number: " + strconv.Itoa(newRFDNumber))

	title := getUserInput("Enter title of RFD: ")
	authors := getUserInput("Enter authors, comma delimited: ")

	err := createRFD(newRFDNumber)
	CheckFatal(err)

	fmt.Println("Title: " + title)
	fmt.Println("Authors: " + authors)
	fmt.Println("RFD ID: " + strconv.Itoa(newRFDNumber))

}

func getMaxRFDNumber() int {

	err, maxRFDBranchId := getMaxBranchRFD()
	CheckFatal(err)
	err, maxRFDDirId := getMaxDirRFD()
	CheckFatal(err)
	err, maxRemoteRFDBranchId := getMaxRemoteBranchRFD()
	CheckFatal(err)

	maxRFDId := maxRFDBranchId
	if maxRFDDirId > maxRFDBranchId {
		maxRFDId = maxRFDDirId
	}
	if maxRemoteRFDBranchId > maxRFDId {
		maxRFDId = maxRemoteRFDBranchId
	}

	return maxRFDDirId
}

func getMaxBranchRFD() (error, int) {

	// Todo - consider using a configuration file to define the path
	r, err := git.PlainOpen(".")
	CheckFatal(err)

	// ... retrieving the branches
	branches, err := r.Branches()
	CheckFatal(err)

	var maxRFDId int = 0

	branches.ForEach(func(p *plumbing.Reference) error {
		rName := p.Name()
		name := rName.String()

		// Remove the first 11 characters as they refer to (I think) internal git
		// identifiers.
		sId := name[11:]

		// A valid branch id is nnnn, e.g. 0007
		entryIsBranchID, err := isMatchRFDId(sId)
		CheckFatal(err)

		if entryIsBranchID {
			rfdId, err := strconv.Atoi(sId)
			if err == nil {
				if rfdId > maxRFDId {
					maxRFDId = rfdId
				}
			}
		}

		return nil
	})
	return err, maxRFDId
}

func getMaxDirRFD() (error, int) {

	var maxRFDId int = 0

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := isMatchRFDId(entry.Name())
		CheckFatal(err)

		if entryIsBranchID {

			sId := entry.Name()

			if entry.IsDir() {

				rfdId, err := strconv.Atoi(sId)
				if err == nil {
					if rfdId > maxRFDId {
						maxRFDId = rfdId
					}
				}

			}

		}

	}

	return nil, maxRFDId
}

func getMaxRemoteBranchRFD() (error, int) {

	// Todo - consider using a configuration file to define the path
	//r, err := git.PlainOpen(".")
	//CheckFatal(err)

	url := "git@github.com:redazzo/rfd.git"
	var publicKey *ssh.PublicKeys
	sshPath := os.Getenv("HOME") + "/.ssh/id_rsa.pub"
	sshKey, _ := ioutil.ReadFile(sshPath)

	print(string(sshKey))
	publicKey, keyError := ssh.NewPublicKeysFromFile("redazzo@dev01", sshPath, "")
	CheckFatal(keyError)

	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
		Auth:     publicKey,
	})
	CheckFatal(err)

	//err = r.Fetch(&git.FetchOptions{
	//	Auth: publicKey,
	//})
	//CheckFatal(err)

	remote, err := r.Remote("origin")
	CheckFatal(err)
	refList, err := remote.List(&git.ListOptions{})
	CheckFatal(err)
	refPrefix := "refs/heads/"
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]
		fmt.Println(branchName)
	}

	return err, 0
}

/*func remoteBranches(s storer.ReferenceStorer) (storer.ReferenceIter, error) {
	refs, err := s.IterReferences()
	if err != nil {
		return nil, err
	}

	return storer.NewReferenceFilteredIter(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	}, refs), nil
}*/

func getUserInput(txt string) string {
	fmt.Println(txt)
	reader := bufio.NewReader(os.Stdin)
	responseTxt, err := reader.ReadString('\n')
	CheckFatal(err)
	return responseTxt
}

func createRFD(rfdNumber int) error {

	return nil //errors.New("test error")
}
