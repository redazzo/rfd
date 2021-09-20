package main

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
	title := "The " + appConfig.Organisation + " Request for Discussion Process"
	authors := "Gerry Kessell-Haak"
	state := "discussion"
	link := ""

	// Branch, write the readme file, stage, commit, push, and set upstream

	// Create a branch named as per "nnnn"
	//err, r, w, _ := createBranch(formattedRFDNumber)
	//CheckFatal(err)

	createReadme(&RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,

	}, newRepoTemplateFileLocation)

	return nil
}
