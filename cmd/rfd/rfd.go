package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"os"
	"runtime"
)

const HOME string = "HOME"
const HOMEDRIVE string = "HOMEDRIVE"
const HOMEPATH string = "HOMEPATH"

var sshDir string
var templateFileLocation string
var newRepoTemplateFileLocation string
var sPathseparator string

var logger Trace
var appConfig *configuration
var appStates *states

type configuration struct {
	RootDirectory      string `yaml:"root-directory"`
	TemplatesDirectory string `yaml:"templates-directory"`
	//RFDRelativeDirectory string                         `yaml:"rfd-relative-directory"`
	PrivateKeyFileName string `yaml:"private-key-file-name"`
	InitialAuthor      string `yaml:"initial-author"`
	Organisation       string `yaml:"organisation"`
	InstigationDate    string `yaml:"instigation-date"`
}

type states struct {
	RFDStates []map[string]map[string]string `yaml:"rfd-states"`
}

func preConfigure() {
	sPathseparator = string(os.PathSeparator)
	logger = TraceLog{}

}

func configure() {

	err := checkConfigurationFile()
	CheckFatal(err)
	appConfig, err = populateConfig()
	if err != nil {
		fmt.Println("Error populating configuration")
		fmt.Println(err)
		os.Exit(1)
	}

}

func postConfigure() {

	initFileLocations()

	populatedStates, err := getConfiguredStates()
	if err != nil {
		fmt.Println("Error populating states")
		fmt.Println(err)
		os.Exit(1)
	}

	appStates = populatedStates

	initSSHDIR()
}

func initFileLocations() {
	initTemplateFileLocation()
	initNewRepoTemplateFileLocation()
}

func main() {
	app := createCommandLineApp()
	err := app.Run(os.Args)
	CheckFatal(err)
}

func createCommandLineApp() *cli.App {
	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index, and manage their status.",
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Check environment is suitable to ensure a clean run when creating a new RFD.",
				Action: func(c *cli.Context) error {
					preConfigure()
					configure()
					postConfigure()
					checkAndReportOnRepositoryState()
					return nil
				},
			},
			{
				Name:  "index",
				Usage: "Output the status of all rfd's to index.md in markdown format.",
				Action: func(c *cli.Context) error {
					preConfigure()
					configure()
					postConfigure()
					Index()
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "Create a new rfd",
				Action: func(c *cli.Context) error {
					if checkAndReportOnRepositoryState() {
						preConfigure()
						configure()
						postConfigure()
						new()
					} else {
						fmt.Println("Creating a new RFD creates and switches to new branch. Commit (or otherwise) unstaged and/or uncommitted work first.")
						fmt.Println()
					}
					return nil
				},
			},

			{
				Name:  "init",
				Usage: "Initialise an RFD repository.",
				Action: func(c *cli.Context) error {
					preConfigure()
					initRepo()

					return nil
				},
			},
			{
				Name:  "environment",
				Usage: "Displays configuration settings and relevant operating system environment variables.",
				Action: func(c *cli.Context) error {
					preConfigure()
					configure()
					postConfigure()
					displayEnvironment()
					return nil
				},
			},
			{
				Name:  "merge",
				Usage: "Transitions an RFD's status to Accepted, captures the discussion link from the user, and merges it into the main branch (NOT IMPLEMENTED)",
				Action: func(c *cli.Context) error {
					preConfigure()
					configure()
					postConfigure()
					doMerge()
					return nil
				},
			},
			{
				Name:  "status",
				Usage: "Displays the status of the current RFD (as per branch). (NOT IMPLEMENTED)",
				Action: func(c *cli.Context) error {

					preConfigure()
					configure()
					postConfigure()

					getDefaultStatus()
					return nil
				},
			},
		},
	}

	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[arguments...]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}{{end}}{{if .Copyright }}
COPYRIGHT:
   {{.Copyright}}
   {{end}}{{if .Version}}
VERSION:
   {{.Version}}
   {{end}}
`
	return app
}

func initSSHDIR() {

	operatingSystem := runtime.GOOS
	switch operatingSystem {
	case "windows":
		sshDir = os.Getenv(HOMEDRIVE) + os.Getenv(HOMEPATH)
	case "linux":
		sshDir = os.Getenv(HOME)
	}

}

func initTemplateFileLocation() {
	templateFileLocation = appConfig.TemplatesDirectory + sPathseparator + "readme.md"
}

func initNewRepoTemplateFileLocation() {
	newRepoTemplateFileLocation = appConfig.TemplatesDirectory + sPathseparator + "0001" + sPathseparator + "readme.md"
}

func displayEnvironment() {
	operatingSystem := runtime.GOOS
	fmt.Println("OS: " + operatingSystem)
	switch operatingSystem {
	case "windows":
		fmt.Println("HOMEDRIVE=" + os.Getenv(HOMEDRIVE))
		fmt.Println("HOMEPATH =" + os.Getenv(HOMEPATH))
	case "linux":
		fmt.Println("HOME=" + os.Getenv(HOME))
	}
	fmt.Println("RFD root directory=" + appConfig.RootDirectory)
	//fmt.Println("RFD relative directory=" + appConfig.RFDRelativeDirectory)
	fmt.Println("Installation directory=" + appConfig.TemplatesDirectory)
	fmt.Println("SSH public key directory=" + getSSHPath())

	publicKey, err := getPublicKey()
	CheckFatal(err)

	bytes := publicKey.Signer.PublicKey().Marshal()

	sPublicKey := base64.StdEncoding.EncodeToString(bytes) + " " + publicKey.User
	CheckFatal(err)

	fmt.Println("SSH Public Key=" + sPublicKey)

	println()
	println()

	fmt.Printf("%+v", appConfig)
}

func populateConfig() (*configuration, error) {

	err := checkConfigurationFile()
	if err != nil {
		return nil, err
	}

	// Create appConfig structure
	config := &configuration{}

	// Open appConfig file
	file, err := os.Open("./config.yml")
	CheckFatal(err)

	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&config)
	CheckFatal(err)

	return config, err
}

func getConfiguredStates() (*states, error) {

	err := checkStatesFile()
	if err != nil {
		return nil, err
	}

	// Create states structure
	states := &states{}

	// Open appConfig file
	file, err := os.Open(appConfig.TemplatesDirectory + "/states.yml")
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

	if _, err := os.Stat(appConfig.TemplatesDirectory + "/states.yml"); os.IsNotExist(err) {
		fmt.Println("templateFileLocation: " + appConfig.TemplatesDirectory)
		fmt.Println("States file does not exist... ")
		return err
	}
	return nil
}
func checkAndReportOnRepositoryState() bool {

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
