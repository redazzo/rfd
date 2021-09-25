package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"text/template"
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

type RFDMetadata struct {
	RFDID   string
	Title   string
	Authors string
	State   string
	Link    string
	RFDStates []map[string]map[string]string
}

func new() {
	logger.traceLog("Creating new RFD")

	newRFDNumber := getMaxRFDNumber() + 1
	logger.traceLog("New RFD Number: " + strconv.Itoa(newRFDNumber))

	title := getUserInput("Enter title of RFD: ")
	authors := getUserInput("Enter authors, comma delimited: ")

	fmt.Println("Title: " + title)
	fmt.Println("Authors: " + authors)
	fmt.Println("RFD ID: " + strconv.Itoa(newRFDNumber))

	err := createRFD(newRFDNumber, title, authors, "prediscussion", "")

	CheckFatal(err)

}

func createRFD(rfdNumber int, title string, authors string, state string, link string) error {

	// Format the number to match nnnn
	formattedRFDNumber := formatToNNNN(rfdNumber)

	// Branch, write the readme file, stage, commit, push, and set upstream

	// Create a branch named as per "nnnn"
	err, r, w, _ := createBranch(formattedRFDNumber)
	CheckFatal(err)

	err, _ = createReadme(&RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,
		nil,
	}, templateFileLocation)

	CheckFatal(err)

	// Stage and commit
	logger.traceLog("Staging ...")
	tmpPrefix := getPathPrefix()

	_, err = w.Add(tmpPrefix + formattedRFDNumber + sPathseparator)
	_, err = w.Add(tmpPrefix + formattedRFDNumber + sPathseparator + "readme.md")
	CheckFatal(err)

	logger.traceLog("Committing ...")
	_, err = w.Commit("Earmark branch", &git.CommitOptions{
		All: true,
	})
	CheckFatal(err)

	// Push to origin and set upstream
	logger.traceLog("Pushing to origin ...")
	err = pushToOrigin(r)
	CheckFatal(err)
	logger.traceLog("Pushed to origin")

	logger.traceLog("Setting upstream ...")
	err = setUpstream(r, formattedRFDNumber)
	logger.traceLog("Upstream set to " + formattedRFDNumber)

	return err
}

func getPathPrefix() string {
	tmpPrefix := sPathseparator
	if len(strings.TrimSpace(appConfig.RFDRelativeDirectory)) == 0 {
		tmpPrefix = ""
	} else {
		tmpPrefix = appConfig.RFDRelativeDirectory + sPathseparator
	}
	return tmpPrefix
}

func pushToOrigin(r *git.Repository) error {

	publicKey, err := getPublicKey()

	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       publicKey,
		Force:      true,
	})
	CheckFatal(err)

	return err
}

func setUpstream(r *git.Repository, formattedRFDNumber string) error {

	r, err := git.PlainOpen(".")
	CheckFatal(err)
	currentConfig, err := r.Config()
	CheckFatal(err)

	branches := currentConfig.Branches

	referenceName := plumbing.NewBranchReferenceName(formattedRFDNumber)

	newBranch := &config.Branch{
		Name:   formattedRFDNumber,
		Remote: "origin",
		Merge:  referenceName,
	}

	branches[formattedRFDNumber] = newBranch

	err = r.Storer.SetConfig(currentConfig)
	CheckFatal(err)

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

func createReadme(metadata *RFDMetadata, tmplate string) (error, *os.File) {

	logger.traceLog("Creating placeholder readme file, and adding to repository")
	// Create readme.md file with template @ template/readme.md

	logger.traceLog(tmplate)
	bTemplate, err := ioutil.ReadFile(tmplate)
	CheckFatal(err)
	sTemplate := string(bTemplate)
	tmpl, err := template.New("test").Parse(sTemplate)
	CheckFatal(err)

	// Create local directory

	err = os.Mkdir(getRFDDirectory(metadata.RFDID), 0755)
	CheckFatal(err)

	// Write out new readme.md to nnnn/readme.md
	// Status on readme.md will be set to "prediscussion"
	fReadme, err := os.Create(getRFDDirectory(metadata.RFDID) + sPathseparator + "readme.md")
	CheckFatal(err)
	defer fReadme.Close()

	err = tmpl.Execute(fReadme, metadata)
	return err, fReadme
}

func undoCreateReadme() {

	fmt.Println("Called undoCreateReadme - but I'm empty :(")

}

func getRFDDirectory(sRfdNumber string) string {
	return appConfig.RFDRootDirectory + "/" + sRfdNumber
}

func createBranch(rfdNumber string) (error, *git.Repository, *git.Worktree, *plumbing.Reference) {

	r, err := git.PlainOpen(".")
	CheckFatal(err)

	headRef, err := r.Head()
	CheckFatal(err)

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
	CheckFatal(err)

	w, err := r.Worktree()
	CheckFatal(err)

	// ... checking out to commit
	logger.traceLog("checking out")

	err = w.Checkout(&git.CheckoutOptions{
		Branch: ref.Name(),
		Keep:   true,
	})
	CheckFatal(err)

	return err, r, w, ref
}

func undoCreateBranch() {
	fmt.Println("Called undoCreateBranch - but I'm empty :(")

}

func getMaxRFDNumber() int {

	err, maxRFDBranchId := getMaxBranchId()
	CheckFatal(err)
	logger.traceLog("Local branch max id: " + strconv.Itoa(maxRFDBranchId))

	err, maxRFDDirId := getMaxDirId()
	CheckFatal(err)
	logger.traceLog("Directory branch max id: " + strconv.Itoa(maxRFDDirId))

	err, maxRemoteRFDBranchId := getMaxRemoteBranchId()
	CheckFatal(err)
	logger.traceLog("Remote branch max id: " + strconv.Itoa(maxRemoteRFDBranchId))

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
	CheckFatal(err)

	// ... retrieving the branches
	branches, err := r.Branches()
	CheckFatal(err)

	var maxRFDId = 0

	branches.ForEach(func(p *plumbing.Reference) error {
		rName := p.Name()
		name := rName.String()

		// Remove the "refs/heads/" identifiers.
		sId := name[11:]

		// A valid branch id is nnnn, e.g. 0007
		entryIsBranchID, err := isRFDIDFormat(sId)
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

func getMaxDirId() (error, int) {

	var maxRFDId = 0

	entries := getDirectories()
	for _, entry := range entries {

		entryIsBranchID, err := isRFDIDFormat(entry.Name())
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

func getMaxRemoteBranchId() (error, int) {

	var maxRemoteBranchId = 0

	r, err := git.PlainOpen(".")
	CheckFatal(err)

	publicKey, err := getPublicKey()

	remote, err := r.Remote("origin")
	CheckFatal(err)
	refList, err := remote.List(&git.ListOptions{
		Auth: publicKey,
	})
	CheckFatal(err)

	refPrefix := "refs/heads/"
	for _, ref := range refList {

		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		branchName := refName[len(refPrefix):]

		entryIsBranchID, err := isRFDIDFormat(branchName)
		CheckFatal(err)

		if entryIsBranchID {
			entryId, err := strconv.Atoi(branchName)
			if (entryId > maxRemoteBranchId) && err == nil {
				maxRemoteBranchId = entryId
			}
		}
	}

	return err, maxRemoteBranchId
}


