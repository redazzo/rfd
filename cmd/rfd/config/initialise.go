package config

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/redazzo/rfd/cmd/rfd/global"
	"github.com/redazzo/rfd/cmd/rfd/util"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"text/template"
	"time"
)

func PreConfigure() {

	global.PATH_SEPARATOR = string(os.PathSeparator)

}

func Configure() {

	err := checkConfigurationFile()
	util.CheckFatal(err)
	global.APP_CONFIG, err = populateConfig()
	if err != nil {
		fmt.Println("Error populating configuration")
		fmt.Println(err)
		os.Exit(1)
	}

}

func PostConfigure() {

	initFileLocations()

	populatedStates, err := getConfiguredStates()
	if err != nil {
		fmt.Println("Error populating states")
		fmt.Println(err)
		os.Exit(1)
	}

	global.APP_STATES = populatedStates

	initSSHDIR()
}

func initFileLocations() {
	initTemplateFileLocation()
	initNewRepoTemplateFileLocation()
}

func initSSHDIR() {

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		global.SSHDIR = os.Getenv(global.HOMEDRIVE) + os.Getenv(global.HOMEPATH)
	case "linux":
		global.SSHDIR = os.Getenv(global.HOME)
	}

}

func initTemplateFileLocation() {
	global.TEMPLATE_FILE_LOCATION = global.APP_CONFIG.TemplatesDirectory + global.PATH_SEPARATOR + "readme.md"
}

func initNewRepoTemplateFileLocation() {
	global.REPO_TEMPLATE_FILE_LOCATION = global.APP_CONFIG.TemplatesDirectory + global.PATH_SEPARATOR + "0001" + global.PATH_SEPARATOR + "readme.md"
}

func populateConfig() (*global.Configuration, error) {

	err := checkConfigurationFile()
	if err != nil {
		return nil, err
	}

	// Create appConfig structure
	config := &global.Configuration{}

	// Open appConfig file
	file, err := os.Open("./config.yml")
	util.CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&config)
	util.CheckFatal(err)

	return config, err
}

func getConfiguredStates() (*global.States, error) {

	err := checkStatesFile()
	if err != nil {
		return nil, err
	}

	// Create states structure
	states := &global.States{}

	// Open appConfig file
	file, err := os.Open(global.APP_CONFIG.TemplatesDirectory + "/states.yml")
	util.CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&states)
	util.CheckFatal(err)

	return states, err
}

func checkStatesFile() error {

	if _, err := os.Stat(global.APP_CONFIG.TemplatesDirectory + "/states.yml"); os.IsNotExist(err) {
		fmt.Println("templateFileLocation: " + global.APP_CONFIG.TemplatesDirectory)
		fmt.Println("States file does not exist... ")
		return err
	}
	return nil
}

func CheckAndReportOnRepositoryState() bool {

	err := checkConfigurationFile()
	if err != nil {
		return false
	}

	var fileStatusMapping = map[git.StatusCode]string{
		git.Unmodified:         "Unmodified",
		git.Untracked:          "Untracked",
		git.Modified:           "==== Modified ====",
		git.Added:              "==== Added ====",
		git.Deleted:            "==== Deleted ====",
		git.Renamed:            "==== Renamed ====",
		git.Copied:             "==== Copied ====",
		git.UpdatedButUnmerged: "==== Updated ====",
	}

	// Check to ensure git status is clean.
	r, err := git.PlainOpen(".")
	util.CheckFatal(err)

	w, err := r.Worktree()
	status, err := w.Status()

	fmt.Println()
	if status.IsClean() {

		fmt.Println("Nothing to commit, working tree clean.")

	} else {

		fmt.Println("There are changes present in the repository.")
		fmt.Println()

		for s := range status {
			fileStatus := status.File(s)
			fmt.Println(s + " is " + fileStatusMapping[fileStatus.Worktree])

		}
		fmt.Println()

	}

	return status.IsClean()
}

func checkConfigurationFile() error {
	// Check to ensure there is a config file present
	_, err := os.Stat("./config.yml")
	if os.IsNotExist(err) {
		fmt.Print(
			"\n  There doesn't appear to be a configuration file present.\n" +
				"  Either run 'rfd init', or if you have, make sure you are\n" +
				"  in the root directory of your rfd repository.\n\n")
		return err
	}
	return nil
}

/**
Initialising a repo steps:

0. Colate initial configuration information from user
1. Create a branch 0001
2. Create a directory 0001
3. Copy template/0001/readme.md into <rfd root>/0001 directory
4. Copy template /0001/readme.md into <rfd root>.

*/

func InitialiseRepo() {

	// Colate initial configuration information from user
	colateInitialConfiguration()

	create0001Rfd()

	// Stage and commit
	util.Logger.TraceLog("Staging ...")
	r, err := git.PlainOpen(".")
	util.CheckFatal(err)
	w, err := r.Worktree()
	util.CheckFatal(err)

	_, err = w.Add("0001" + global.PATH_SEPARATOR)
	util.CheckFatal(err)
	_, err = w.Add("0001" + global.PATH_SEPARATOR + "readme.md")
	util.CheckFatal(err)
	_, err = w.Add("readme.md")
	util.CheckFatal(err)

	util.Logger.TraceLog("Committing ...")
	_, err = w.Commit("Initialising repository", &git.CommitOptions{
		All: true,
	})
	util.CheckFatal(err)

	// Push to origin
	util.Logger.TraceLog("Pushing to origin ...")
	err = util.PushToOrigin(r)
	util.CheckFatal(err)
	util.Logger.TraceLog("Pushed to origin")

}

