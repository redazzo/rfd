package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	localConfig "github.com/redazzo/rfd/cmd/rfd/config"
	"github.com/redazzo/rfd/cmd/rfd/global"
	"github.com/redazzo/rfd/cmd/rfd/util"
	"strconv"
	"strings"
)

/*

An RFD can be in one of two branch states:

1. Newly created, or still being updated and not yet ready for mainlining into the trunk. In this case
   there won't yet be a separate RFD directory in the trunk. Instead there will be a branch called nnnn,
   where nnnn is a 4-digit number. On this branch there will be a directory called /nnnn, and a readme.md file
   located at /nnnn/readme.md

2. Mainlined - the nnnn branch will have been merged into the trunk (master), with status set to accepted (or beyond e.g. committed).

To create a new RFD:
1. Fetch all local branches that match nnnn naming format, and keep a record of the greatest as L
2. Fetch all remote branches that match nnnn naming format, and keep a record of the greatest as R
3. Fetch all directories on the trunk that match nnnn naming format, and keep a record of the greatest as D
4. Find max(L, R, D), and create a branch max(L, R, D)+1 --> mmmm
5. Create a readme.md file --> mmmm\readme.md
6. Stage, commit, push to remote, and update upstream tracking

*/

func new() {
	util.Logger.TraceLog("Creating new RFD")

	newRFDNumber := getMaxRFDNumber() + 1
	util.Logger.TraceLog("New RFD Number: " + strconv.Itoa(newRFDNumber))

	title := util.GetUserInput("Enter title of RFD: ")
	authors := util.GetUserInput("Enter authors, comma delimited: ")

	fmt.Println("Title: " + title)
	fmt.Println("Authors: " + authors)
	fmt.Println("RFD ID: " + strconv.Itoa(newRFDNumber))

	defaultStatus := getDefaultStatus()

	err := createRFD(newRFDNumber, title, authors, defaultStatus, "")

	util.CheckFatal(err)

}

func getDefaultStatus() string {

	var result string = "ERROR"

	for _, state := range global.APP_STATES.RFDStates {
		for _, m := range state {
			if m["id"] == "1" {
				result = m["name"]
			}
		}
	}

	return result
}

func createRFD(rfdNumber int, title string, authors string, state string, link string) error {

	// Format the number to match nnnn
	formattedRFDNumber := formatToNNNN(rfdNumber)

	// Branch, write the readme file, stage, commit, push, and set upstream

	// Create a branch named as per "nnnn"
	err, r, w, _ := createBranch(formattedRFDNumber)
	util.CheckFatal(err)

	err, _ = localConfig.CreateReadme(&global.RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,
		nil,
	}, global.TEMPLATE_FILE_LOCATION)

	// Update index
	Index()

	util.CheckFatal(err)

	// Stage and commit
	util.Logger.TraceLog("Staging ...")

	_, err = w.Add(formattedRFDNumber + global.PATH_SEPARATOR)
	util.CheckFatal(err)
	_, err = w.Add(formattedRFDNumber + global.PATH_SEPARATOR + "readme.md")
	util.CheckFatal(err)
	_, err = w.Add("index.md")
	util.CheckFatal(err)

	util.Logger.TraceLog("Committing ...")
	_, err = w.Commit("Earmark branch", &git.CommitOptions{
		All: true,
	})
	util.CheckFatal(err)

	// Push to origin and set upstream
	util.Logger.TraceLog("Pushing to origin ...")
	err = util.PushToOrigin(r)
	util.CheckFatal(err)
	util.Logger.TraceLog("Pushed to origin")

	util.Logger.TraceLog("Setting upstream ...")
	err = setUpstream(r, formattedRFDNumber)
	util.Logger.TraceLog("Upstream set to " + formattedRFDNumber)

	return err
}

func setUpstream(r *git.Repository, formattedRFDNumber string) error {

	r, err := git.PlainOpen(".")
	util.CheckFatal(err)
	currentConfig, err := r.Config()
	util.CheckFatal(err)

	branches := currentConfig.Branches

	referenceName := plumbing.NewBranchReferenceName(formattedRFDNumber)

	newBranch := &config.Branch{
		Name:   formattedRFDNumber,
		Remote: "origin",
		Merge:  referenceName,
	}

	branches[formattedRFDNumber] = newBranch

	err = r.Storer.SetConfig(currentConfig)
	util.CheckFatal(err)

	return err
}

