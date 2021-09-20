package main

import (
	"encoding/base64"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"runtime"
)

const HOME string = "HOME"
const HOMEDRIVE string = "HOMEDRIVE"
const HOMEPATH string = "HOMEPATH"

var sshDir string
var templateFileLocation string
var sPathseparator string

var logger Trace
var appConfig *configuration

/* 	TODO/IDEA
A set of functions that will be called to roll back, in case of a fatal error.
Currently just a set of empty functions and cleaning up requires a bit of manual git commandline
work.

The idea is that functions that need to be called to roll back a request are added to a queue.
These are then popped off in reverse order and executed if there is a fatal error.
*/
var rollbackFunctions []func()

type configuration struct {
	RFDRootDirectory     string `yaml:"rfd-root-directory"`
	InstallDirectory     string `yaml:"install-directory"`
	RFDRelativeDirectory string `yaml:"rfd-relative-directory"`
	PrivateKeyFileName   string `yaml:"private-key-file-name"`
}

func init() {

	sPathseparator = string(os.PathSeparator)

	rollbackFunctions = append(rollbackFunctions, undoCreateReadme)
	rollbackFunctions = append(rollbackFunctions, undoCreateBranch)

	logger = TraceLog{}
	appConfig = populateConfig()
	initSSHDIR()
	initTemplateFileLocation()

	err := checkConfig()
	CheckFatal(err)
}

func main() {
	app := createCommandLineApp()
	err := app.Run(os.Args)
	CheckFatal(err)
}

//
func createCommandLineApp() *cli.App {
	app := &cli.App{
		Name:  "rfd",
		Usage: "Create new rfd's, index and output their status, and manage their .",
		Commands: []*cli.Command{
			{
				Name:  "check",
				Usage: "Check environment is suitable to ensure a clean run of 'rfd new'",
				Action: func(c *cli.Context) error {
					checkAndReportOnRepositoryState()
					return nil
				},
			},
			{
				Name:  "index",
				Usage: "Output the status of all rfd's to index.md in markdown format.",
				Action: func(c *cli.Context) error {
					Index()
					return nil
				},
			},
			{
				Name:  "new",
				Usage: "Create a new rfd",
				Action: func(c *cli.Context) error {
					if checkAndReportOnRepositoryState() {
						new()
					} else {
						fmt.Println("Creating a new RFD creates and switches to new branch. Please commit (or otherwise) unstaged and/or uncommitted work.")
						fmt.Println()
					}
					return nil
				},
			},

			{
				Name:  "init",
				Usage: "Initialise an RFD repository. Note - the repository must be empty, and it is assumed that there is a remote target repository.",
				Action: func(c *cli.Context) error {
					initRepo()
					return nil
				},
			},
			{
				Name:  "environment",
				Usage: "Displays configuration settings and relevant operating system environment variables.",
				Action: func(c *cli.Context) error {
					displayEnvironment()
					return nil
				},
			},
			{
				Name:  "transition",
				Usage: "Displays the status of <nnnn>. Will output the status of every RFD if it isn't provided an RFD ID.",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:  "show",
				Usage: "Displays the status of <nnnn>. Will output the status of every RFD if it isn't provided an RFD ID.",
				Action: func(c *cli.Context) error {
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
	templateFileLocation = appConfig.InstallDirectory + sPathseparator + "template" + sPathseparator + "readme.md"
}

func displayEnvironment() {
	operatingSystem := runtime.GOOS
	fmt.Println(operatingSystem)
	switch operatingSystem {
	case "windows":
		fmt.Println("HOMEDRIVE=" + os.Getenv(HOMEDRIVE))
		fmt.Println("HOMEPATH =" + os.Getenv(HOMEPATH))
	case "linux":
		fmt.Println("HOME=" + os.Getenv(HOME))
	}
	fmt.Println("RFD root directory=" + appConfig.RFDRootDirectory)
	fmt.Println("RFD relative directory=" + appConfig.RFDRelativeDirectory)
	fmt.Println("Installation directory=" + appConfig.InstallDirectory)
	fmt.Println("SSH public key directory=" + getSSHPath())

	publicKey, err := getPublicKey()
	CheckFatal(err)

	bytes := publicKey.Signer.PublicKey().Marshal()

	sPublicKey := base64.StdEncoding.EncodeToString(bytes) + " " + publicKey.User
	CheckFatal(err)

	fmt.Println("SSH Public Key=" + sPublicKey)
}

func populateConfig() *configuration {
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

	return config
}

func checkConfig() error {

	// Check location of template file is correct
	_, err := ioutil.ReadFile(templateFileLocation)
	if err != nil {
		fmt.Println("Attempted to read " + templateFileLocation)
		fmt.Println("Can't read readme template file. Please check the config.yml file.\n")
	}

	return err
}

func checkAndReportOnRepositoryState() bool {

	var fileStatusMapping = map[git.StatusCode]string{
		git.Unmodified:         "Unmodified",
		git.Untracked:          "Untracked",
		git.Modified:           "Modified",
		git.Added:              "Added",
		git.Deleted:            "Deleted",
		git.Renamed:            "Renamed",
		git.Copied:             "Copied",
		git.UpdatedButUnmerged: "Updated",
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