func colateInitialConfiguration() {

	// Collect information from user on where the rfd repo will be created.
	// Default to the current directory.

	repositoryRoot := util.GetUserInput("Enter the path to the directory where you want to create the rfd repository (default: current directory):")
	if repositoryRoot == "" {
		// if the repository root is empty, then use the working directory
		repositoryRoot, _ = os.Getwd()
	}

	// Check to see if the directory exists,and if not, exit.
	if !util.Exists(repositoryRoot) {
		util.Logger.TraceLog("The directory " + repositoryRoot + " does not exist.")
		os.Exit(1)
	}

	fmt.Println("Using repository root: " + repositoryRoot)

	// Get template directory from user
	templatesDirectory := util.GetUserInput("Enter the path to the directory where the rfd templates are located (default: <current directory>/template):")
	if templatesDirectory == "" {
		// if the repository root is empty, then use the working directory
		templatesDirectory, _ = os.Getwd()
		templatesDirectory = templatesDirectory + global.PATH_SEPARATOR + "template"
	}

	// Check to see if the directory exists,and if not, exit.
	if !util.Exists(templatesDirectory) {
		util.Logger.TraceLog("The directory " + templatesDirectory + " does not exist.")
		os.Exit(1)
	}

	fmt.Println("Using templates directory: " + templatesDirectory)

	RSA_OR_DSA := util.GetUserInput("Enter the type of SSH key you are using (RSA/dsa):")
	if RSA_OR_DSA == "" {
		// If it's empty, deault to RSA
		RSA_OR_DSA = "RSA"
	}

	if RSA_OR_DSA != "RSA" && RSA_OR_DSA != "rsa" && RSA_OR_DSA != "DSA" && RSA_OR_DSA != "dsa" {
		util.Logger.TraceLog("Invalid key type. Exiting.")
		os.Exit(1)
	}

	keyType := "id_rsa"
	if RSA_OR_DSA == "DSA" || RSA_OR_DSA == "dsa" {
		keyType = "id_ed25519"
	}

	fmt.Println("Using " + keyType + " key type.")

	// Get the name of the first user
	userName := util.GetUserInput("Enter the name of the first user (default: the current user name):")
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
	organisation := util.GetUserInput("Enter the name of the organisation (default: MyOrg):")
	if organisation == "" {
		organisation = "MyOrg"
	}

	fmt.Println("Using " + organisation + " as the organisation.")

	// Write the configuration file
	writeConfigFile(repositoryRoot, templatesDirectory, keyType, userName, organisation)

	// Configure the repo
	Configure()
	PostConfigure()

}

func writeConfigFile(repositoryRoot string, templatesDirectory string, keyType string, userName string, organisation string) {

	global.APP_CONFIG = &global.Configuration{
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
	yamlData, err := yaml.Marshal(global.APP_CONFIG)

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

	var fileExists = util.Exists(util.GetRFDDirectory("0001") + global.PATH_SEPARATOR + "readme.md")

	if fileExists {

		response := util.GetUserInput("File exists. Overwrite (y/N)?")
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
	title := "The " + global.APP_CONFIG.Organisation + " Request for Discussion Process"
	authors := global.APP_CONFIG.InitialAuthor
	state := "discussion"
	link := ""

	readmeFile := util.GetRFDDirectory(formattedRFDNumber) + global.PATH_SEPARATOR + "readme.md"

	if util.Exists(util.GetRFDDirectory(formattedRFDNumber)) {

		if util.Exists(readmeFile) {
			err := os.Remove(readmeFile)
			util.CheckFatal(err)
		}

		err := os.Remove(util.GetRFDDirectory(formattedRFDNumber))
		util.CheckFatal(err)

	}

	CreateReadme(&global.RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,
		global.APP_STATES.RFDStates,
	}, global.REPO_TEMPLATE_FILE_LOCATION)

	util.CopyToRoot(readmeFile, "readme.md", true)

}

func CreateReadme(metadata *global.RFDMetadata, tmplate string) (error, *os.File) {

	util.Logger.TraceLog("Creating placeholder readme file, and adding to repository")
	// Create readme.md file with template @ template/readme.md

	util.Logger.TraceLog("Template:" + tmplate)

	bTemplate, err := os.ReadFile(tmplate)
	util.CheckFatal(err)
	sTemplate := string(bTemplate)
	tmpl, err := template.New("test").Parse(sTemplate)
	util.CheckFatal(err)

	// Create local directory

	err = os.Mkdir(util.GetRFDDirectory(metadata.RFDID), 0755)
	util.CheckFatal(err)

	// Write out new readme.md to nnnn/readme.md
	// Status on readme.md will be set to "prediscussion"
	fReadme, err := os.Create(util.GetRFDDirectory(metadata.RFDID) + global.PATH_SEPARATOR + "readme.md")
	util.CheckFatal(err)
	defer fReadme.Close()

	err = tmpl.Execute(fReadme, metadata)
	return err, fReadme
}

func printCancelled() {
	println("Operation cancelled.")
}