func formatToNNNN(rfdNumber int) string {
	sRfdNumber := strconv.Itoa(rfdNumber)

	strLength := len(sRfdNumber)
	for strLength < 4 {
		sRfdNumber = "0" + sRfdNumber
		strLength++
	}
	return sRfdNumber
}

func undoCreateReadme() {

	fmt.Println("Called undoCreateReadme - but I'm empty :(")

}

func createBranch(rfdNumber string) (error, *git.Repository, *git.Worktree, *plumbing.Reference) {

	r, err := git.PlainOpen(".")
	util.CheckFatal(err)

	headRef, err := r.Head()
	util.CheckFatal(err)

	// Create a new plumbing.HashReference object with the name of the branch
	// and the hash from the HEAD. The reference name should be a full reference
	// name and not an abbreviated one, as is used on the git cli.
	//
	// For tags we should use `refs/tags/%s` instead of `refs/heads/%s` used
	// for branches.
	rfName := "refs/heads/" + rfdNumber
	ref := plumbing.NewHashReference(plumbing.ReferenceName(rfName), headRef.Hash())

	// The created reference is saved in the storage.
	err = r.Storer.SetReference(ref)
	util.CheckFatal(err)

	w, err := r.Worktree()
	util.CheckFatal(err)

	// ... checking out to commit
	util.Logger.TraceLog("checking out")

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref.Name(),
		Keep:   true,
	})
	util.CheckFatal(err)

	return err, r, w, ref
}

func undoCreateBranch() {
	fmt.Println("Called undoCreateBranch - but I'm empty :(")

}

func getMaxRFDNumber() int {

	err, maxRFDBranchId := getMaxBranchId()
	util.CheckFatal(err)
	util.Logger.TraceLog("Local branch max id: " + strconv.Itoa(maxRFDBranchId))

	err, maxRFDDirId := getMaxDirId()
	util.CheckFatal(err)
	util.Logger.TraceLog("Directory branch max id: " + strconv.Itoa(maxRFDDirId))

	err, maxRemoteRFDBranchId := getMaxRemoteBranchId()
	util.CheckFatal(err)
	util.Logger.TraceLog("Remote branch max id: " + strconv.Itoa(maxRemoteRFDBranchId))

	maxRFDId := maxRFDBranchId
	if maxRFDDirId > maxRFDBranchId {
		maxRFDId = maxRFDDirId
	}
	if maxRemoteRFDBranchId > maxRFDId {
		maxRFDId = maxRemoteRFDBranchId
	}

	return maxRFDId
}

func getMaxBranchId() (error, int) {

	r, err := git.PlainOpen(".")
	util.CheckFatal(err)

	// ... retrieving the branches
	branches, err := r.Branches()
	util.CheckFatal(err)

	var maxRFDId = 0

	branches.ForEach(func(p *plumbing.Reference) error {
		rName := p.Name()
		name := rName.String()

		// Remove the "refs/heads/" identifiers.
		sId := name[11:]

		// A valid branch id is nnnn, e.g. 0007
		entryIsBranchID, err := util.IsRFDIDFormat(sId)
		util.CheckFatal(err)

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

func getMaxDirId() (error, int) {

	var maxRFDId = 0

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := util.IsRFDIDFormat(entry.Name())
		util.CheckFatal(err)

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

func getMaxRemoteBranchId() (error, int) {

	var maxRemoteBranchId = 0

	r, err := git.PlainOpen(".")
	util.CheckFatal(err)

	publicKey, err := util.GetPublicKey()

	remote, err := r.Remote("origin")
	util.CheckFatal(err)
	refList, err := remote.List(&git.ListOptions{
		Auth: publicKey,
	})
	util.CheckFatal(err)

	refPrefix := "refs/heads/"
	for _, ref := range refList {

		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]

		entryIsBranchID, err := util.IsRFDIDFormat(branchName)
		util.CheckFatal(err)

		if entryIsBranchID {
			entryId, err := strconv.Atoi(branchName)
			if (entryId > maxRemoteBranchId) && err == nil {
				maxRemoteBranchId = entryId
			}
		}
	}

	return err, maxRemoteBranchId
}
