package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/user"
	"strings"
	"time"
)

/**
Initialising a repo steps:

0. Colate initial configuration information from user
1. Create a branch 0001
2. Create a directory 0001
3. Copy template/0001/readme.md into <rfd root>/0001 directory
4. Copy template /0001/readme.md into <rfd root>.

*/

func initRepo() {

	// Colate initial configuration information from user
	colateInitialConfiguration()

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

func colateInitialConfiguration() {

	// Collect information from user on where the rfd repo will be created.
	// Default to the current directory.

	repositoryRoot := getUserInput("Enter the path to the directory where you want to create the rfd repository (default: current directory):")
	if repositoryRoot == "" {
		// if the repository root is empty, then use the working directory
		repositoryRoot, _ = os.Getwd()
	}

	// Check to see if the directory exists,and if not, exit.
	if !exists(repositoryRoot) {
		logger.traceLog("The directory " + repositoryRoot + " does not exist.")
		os.Exit(1)
	}

	fmt.Println("Using repository root: " + repositoryRoot)

	// Get template directory from user
	templatesDirectory := getUserInput("Enter the path to the directory where the rfd templates are located (default: <current directory>/template):")
	if templatesDirectory == "" {
		// if the repository root is empty, then use the working directory
		templatesDirectory, _ = os.Getwd()
		templatesDirectory = templatesDirectory + sPathseparator + "template"
	}

	// Check to see if the directory exists,and if not, exit.
	if !exists(templatesDirectory) {
		logger.traceLog("The directory " + templatesDirectory + " does not exist.")
		os.Exit(1)
	}

	fmt.Println("Using templates directory: " + templatesDirectory)

	RSA_OR_DSA := getUserInput("Enter the type of SSH key you are using (RSA/dsa):")
	if RSA_OR_DSA == "" {
		// If it's empty, deault to RSA
		RSA_OR_DSA = "RSA"
	}

	if RSA_OR_DSA != "RSA" && RSA_OR_DSA != "rsa" && RSA_OR_DSA != "DSA" && RSA_OR_DSA != "dsa" {
		logger.traceLog("Invalid key type. Exiting.")
		os.Exit(1)
	}

	keyType := "id_rsa"
	if RSA_OR_DSA == "DSA" || RSA_OR_DSA == "dsa" {
		keyType = "id_ed25519"
	}

	fmt.Println("Using " + keyType + " key type.")

	// Get the name of the first user
	userName := getUserInput("Enter the name of the first user (default: the current user name):")
	if userName == "" {
		// If it's empty, deault to the current user
		user, err := user.Current()
		if err != nil {
			log.Fatalf(err.Error())
		}

		userName = user.Username
	}

	fmt.Println("Using " + userName + " as the first author.")

	// Get the name of the organisation
	organisation := getUserInput("Enter the name of the organisation (default: MyOrg):")
	if organisation == "" {
		organisation = "MyOrg"
	}

	fmt.Println("Using " + organisation + " as the organisation.")

	// Write the configuration file
	writeConfigFile(repositoryRoot, templatesDirectory, keyType, userName, organisation)

	// Configure the repo
	configure()
	postConfigure()

}

func writeConfigFile(repositoryRoot string, templatesDirectory string, keyType string, userName string, organisation string) {

	appConfig = &configuration{
		RootDirectory:      repositoryRoot,
		TemplatesDirectory: templatesDirectory,
		PrivateKeyFileName: keyType,
		InitialAuthor:      userName,
		Organisation:       organisation,
		InstigationDate:    time.Now().Format(time.DateOnly),
	}

	// Write the configuration file to the root directory in YAML format
	// Open appConfig file
	//file, err := os.Open("./config.yml")
	//CheckFatal(err)

	//defer file.Close()

	// Write the configuration file
	yamlData, err := yaml.Marshal(appConfig)

	if err != nil {
		log.Fatal("Error while Marshaling. %v", err)
	}

	fileName := "config.yml"
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		log.Fatal("Unable to write data into the file", err)
	}
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
	authors := appConfig.InitialAuthor
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
		appStates.RFDStates,
	}, newRepoTemplateFileLocation)

	copyToRoot(readmeFile, "readme.md", true)

}

func printCancelled() {
	println("Operation cancelled.")
}
