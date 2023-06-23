package config

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"text/template"
	"time"
)

func PostConfigure() {

	populatedStates, err := getConfiguredStates()
	if err != nil {
		fmt.Println("Error populating states")
		fmt.Println(err)
		os.Exit(1)
	}

	APP_STATES = populatedStates

	initSSHDIR()
}

func initSSHDIR() {

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		SSHDIR = os.Getenv(HOMEDRIVE) + os.Getenv(HOMEPATH)
	case "linux":
		SSHDIR = os.Getenv(HOME)
	}

}

func getConfiguredStates() (*States, error) {

	err := checkStatesFile()
	if err != nil {
		return nil, err
	}

	// Create states structure
	states := &States{}

	// Open appConfig file
	file, err := os.Open(APP_CONFIG.TemplatesDirectory + "/states.yml")
	CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&states)
	CheckFatal(err)

	return states, err
}

func checkStatesFile() error {

	if _, err := os.Stat(APP_CONFIG.TemplatesDirectory + "/states.yml"); os.IsNotExist(err) {
		fmt.Println("templateFileLocation: " + APP_CONFIG.TemplatesDirectory)
		fmt.Println("States file does not exist... ")
		return err
	}
	return nil
}

func CheckAndReportOnRepositoryState() bool {

	err := CheckConfigurationFilePresence()
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
	CheckFatal(err)

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
	repository, worktree, err := stage()
	err = commit(err, worktree)
	pushToOrigin(err, repository)
	//FetchTemplateDirectory()

}

func colateInitialConfiguration() {

	// Collect information from user on where the rfd repo will be created.
	repositoryRoot, templatesDirectory, keyType, userName, organisation := getConfigurationInfoFromUser()

	// Write the configuration file
	writeConfigFile(repositoryRoot, templatesDirectory, keyType, userName, organisation)

	// Configure the repository
	Configure()

	// Write the template directory
	_, err := WriteTemplates()
	if err != nil {
		log.Fatal(err)
	}

	PostConfigure()
}

func create0001Rfd() {

	var fileExists = Exists(GetRFDDirectory("0001") + PATH_SEPARATOR + "readme.md")

	if fileExists {

		response := GetUserInput("File exists. Overwrite (y/N)?")
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

func pushToOrigin(err error, repository *git.Repository) {
	Logger.TraceLog("Pushing to origin ...")
	err = PushToOrigin(repository)
	CheckFatal(err)
	Logger.TraceLog("Pushed to origin")
}

func commit(err error, worktree *git.Worktree) error {
	Logger.TraceLog("Committing ...")
	_, err = worktree.Commit("Initialising repository", &git.CommitOptions{
		All: true,
	})
	CheckFatal(err)
	return err
}

func stage() (*git.Repository, *git.Worktree, error) {
	// Stage and commit
	Logger.TraceLog("Staging ...")
	repository, err := git.PlainOpen(".")
	CheckFatal(err)
	worktree, err := repository.Worktree()
	CheckFatal(err)

	_, err = worktree.Add("0001" + PATH_SEPARATOR)
	CheckFatal(err)
	_, err = worktree.Add("0001" + PATH_SEPARATOR + "readme.md")
	CheckFatal(err)
	_, err = worktree.Add("readme.md")
	CheckFatal(err)
	return repository, worktree, err
}

func getConfigurationInfoFromUser() (string, string, string, string, string) {
	// Default to the current directory.

	repositoryRoot := GetUserInput("Enter the path to the directory where you want to create the rfd repository (default: current directory):")
	if repositoryRoot == "" {
		// if the repository root is empty, then use the working directory
		repositoryRoot, _ = os.Getwd()
	}

	// Check to see if the directory exists,and if not, exit.
	if !Exists(repositoryRoot) {
		Logger.TraceLog("The directory " + repositoryRoot + " does not exist.")
		os.Exit(1)
	}

	fmt.Println("Using repository root: " + repositoryRoot)

	templatesDirectory := repositoryRoot + PATH_SEPARATOR + "template"

	fmt.Println("Using templates directory: " + templatesDirectory)

	RSA_OR_DSA := GetUserInput("Enter the type of SSH key you are using (RSA/dsa):")
	if RSA_OR_DSA == "" {
		// If it's empty, deault to RSA
		RSA_OR_DSA = "RSA"
	}

	if RSA_OR_DSA != "RSA" && RSA_OR_DSA != "rsa" && RSA_OR_DSA != "DSA" && RSA_OR_DSA != "dsa" {
		Logger.TraceLog("Invalid key type. Exiting.")
		os.Exit(1)
	}

	keyType := "id_rsa"
	if RSA_OR_DSA == "DSA" || RSA_OR_DSA == "dsa" {
		keyType = "id_ed25519"
	}

	fmt.Println("Using " + keyType + " key type.")

	// Get the name of the first user
	userName := GetUserInput("Enter the name of the first user (default: the current user name):")
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
	organisation := GetUserInput("Enter the name of the organisation (default: MyOrg):")
	if organisation == "" {
		organisation = "MyOrg"
	}

	fmt.Println("Using " + organisation + " as the organisation.")
	return repositoryRoot, templatesDirectory, keyType, userName, organisation
}

func writeConfigFile(repositoryRoot string, templatesDirectory string, keyType string, userName string, organisation string) {

	APP_CONFIG = &Configuration{
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
	yamlData, err := yaml.Marshal(APP_CONFIG)

	if err != nil {
		log.Fatalf("Error while Marshaling. %v", err)
	}

	fileName := "config.yml"
	err = os.WriteFile(fileName, yamlData, 0644)
	if err != nil {
		log.Fatal("Unable to write data into the file", err)
	}
}

func initReadme() {

	formattedRFDNumber := "0001"
	title := "The " + APP_CONFIG.Organisation + " Request for Discussion Process"
	authors := APP_CONFIG.InitialAuthor
	state := "discussion"
	link := ""

	readmeFile := GetRFDDirectory(formattedRFDNumber) + PATH_SEPARATOR + "readme.md"

	if Exists(GetRFDDirectory(formattedRFDNumber)) {

		if Exists(readmeFile) {
			err := os.Remove(readmeFile)
			CheckFatal(err)
		}

		err := os.Remove(GetRFDDirectory(formattedRFDNumber))
		CheckFatal(err)

	}

	CreateReadme(&RFDMetadata{
		formattedRFDNumber,
		title,
		authors,
		state,
		link,
		APP_STATES.RFDStates,
	}, APP_CONFIG.Get001ReadmeFileLocation())

	CopyToRoot(readmeFile, "readme.md", true)

}

func CreateReadme(metadata *RFDMetadata, tmplate string) (error, *os.File) {

	Logger.TraceLog("Creating placeholder readme file, and adding to repository")
	// Create readme.md file with template @ template/readme.md

	Logger.TraceLog("Template:" + tmplate)

	bTemplate, err := os.ReadFile(tmplate)
	CheckFatal(err)
	sTemplate := string(bTemplate)
	tmpl, err := template.New("test").Parse(sTemplate)
	CheckFatal(err)

	// Create local directory

	err = os.Mkdir(GetRFDDirectory(metadata.RFDID), 0755)
	CheckFatal(err)

	// Write out new readme.md to nnnn/readme.md
	// Status on readme.md will be set to "prediscussion"
	fReadme, err := os.Create(GetRFDDirectory(metadata.RFDID) + PATH_SEPARATOR + "readme.md")
	CheckFatal(err)
	defer fReadme.Close()

	err = tmpl.Execute(fReadme, metadata)
	return err, fReadme
}

func printCancelled() {
	println("Operation cancelled.")
}
